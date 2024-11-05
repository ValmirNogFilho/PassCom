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
	Log         []models.LogMessage
	Buffer      chan models.LogMessage
	VectorClock map[string]int
	Connections map[string]models.Connection
	Lock        sync.RWMutex
	wg          sync.WaitGroup // WaitGroup para controlar goroutines
	shutdown    chan os.Signal // Canal para sinalizar o encerramento
}

const (
	SERVER_NAME        = "giro"
	ADDRESS            = "localhost"
	PORT               = "9999"
	CLIPORT            = ":7772"
	INSTANCE_PATH      = "systemvars.json"
	BUFFER_SIZE        = 100
	LOG_SIZE           = 1000
	CONNECTION_TIMEOUT = 10 * time.Second
	HEARTBEAT_TIMER    = 1 * time.Second
	SESSION_TIME_LIMIT = 30 * time.Minute
	URL_PREFIX         = "http://"
)

const (
	EQUAL = iota
	CONCURRENT
	NEWER
	OLDER
)

var (
	instance *System
	once     sync.Once
)

// storeSystemVars saves the system variables to a JSON file named "systemvars.json".
// It creates the file if it doesn't exist, overwriting any existing content.
// The function serializes the system variables into a JSON object and writes it to the file.
// If any error occurs during file creation, writing, or serialization, the function logs the error and terminates the program.
func (s *System) storeSystemVars(path string) {
	file, err := os.Create(path)
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	systemVars := make(map[string]interface{})
	systemVars["ServerName"] = s.ServerName
	systemVars["ServerId"] = s.ServerId
	systemVars["Address"] = s.Address
	systemVars["Log"] = s.Log
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

// LoadInstanceFromFile reads the system variables from a JSON file named "systemvars.json" and returns a new System instance.
// If the file does not exist or cannot be read, it returns nil and an error.
//
// The function performs the following steps:
// 1. Opens the "systemvars.json" file.
// 2. If an error occurs while opening the file, it returns nil and the error.
// 3. Decodes the JSON content into a new System instance.
// 4. If an error occurs while decoding the JSON, it returns nil and the error.
// 5. If the file exists and can be read successfully, it restores the channels and maps in the loaded instance.
// 6. Finally, it returns a pointer to the loaded System instance and nil as the error.
func LoadInstanceFromFile(path string) (*System, error) {
	file, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}

	var loadedInstance System
	if err := json.Unmarshal(file, &loadedInstance); err != nil {
		return nil, err
	}

	// Restaurar canais e mapas, se necessário
	loadedInstance.Address = getLocalIP()
	loadedInstance.Port = getPort()
	loadedInstance.Buffer = make(chan models.LogMessage, BUFFER_SIZE)
	loadedInstance.shutdown = make(chan os.Signal, 1)

	return &loadedInstance, nil
}

// GetInstance returns a singleton instance of the System struct.
// It initializes the instance with default values if it hasn't been created yet.
// The function uses a sync.Once to ensure that the instance is only created once.
//
// The function returns a pointer to the System instance.
func GetInstance() *System {
	once.Do(func() {
		loadedInstance, err := LoadInstanceFromFile(INSTANCE_PATH)
		if err == nil {
			log.Printf("Server instance loaded from file %v", INSTANCE_PATH)
			instance = loadedInstance
		} else {
			log.Printf("Failed to load server instance from file %v due to: %v. Creating new instance...", INSTANCE_PATH, err)
			instance = &System{
				ServerName:  SERVER_NAME,
				ServerId:    uuid.New(),
				Address:     getLocalIP(),
				Port:        getPort(),
				Log:         make([]models.LogMessage, LOG_SIZE),
				Buffer:      make(chan models.LogMessage, BUFFER_SIZE),
				VectorClock: make(map[string]int),
				Connections: make(map[string]models.Connection),
				shutdown:    make(chan os.Signal, 1),
			}

			instance.VectorClock[instance.ServerId.String()] = 0
		}
	})
	return instance
}

// StartServer initializes and starts the server, handling HTTP requests and messages.
// It also listens for system signals to gracefully shut down the server.
//
// The function starts a cleanup goroutine to remove expired sessions.
// It registers HTTP handlers for client requests and server messages.
// It sets up an HTTP server with the specified address and timeouts.
// It starts goroutines to send heartbeats, handle CLI connections, and listen for system signals.
//
// The function returns an error if the server fails to start or if an error occurs during shutdown.
func (s *System) StartServer() error {
	signal.Notify(s.shutdown, syscall.SIGINT, syscall.SIGTERM)

	go s.CleanupSessions(SESSION_TIME_LIMIT)

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
	http.HandleFunc("/server/heartbeat", s.handleHeartbeat)
	http.HandleFunc("/server/connect", s.handleConnect)
	http.HandleFunc("/server/database", s.handleDatabase)

	httpServer := &http.Server{
		Addr:         s.Address + ":" + s.Port,
		ReadTimeout:  CONNECTION_TIMEOUT,
		WriteTimeout: CONNECTION_TIMEOUT,
	}

	log.Println("HTTP Server listening on", httpServer.Addr)

	errCh := make(chan error, 1)
	go func() {
		if err := httpServer.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			errCh <- err
		}
	}()

	go s.sendHeartbeats()

	go s.HandleCLIServer()

	select {
	case <-s.shutdown:
		err := s.shutdownServer(httpServer)
		if err != nil {
			log.Fatalf("Error shutting down server gracefully: %v", err)
		}
		return nil
	case err := <-errCh:
		log.Fatalf("Server error: %v", err)
		return err
	}
}

func (s *System) shutdownServer(server *http.Server) error {
	log.Println("Server shutting down...")

	// Espere todas as goroutines terminarem
	s.wg.Wait()

	// Bloqueie o mutex para escrever em arquivo
	s.Lock.Lock()
	defer s.Lock.Unlock()

	// Salve os dados do sistema em um arquivo
	s.storeSystemVars(INSTANCE_PATH)

	// Fecha o servidor
	return server.Close()
}

// HandleCLIServer starts a TCP server listening on the specified CLIPORT.
// This server accepts incoming connections and handles them using the handleCLIConnection function.
// The server logs any errors during initialization, listening, or accepting connections.
//
// The function performs the following steps:
// 1. Listens for incoming TCP connections on the specified CLIPORT.
// 2. If an error occurs during listening, logs the error and terminates the program.
// 3. Accepts incoming connections and starts a new goroutine to handle each connection using the handleCLIConnection function.
// 4. Logs any errors that occur during connection acceptance.
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

// handleCLIConnection handles incoming connections from the CLI server.
// It reads commands from the connection, processes them, and responds accordingly.
//
// Parameters:
// - conn: The net.Conn object representing the connection from the client.
//
// Return:
// This function does not return any value.
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
					"\n- addconn <address> <port>: to add a new connection" +
					"\n- quit: to close the connection" +
					"\n- shutdown: to shut down the server\n"))

		case "info":
			conn.Write([]byte(s.getServerInfo()))

		case "addconn":
			if len(args) < 2 {
				conn.Write([]byte("Error: 'addconn' requires two arguments (address, port).\n"))
			} else {
				address := args[0]
				connPort := args[1]
				s.RequestConnection(address, connPort)
				conn.Write([]byte("Requesting connection to " + address + ":" + connPort + "...\n"))
				s.RequestDatabase(address, connPort)
				conn.Write([]byte("Requesting database from " + address + ":" + connPort + "...\n"))
			}

		case "rmconn":
			if len(args) < 1 {
				conn.Write([]byte("Error: 'rmconn' requires one argument (connection ID).\n"))
			} else {
				connId := args[0]
				id, serverConn := s.FindConnectionByName(connId)
				if serverConn == nil {
					conn.Write([]byte("Connection not found.\n"))
				} else {
					name := serverConn.Name
					s.RequestDisconnection(serverConn.Address, serverConn.Port)
					conn.Write([]byte("Requested disconnection from " + serverConn.Address + ":" + serverConn.Port + "...\n"))
					s.RemoveConnection(id)
					conn.Write([]byte("Removing connection from " + name + "...\n"))
					s.RemoveDatabase(connId)
					conn.Write([]byte("Removing database from " + name + "...\n"))
				}
			}

		case "quit":
			conn.Write([]byte("Closing CLI...\n"))
			return

		case "shutdown":
			conn.Write([]byte("Shutting down the server...\n"))
			go func() {
				s.shutdown <- syscall.SIGTERM
			}()
			return

		default:
			conn.Write([]byte("Command not found.\n"))
		}
	}
}

// getServerInfo returns a formatted string containing information about the server.
// It includes the server's name, address, port, server ID, connections, and vector clock.
// The function is safe for concurrent use and ensures that the server's data is read atomically.
//
// Parameters:
// - s: A pointer to the System struct representing the server.
//
// Return:
// - A string containing the server information.
func (s *System) getServerInfo() string {

	s.Lock.RLock()
	defer s.Lock.RUnlock()
	return fmt.Sprintf("\nName: %s\nAddress: %s\nPort: %s\nServerId: %s\n"+
		"Connections: %v\nVector Clock:%v\n", s.ServerName, s.Address, s.Port,
		s.ServerId, utils.PrintMap(s.Connections),
		utils.PrintMap(s.VectorClock))
}
