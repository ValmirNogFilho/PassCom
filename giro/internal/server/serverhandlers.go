package server

import (
	"encoding/json"
	"giro/internal/models"
	"log"
	"net/http"

	"github.com/google/uuid"
)

func (s *System) handleHeartbeat(w http.ResponseWriter, r *http.Request) {
	// Define o tipo de resposta e status OK
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	// Decodifica o *heartbeat* recebido
	var receivedMessage models.Message
	if err := json.NewDecoder(r.Body).Decode(&receivedMessage); err != nil {
		http.Error(w, "Invalid request payload", http.StatusBadRequest)
		return
	}

	// Bloqueia o mutex para manipular o VectorClock de maneira segura
	s.Lock.Lock()
	defer s.Lock.Unlock()

	// Incrementa o relógio do sistema para indicar que o heartbeat foi recebido
	s.IncrementClock()

	log.Print("Received heartbeat from ", receivedMessage.From, "with VectorClock: ", receivedMessage.VectorClock)

	// Atualiza o VectorClock com base no *heartbeat* recebido
	s.UpdateClock(receivedMessage.VectorClock)

	// Cria uma nova mensagem de resposta com o VectorClock atualizado
	responseMessage := models.Message{
		Id:          uuid.New().String(),
		From:        s.ServerId.String(),
		To:          receivedMessage.From,
		VectorClock: s.VectorClock,
		Body:        "Healthy",
	}

	// Codifica a mensagem de resposta como JSON
	if err := json.NewEncoder(w).Encode(responseMessage); err != nil {
		log.Printf("Error encoding heartbeat response: %v", err)
	}

	log.Print("Sent heartbeat response to ", receivedMessage.From)
}
