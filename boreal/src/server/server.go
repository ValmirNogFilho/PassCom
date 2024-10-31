package server

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"passcom/boreal/src/models"
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
	Lock        sync.Mutex
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

// TODO: Verificar se a mensagem foi redirecionada para evitar loops
func (s *System) handleGetMessage(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	if err := json.NewEncoder(w).Encode(map[string]string{"status": "Message received"}); err != nil {
		http.Error(w, "Failed to encode response", http.StatusInternalServerError)
	}
}

func (s *System) handlePostMessage(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	var msg models.Message
	if err := json.NewDecoder(r.Body).Decode(&msg); err != nil {
		http.Error(w, "Failed to decode request", http.StatusBadRequest)
		return
	}
	fmt.Println("Message received:", msg)

	if err := json.NewEncoder(w).Encode(map[string]string{"status": "Message received"}); err != nil {
		http.Error(w, "Failed to encode response", http.StatusInternalServerError)
	}
}

func (s *System) handleHTTPMessage(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()

	switch r.Method {
	case http.MethodGet:
		s.handleGetMessage(w, r)
	case http.MethodPost:
		s.handlePostMessage(w, r)
	default:
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}
}

func (s *System) handleHeartbeat(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	ip := r.RemoteAddr

	log.Printf("Heartbeat check successful from %s\n", ip)
}

func (s *System) IsReceivedClockNewer(clock map[uuid.UUID]int) int {
	// Assuma que VC(x) é o relógio do sistema e VC(y) é o relógio recebido.
	// VC(x) < VC(y) ⇔ ∀z[VC(x)[z] ≤ VC(y)[z]] e ∃w[VC(x)[w] < VC(y)[w]]
	// Lê-se: VC(x) é mais antigo que VC(y) se, para todo z de VC(x), eles são menores ou iguais
	// para o correspondente z em VC(y), e existe um w onde VC(x)[w] é estritamente menor que VC(y)[w].
	// Se a condição acima for atendida, então VC(y) é mais novo que VC(x).
	// Se a condição acima for atendida para VC(y), então VC(y) é mais antigo que VC(x).
	// Senão, se ∃z'[VC(x)[z'] > VC(y)[z']], então VC(x) e VC(y) são concorrentes.
	// Senão, VC(x) e VC(y) são iguais.

	vx := s.VectorClock
	vy := clock

	isLess := false
	isGreater := false

	// Itera sobre as chaves em ambos os relógios
	for id, x := range vx {
		y, exists := vy[id]
		if !exists {
			y = 0 // Se o ID não existe em vy, assume que seu valor é 0
		}

		if x < y {
			isLess = true
		} else if x > y {
			isGreater = true
		}

		// Se ambos isLess e isGreater são verdadeiros, eles são concorrentes
		if isLess && isGreater {
			return CONCURRENT
		}
	}

	// Verifica quaisquer IDs em vy que estão ausentes em vx
	for id, y := range vy {
		if _, exists := vx[id]; !exists && y > 0 {
			isLess = true
			if isGreater {
				return CONCURRENT
			}
		}
	}

	// Agora determina a relação com base nas flags
	if isLess {
		return NEWER
	}
	if isGreater {
		return OLDER
	}

	// Se nenhuma flag foi definida, então eles são iguais
	return EQUAL
}

func (s *System) updateVectorClock(receivedClock map[uuid.UUID]int) {
	for id, time := range receivedClock {
		if _, exists := s.VectorClock[id]; !exists || time > s.VectorClock[id] {
			s.VectorClock[id] = time
		}
	}
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
