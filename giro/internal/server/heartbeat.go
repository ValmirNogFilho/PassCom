package server

import (
	"bytes"
	"encoding/json"
	"fmt"
	"giro/internal/models"
	"log"
	"net/http"
	"time"
)

// Envia heartbeats constantemente para todos os servidores no mapa Connections
func (s *System) sendHeartbeats() {
	ticker := time.NewTicker(5 * time.Second) // Intervalo de envio do heartbeat
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			s.Lock.RLock()
			for id, conn := range s.Connections {
				s.wg.Add(1)
				go s.sendHeartbeatToConnection(id, conn)
			}
			s.Lock.RUnlock()
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
