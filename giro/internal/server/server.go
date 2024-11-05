package server

import (
	"encoding/json"
	"fmt"
	"giro/internal/models"
	"giro/internal/utils"
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
	PORT               = "8888"
	CLIPORT            = ":7771"
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
	http.HandleFunc("/server/ticket/purchase", s.HandleServerTicketPurchase)
	http.HandleFunc("/server/ticket/cancel", s.HandleServerTicketCancel)
	http.HandleFunc("/server/broadcast", s.HandleBroadcast)

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

// shutdownServer gracefully shuts down the server by waiting for all goroutines to finish,
// saving the system variables to a file, and closing the HTTP server.
//
// Parameters:
// - server: A pointer to the http.Server instance representing the server.
//
// Return:
// - An error if the server fails to close gracefully. Returns nil if the shutdown is successful.
func (s *System) shutdownServer(server *http.Server) error {
	log.Println("Server shutting down...")

	// Wait for all goroutines to finish
	s.wg.Wait()

	// Lock the mutex to ensure atomic access when writing to file
	s.Lock.Lock()
	defer s.Lock.Unlock()

	// Save the system variables to a file
	s.storeSystemVars(INSTANCE_PATH)

	// Close the HTTP server
	return server.Close()
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
