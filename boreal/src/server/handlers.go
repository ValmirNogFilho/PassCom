package server

import (
	"boreal/src/models"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
)

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
