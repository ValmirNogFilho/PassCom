package server

import (
	"bytes"
	"encoding/json"
	"fmt"
	"giro/internal/models"
	"log"
	"net/http"
)

func (s *System) handleConnect(w http.ResponseWriter, r *http.Request) {
	allowCrossOrigin(w, r)

	if r.Method != http.MethodPost {
		http.Error(w, "only POST allowed", http.StatusMethodNotAllowed)
		return
	}

	var message models.Message
	err := json.NewDecoder(r.Body).Decode(&message)
	if err != nil {
		http.Error(w, "Invalid request payload", http.StatusBadRequest)
		return
	}

	// Decodifica Body para map[string]interface{} para acessar Address e Port
	body, ok := message.Body.(map[string]interface{})
	if !ok {
		http.Error(w, "Invalid body format", http.StatusBadRequest)
		return
	}

	name, ok := body["Name"].(string)
	if !ok {
		http.Error(w, "Invalid name format", http.StatusBadRequest)
		return
	}

	address, ok := body["Address"].(string)
	if !ok {
		http.Error(w, "Invalid Address format", http.StatusBadRequest)
		return
	}

	port, ok := body["Port"].(string)
	if !ok {
		http.Error(w, "Invalid Port format", http.StatusBadRequest)
		return
	}

	// Cria uma nova conexão com os dados extraídos
	newConnection := models.Connection{
		Name:     name,
		Address:  address,
		Port:     port,
		IsOnline: true,
	}

	// Adiciona a nova conexão ao sistema
	err = s.AddConnection(message.From, newConnection)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	log.Printf("New connection from %s: %s:%s", message.From, address, port)

	w.WriteHeader(http.StatusCreated)
}

func (s *System) AddConnection(id string, conn models.Connection) error {
	s.Lock.Lock()
	defer s.Lock.Unlock()

	s.Connections[id] = conn
	return nil
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

func (s *System) RequestConnection(name string, address string, port string) {
	// Monta a mensagem de conexão
	message := models.Message{
		From: s.ServerId.String(),
		To:   "", // Destinatário ainda desconhecido
		Body: map[string]interface{}{
			"Name":    s.ServerName,
			"Address": s.Address,
			"Port":    s.Port,
		},
	}

	// Serializa a mensagem em JSON
	jsonData, err := json.Marshal(message)
	if err != nil {
		log.Printf("Error encoding connection request message: %v", err)
		return
	}

	// Cria a URL com endereço e porta do destino
	url := fmt.Sprintf("%s%s:%s/connect", urlPrefix, address, port)
	req, err := http.NewRequest("POST", url, bytes.NewBuffer(jsonData))
	if err != nil {
		log.Printf("Error creating connection request: %v", err)
		return
	}
	req.Header.Set("Content-Type", "application/json")

	// Realiza a solicitação ao destino
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		log.Printf("Error connecting to %s:%s: %v", address, port, err)
		return
	}
	defer resp.Body.Close()

	// Verifica o status da resposta
	if resp.StatusCode != http.StatusCreated {
		log.Printf("Failed to connect to %s:%s - status: %s", address, port, resp.Status)
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

	updatedName, ok := body["Name"].(string)
	if !ok {
		log.Printf("Invalid name format in connection response")
		return
	}

	updatedAddress, ok := body["Address"].(string)
	if !ok {
		log.Printf("Invalid address format in connection response")
		return
	}

	updatedPort, ok := body["Port"].(string)
	if !ok {
		log.Printf("Invalid port format in connection response")
		return
	}

	// Atualiza o destinatário com o ID do servidor recebido na resposta
	newConnection := models.Connection{
		Name:     updatedName,
		Address:  updatedAddress,
		Port:     updatedPort,
		IsOnline: true,
	}

	// Adiciona a nova conexão ao mapa de conexões do sistema
	s.Lock.Lock()
	s.Connections[responseMessage.From] = newConnection
	s.Lock.Unlock()

	log.Printf("Successfully connected to server %s at %s:%s", responseMessage.From, address, port)
}
