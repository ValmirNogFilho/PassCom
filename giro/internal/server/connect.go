package server

import (
	"bytes"
	"encoding/json"
	"giro/internal/models"
	"log"
	"net/http"
)

func (s *System) AddConnection(id string, conn models.Connection) error {
	s.Lock.Lock()
	defer s.Lock.Unlock()

	s.Connections[id] = conn
	return nil
}

func (s *System) removeConnection(id string) {
	s.Lock.Lock()
	defer s.Lock.Unlock()

	delete(s.Connections, id)
}

func (s *System) updateConnectionStatus(id string, isOnline bool) {
	s.Lock.Lock()
	defer s.Lock.Unlock()

	conn, exists := s.Connections[id]
	if exists {
		conn.IsOnline = isOnline
		s.Connections[id] = conn
	}
}

func (s *System) RequestConnection(address string, port string) {
	message, err := models.CreateMessage(s.ServerId.String(), "", s.VectorClock, map[string]interface{}{
		"Name":    s.ServerName,
		"Address": s.ServerName, // TODO: corrigir depois, pois assim está usando o mesmo nome do container
		"Port":    s.Port,
	})

	if err != nil {
		log.Printf("Error creating connection request message: %v", err)
		return
	}

	// Serializa a mensagem em JSON
	jsonData, err := json.Marshal(message)
	if err != nil {
		log.Printf("Error encoding connection request message: %v", err)
		return
	}

	url := URL_PREFIX + address + ":" + port

	log.Print("url being used is: ", url)
	// Cria a URL com endereço e porta do destino
	req, err := http.NewRequest("POST", url+"/server/connect", bytes.NewBuffer(jsonData))
	if err != nil {
		log.Printf("Error creating connection request: %v", err)
		return
	}
	req.Header.Set("Content-Type", "application/json")

	// Realiza a solicitação ao destino
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		log.Printf("Error connecting to %s: %v", url, err)
		return
	}
	defer resp.Body.Close()

	// Verifica o status da resposta
	if resp.StatusCode != http.StatusCreated {
		log.Printf("Failed to connect to %s - status: %s", url, resp.Status)
		return
	}

	// Lê a resposta e extrai o ID do servidor conectado
	var responseMessage models.Message
	if err := json.NewDecoder(resp.Body).Decode(&responseMessage); err != nil {
		log.Printf("Error decoding connection response: %v", err)
		return
	}

	// Extrai o corpo da resposta e valida os campos
	body, ok := responseMessage.Body.(map[string]interface{})
	if !ok {
		log.Printf("Invalid body format in connection response")
		return
	}

	name, ok := body["Name"].(string)
	if !ok {
		log.Printf("Invalid name format in connection response")
		return
	}

	// Atualiza o destinatário com o ID do servidor recebido na resposta
	newConnection := models.Connection{
		Name:     name,
		Address:  address,
		Port:     port,
		IsOnline: true,
	}

	// Adiciona a nova conexão ao mapa de conexões do sistema
	s.Lock.Lock()
	s.Connections[responseMessage.From] = newConnection
	s.Lock.Unlock()

	log.Printf("Successfully connected to server %s at %s", name, url)
}
