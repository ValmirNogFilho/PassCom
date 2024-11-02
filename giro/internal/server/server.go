package server

import (
	"fmt"
	"giro/internal/dao"
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
	ServerId    uuid.UUID
	Buffer      chan models.Message
	VectorClock map[uuid.UUID]int
	Connections map[uuid.UUID]models.Connection
	Lock        sync.RWMutex
	wg          sync.WaitGroup // WaitGroup para controlar goroutines
	shutdown    chan os.Signal // Canal para sinalizar o encerramento
}

const (
	ADDRESS     = ""
	PORT        = ":9999"
	BUFFER_SIZE = 100
	TIMEOUT     = 10 * time.Second
	timeLimit   = 30 * time.Minute
	EQUAL       = iota
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
			ServerId:    uuid.New(),
			Buffer:      make(chan models.Message, BUFFER_SIZE),
			VectorClock: make(map[uuid.UUID]int),
			Connections: make(map[uuid.UUID]models.Connection),
			shutdown:    make(chan os.Signal, 1),
		}
	})
	return instance
}

func (s *System) StartServer() error {
	signal.Notify(s.shutdown, syscall.SIGINT, syscall.SIGTERM)

	go CleanupSessions(timeLimit)

	http.HandleFunc("/login", handleLogin)
	http.HandleFunc("/logout", handleLogout)
	http.HandleFunc("/user", handleGetUser)
	http.HandleFunc("/route", handleGetRoute)
	http.HandleFunc("/flights", handleGetFlights)
	http.HandleFunc("/ticket", handleTicket)
	http.HandleFunc("/tickets", handleGetTickets)
	http.HandleFunc("/airports", handleGetAirports)
	http.HandleFunc("/heartbeat", s.handleHeartbeat)

	server := &http.Server{
		Addr:         ADDRESS + PORT,
		ReadTimeout:  TIMEOUT,
		WriteTimeout: TIMEOUT,
	}

	log.Println("HTTP Server listening on", server.Addr)

	errCh := make(chan error, 1)
	go func() {
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			errCh <- err
		}
	}()

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

// CleanupSessions periodically checks for inactive sessions and reservations, and cleans them up.
// It runs every minute and checks each session and its reservations against the provided timeout.
// If a session or a reservation is inactive (i.e., its last activity time is older than the timeout),
// it is deleted from the system.
//
// Parameters:
//   - timeout: The duration after which a session or a reservation is considered inactive.
func CleanupSessions(timeout time.Duration) {
	ticker := time.NewTicker(time.Minute)
	defer ticker.Stop()

	for range ticker.C {
		for _, session := range dao.GetSessionDAO().FindAll() {
			if time.Since(session.LastTimeActive) > timeout {
				fmt.Printf("Encerrando sessão %s por inatividade\n", session.ID)
				dao.GetSessionDAO().Delete(session)
			}
		}
	}
}

// SessionIfExists checks if a session exists for the given token.
// If a session is found, it updates the session's last activity time and returns the session along with a boolean value of true.
// If no session is found or an error occurs during the process, it returns nil and false.
//
// Parameters:
//   - token: A string representing the session token to be checked.
//
// Return:
//   - *models.Session: A pointer to the found session if it exists, or nil if no session is found or an error occurs.
//   - bool: A boolean value indicating whether a session was found (true) or not (false).
func SessionIfExists(token string) (*models.Session, bool) {
	uuid, err := uuid.Parse(token)
	if err != nil {
		return nil, false
	}
	session, err := dao.GetSessionDAO().FindById(uuid)
	if err != nil {
		return nil, false
	}
	session.LastTimeActive = time.Now()
	dao.GetSessionDAO().Update(session)
	return session, true
}
