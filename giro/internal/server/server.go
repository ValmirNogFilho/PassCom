package server

import (
	"bufio"
	"encoding/json"
	"fmt"
	"giro/internal/models"
	"giro/internal/utils"
	"log"
	"net"
	"net/http"
	"os"
	"os/signal"
	"strings"
	"sync"
	"syscall"
	"time"

	"github.com/google/uuid"
)

type System struct {
	ServerName  string
	ServerId    uuid.UUID
	Address     string
	Port        string
	Buffer      chan models.Message
	VectorClock map[string]int
	Connections map[string]models.Connection
	Lock        sync.RWMutex
	wg          sync.WaitGroup // WaitGroup para controlar goroutines
	shutdown    chan os.Signal // Canal para sinalizar o encerramento
}

const (
	SERVERNAME        = "giro"
	PORT              = "9999"
	CLIPORT           = ":7772"
	BUFFERSIZE        = 100
	CONNECTIONTIMEOUT = 10 * time.Second
	HEARTBEATTIMER    = 1 * time.Second
	SESSIONTIMELIMIT  = 30 * time.Minute
	URLPREFIX         = "http://"
	EQUAL             = iota
	CONCURRENT
	NEWER
	OLDER
)

var (
	instance *System
	once     sync.Once
)

func GetInstance() *System {
	once.Do(func() {
		instance = &System{
			ServerName:  SERVERNAME,
			ServerId:    uuid.New(),
			Address:     getLocalIP(),
			Port:        getPort(),
			Buffer:      make(chan models.Message, BUFFERSIZE),
			VectorClock: make(map[string]int),
			Connections: make(map[string]models.Connection),
			shutdown:    make(chan os.Signal, 1),
		}

		instance.VectorClock[instance.ServerId.String()] = 0
	})
	return instance
}

func (s *System) StartServer() error {
	signal.Notify(s.shutdown, syscall.SIGINT, syscall.SIGTERM)

	go s.CleanupSessions(SESSIONTIMELIMIT)

	// Usam requests dos clientes
	http.HandleFunc("/login", handleLogin)
	http.HandleFunc("/logout", handleLogout)
	http.HandleFunc("/user", handleGetUser)
	http.HandleFunc("/route", handleGetRoute)
	http.HandleFunc("/flights", handleGetFlights)
	http.HandleFunc("/ticket", handleTicket)
	http.HandleFunc("/tickets", handleGetTickets)
	http.HandleFunc("/airports", handleGetAirports)
	http.HandleFunc("/wishlist", handleWishlist)

	// Usam messages dos servidores
	http.HandleFunc("/heartbeat", s.handleHeartbeat)
	http.HandleFunc("/connect", s.handleConnect)
	http.HandleFunc("/request/connect", s.handleRequestConnection)

	server := &http.Server{
		Addr:         s.Address + ":" + s.Port,
		ReadTimeout:  CONNECTIONTIMEOUT,
		WriteTimeout: CONNECTIONTIMEOUT,
	}

	log.Println("HTTP Server listening on", server.Addr)

	errCh := make(chan error, 1)
	go func() {
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			errCh <- err
		}
	}()

	go s.sendHeartbeats()

	go s.HandleCLIServer()

	select {
	case <-s.shutdown:
		s.storeSystemVars()
		log.Println("Server shutting down...")
		s.wg.Wait()
		return server.Close()
	case err := <-errCh:
		log.Fatalf("server error: %v", err)
		return err
	}
}

func (s *System) storeSystemVars() {
	file, err := os.Create("systemvars.json")
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	systemVars := make(map[string]interface{})
	systemVars["ServerName"] = s.ServerName
	systemVars["ServerId"] = s.ServerId
	systemVars["Address"] = s.Address
	systemVars["Port"] = s.Port
	systemVars["VectorClock"] = s.VectorClock
	systemVars["Connections"] = s.Connections

	jsonData, err := json.MarshalIndent(systemVars, "", "  ") // identação
	if err != nil {
		log.Fatal("Error identing JSON:", err)
	}

	if _, err := file.Write(jsonData); err != nil {
		log.Fatal("Error serializing vars to JSON:", err)
	}

	log.Println("System vars saved.")
}

func (s *System) HandleCLIServer() {
	listener, err := net.Listen("tcp", CLIPORT)
	if err != nil {
		log.Fatal("Error initiating server:", err)
	}
	defer listener.Close()

	log.Println("See your server working on the TCP CLI server on port", CLIPORT)

	for {
		conn, err := listener.Accept()
		if err != nil {
			log.Println("Error accepting connection:", err)
			continue
		}

		go s.handleCLIConnection(conn)
	}
}

func (s *System) handleCLIConnection(conn net.Conn) {
	defer conn.Close()
	conn.Write([]byte("Welcome to the CLI server!\nType 'help' to see the commands.\n"))

	scanner := bufio.NewScanner(conn)
	for scanner.Scan() {
		input := strings.TrimSpace(scanner.Text())
		parts := strings.Fields(input)

		if len(parts) == 0 {
			conn.Write([]byte("Empty command.\n"))
			continue
		}

		command := parts[0]
		args := parts[1:]

		switch command {
		case "help":
			conn.Write(
				[]byte("Available commands:" +
					"\n- help: to see commands" +
					"\n- info: to see server informations" +
					"\n- addconn <server_name> <address> <port>: to add a new connection" +
					"\n- quit: to close the connection\n"))

		case "info":
			conn.Write([]byte(s.getServerInfo()))

		case "addconn":
			if len(args) < 1 {
				conn.Write([]byte("Error: 'addconn' requires three arguments (server name, address, port).\n"))
			} else {
				name := args[0]
				address := args[1]
				connPort := args[2]
				if s.ServerName != name {
					s.RequestConnection(name, address, connPort)
					conn.Write([]byte("new connection added with port: " + connPort + "\n"))
				} else {
					conn.Write([]byte("self connection not allowed. \n"))
				}
			}

		case "quit":
			conn.Write([]byte("Closing CLI...\n"))
			return

		default:
			conn.Write([]byte("Command not found.\n"))
		}
	}

}

func (s *System) getServerInfo() string {

	s.Lock.RLock()
	defer s.Lock.RUnlock()
	return fmt.Sprintf("\nName: %s\nAddress: %s\nPort: %s\nServerId: %s\n"+
		"Connections: %v\nVector Clock:%v\n", s.ServerName, s.Address, s.Port,
		s.ServerId, utils.PrintMap(s.Connections),
		utils.PrintMap(s.VectorClock))
}
