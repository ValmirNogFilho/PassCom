package server

import (
	"boreal/internal/models"
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
	ADDRESS     = "localhost"
	PORT        = ":8080"
	BUFFER_SIZE = 100
	TIMEOUT     = 10 * time.Second
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

	http.HandleFunc("/heartbeat", s.handleHeartbeat)
	http.HandleFunc("/messages", s.handleHTTPMessage)

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
