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
	serverName        = "giro"
	port              = ":9999"
	bufferSize        = 100
	connectionTimeout = 10 * time.Second
	heartbeatTimer    = 1 * time.Second
	sessionTimeLimit  = 30 * time.Minute
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
			ServerName:  serverName,
			ServerId:    uuid.New(),
			Address:     getLocalIP(),
			Port:        port,
			Buffer:      make(chan models.Message, bufferSize),
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

	go s.CleanupSessions(sessionTimeLimit)

	// Usam requests dos clientes
	http.HandleFunc("/login", handleLogin)
	http.HandleFunc("/logout", handleLogout)
	http.HandleFunc("/user", handleGetUser)
	http.HandleFunc("/route", handleGetRoute)
	http.HandleFunc("/flights", handleGetFlights)
	http.HandleFunc("/ticket", handleTicket)
	http.HandleFunc("/tickets", handleGetTickets)
	http.HandleFunc("/airports", handleGetAirports)

	// Usam messages dos servidores
	http.HandleFunc("/heartbeat", s.handleHeartbeat)
	http.HandleFunc("/connect", s.handleConnect)

	server := &http.Server{
		Addr:         s.Address + s.Port,
		ReadTimeout:  connectionTimeout,
		WriteTimeout: connectionTimeout,
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
