package server

import (
	"giro/internal/models"
	"log"
	"net/http"
	"os"
	"os/signal"
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

	select {
	case <-s.shutdown:
		log.Println("Server shutting down...")
		s.wg.Wait()
		return server.Close()
	case err := <-errCh:
		log.Fatalf("server error: %v", err)
		return err
	}
}
