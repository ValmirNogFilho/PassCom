package server

import (
	"encoding/json"
	"fmt"
	"net"
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
	Connections map[uuid.UUID]net.Conn
	Lock        sync.Mutex
	wg          sync.WaitGroup // WaitGroup para controlar goroutines
	shutdown    chan os.Signal // Canal para sinalizar o encerramento
}

const (
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
			Buffer:      make(chan models.Message, 100), // Exemplo de tamanho de buffer
			VectorClock: make(map[uuid.UUID]int),
			Connections: make(map[uuid.UUID]net.Conn),
			shutdown:    make(chan os.Signal, 1),
		}
	})
	return instance
}

func (s *System) AddConnection(serverId uuid.UUID, conn net.Conn) {
	s.Lock.Lock()
	defer s.Lock.Unlock()
	s.Connections[serverId] = conn
}

func (s *System) handleGetMessage(w http.ResponseWriter, r *http.Request) {
	// TODO: Verificar se a mensagem foi redirecionada para evitar o loop

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

func (s *System) IsReceivedClockNewer(clock map[uuid.UUID]int) int {
	// Assume que VCx é o relógio do sistema e VCy é o relógio recebido.
	// VCx < VCy ⇔ ∀z[VCx[z] ≤ VCy[z]] e ∃z[VCx[w] < VCy[w]]
	// lê-se: VCx é mais antigo que VCy se, para todo z de VCx, eles são menores ou iguais
	// para o correspondente z em VCy, e existe um w onde VCx[w] é estritamente menor que VCy[w].
	// se a condição acima for atendida, então VCx é mais antigo que VCy.
	// se a condição acima for atendida para VCy, então o relógio do sistema é mais novo que VCy.
	// caso contrário, se ∃[VCx[z] > VCy[z]], então VCx e VCy são concorrentes.

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
		if time > s.VectorClock[id] {
			s.VectorClock[id] = time
		}
	}
}

func (s *System) StartServer() error {
	signal.Notify(s.shutdown, syscall.SIGINT, syscall.SIGTERM)

	http.HandleFunc("/messages", s.handleHTTPMessage)

	fmt.Println("HTTP Server listening on", PORT)
	server := &http.Server{
		Addr:         PORT,
		ReadTimeout:  TIMEOUT,
		WriteTimeout: TIMEOUT,
	}

	errCh := make(chan error, 1)
	go func() {
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			errCh <- err
		}
	}()

	select {
	case <-s.shutdown:
		fmt.Println("Server shutting down...")
		s.wg.Wait()
		return server.Close()
	case err := <-errCh:
		return fmt.Errorf("server error: %v", err)
	}
}
