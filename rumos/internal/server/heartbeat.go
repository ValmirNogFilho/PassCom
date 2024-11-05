package server

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"rumos/internal/models"
	"time"
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

	log.Print("Received heartbeat from ", s.Connections[receivedMessage.From].Name)

	// Atualiza o VectorClock com base no *heartbeat* recebido
	s.UpdateClock(receivedMessage.VectorClock)

	// Cria uma nova mensagem de resposta com o VectorClock atualizado
	responseMessage, err := models.CreateMessage(s.ServerId.String(), receivedMessage.From, s.VectorClock, "Healthy")

	if err != nil {
		log.Printf("Error creating heartbeat response message: %v", err)
		return
	}

	// Codifica a mensagem de resposta como JSON
	if err := json.NewEncoder(w).Encode(responseMessage); err != nil {
		log.Printf("Error encoding heartbeat response: %v", err)
	}

	log.Print("Sent heartbeat response to ", s.Connections[receivedMessage.From].Name)
}

func (s *System) sendHeartbeats() {
	ticker := time.NewTicker(5 * time.Second) // Intervalo de envio do heartbeat
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			s.Lock.Lock()
			s.IncrementClock()
			for id, conn := range s.Connections {
				s.wg.Add(1)
				go s.sendHeartbeatToConnection(id, conn)
			}
			s.Lock.Unlock()
		case <-s.shutdown:
			return
		}
	}
}

// Função para enviar um heartbeat a uma única conexão e atualizar o status
func (s *System) sendHeartbeatToConnection(id string, conn models.Connection) {
	defer s.wg.Done()

	heartbeat, err := models.CreateMessage(s.ServerId.String(), id, s.VectorClock, "Heartbeat")

	if err != nil {
		log.Printf("Error creating heartbeat message: %v", err)
		return
	}

	// Serializar a mensagem de heartbeat
	jsonData, err := json.Marshal(heartbeat)
	if err != nil {
		log.Printf("Error encoding heartbeat message: %v", err)
		return
	}

	// Construir a URL com endereço e porta
	url := fmt.Sprintf("%s%s:%s/server/heartbeat", URL_PREFIX, conn.Address, conn.Port)
	req, err := http.NewRequest("POST", url, bytes.NewBuffer(jsonData))
	if err != nil {
		log.Printf("Error creating heartbeat request: %v", err)
		s.UpdateConnectionStatus(id, false)
		return
	}
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{Timeout: 3 * time.Second} // Define timeout para a resposta

	log.Printf("Sending heartbeat to %s", s.Connections[id].Name)
	resp, err := client.Do(req)

	// Verifica se houve erro na resposta ou se o servidor está offline
	if err != nil || resp.StatusCode != http.StatusOK {
		log.Printf("Connection %s is offline", s.Connections[id].Name)
		s.UpdateConnectionStatus(id, false)
	} else {
		s.UpdateConnectionStatus(id, true)
	}
	if resp != nil {
		resp.Body.Close()
	}
}
