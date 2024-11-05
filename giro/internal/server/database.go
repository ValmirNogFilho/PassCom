package server

import (
	"bytes"
	"encoding/json"
	"giro/internal/dao"
	"giro/internal/models"
	"giro/internal/utils"
	"log"
	"net/http"
)

func (s *System) handleDatabase(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	switch r.Method {
	case http.MethodGet:
		s.handleGetDatabase(w, r)
	case http.MethodPut:
		s.HandlePutDatabase(w, r)
	case http.MethodDelete:
		s.HandleDeleteDatabase(w, r)
	default:
		http.Error(w, "Unsupported method", http.StatusMethodNotAllowed)
		return
	}
}

func (s *System) handleGetDatabase(w http.ResponseWriter, r *http.Request) {
	s.Lock.RLock()
	defer s.Lock.RUnlock()

	var msg models.Message
	err := json.NewDecoder(r.Body).Decode(&msg)

	if err != nil {
		http.Error(w, "Invalid request payload", http.StatusBadRequest)
		return
	}

	to := msg.To

	db := dao.GetFlightDAO()

	flights, err := db.FindByCompany(s.ServerName)
	if err != nil {
		log.Printf("Error searching flights: %v", err)
		http.Error(w, "Failed to find flights", http.StatusInternalServerError)
		return
	}

	responseMsg, err := models.CreateMessage(s.ServerName, to, s.VectorClock, flights)

	if err != nil {
		log.Printf("Error creating response message: %v", err)
		http.Error(w, "Failed to create response message", http.StatusInternalServerError)
		return
	}

	utils.SendJSONResponse(w, responseMsg, http.StatusOK)
}

func (s *System) HandlePutDatabase(w http.ResponseWriter, r *http.Request) {
	s.Lock.Lock()
	defer s.Lock.Unlock()

	var msg models.Message
	err := json.NewDecoder(r.Body).Decode(&msg)

	if err != nil {
		http.Error(w, "Invalid request payload", http.StatusBadRequest)
		return
	}

	to := msg.To

	jsonFlights, err := json.Marshal(msg.Body)
	if err != nil {
		log.Printf("Error marshalling flight data: %v", err)
		return
	}

	var flights []models.Flight
	err = json.Unmarshal(jsonFlights, &flights)
	if err != nil {
		log.Printf("Error unmarshalling flight data: %v", err)
		return
	}

	AddFlights(flights)

	responseMsg, err := models.CreateMessage(s.ServerName, to, s.VectorClock, "Received database")

	if err != nil {
		log.Printf("Error creating response message: %v", err)
		http.Error(w, "Failed to create response message", http.StatusInternalServerError)
		return
	}

	utils.SendJSONResponse(w, responseMsg, http.StatusOK)
}

func (s *System) HandleDeleteDatabase(w http.ResponseWriter, r *http.Request) {
	s.Lock.Lock()
	defer s.Lock.Unlock()

	var msg models.Message
	err := json.NewDecoder(r.Body).Decode(&msg)

	if err != nil {
		http.Error(w, "Invalid request payload", http.StatusBadRequest)
		return
	}

	to := msg.To
	RemoveFlights(s.Connections[to].Name)
	responseMsg, err := models.CreateMessage(s.ServerName, to, s.VectorClock, "Database deleted")

	if err != nil {
		log.Printf("Error creating response message: %v", err)
		http.Error(w, "Failed to create response message", http.StatusInternalServerError)
		return
	}

	utils.SendJSONResponse(w, responseMsg, http.StatusOK)
}

func (s *System) RequestDatabase(id string, address string, port string) {
	s.Lock.Lock()
	defer s.Lock.Unlock()

	url := URL_PREFIX + address + ":" + port + "/server/database"

	requestMsg, err := models.CreateMessage(s.ServerName, id, s.VectorClock, "")

	if err != nil {
		log.Printf("Error creating database request message: %v", err)
		return
	}
	jsonData, err := json.Marshal(requestMsg)

	if err != nil {
		log.Printf("Error encoding database request message: %v", err)
		return
	}

	req, err := http.NewRequest("GET", url, bytes.NewBuffer(jsonData))
	if err != nil {
		log.Printf("Error creating database request: %v", err)
		return
	}

	// Envia a solicitação ao servidor remoto
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		log.Printf("Error requesting database from %s", url)
		return
	}
	defer resp.Body.Close()

	// Verifica o status da resposta
	if resp.StatusCode != http.StatusOK {
		log.Printf("Failed to retrieve database from %s", url)
		return
	}

	var msg models.Message
	err = json.NewDecoder(resp.Body).Decode(&msg)

	if err != nil {
		log.Printf("Error decoding database response: %v", err)
		return
	}

	jsonFlights, err := json.Marshal(msg.Body)
	if err != nil {
		log.Printf("Error marshalling flight data: %v", err)
		return
	}

	var flights []models.Flight
	err = json.Unmarshal(jsonFlights, &flights)
	if err != nil {
		log.Printf("Error unmarshalling flight data: %v", err)
		return
	}

	// Insere ou atualiza cada registro de voo recebido no banco de dados local
	AddFlights(flights)
}

func (s *System) SendDatabase(id string, address string, port string) {
	s.Lock.Lock()
	defer s.Lock.Unlock()

	url := URL_PREFIX + address + ":" + port + "/server/database"

	// Obtém os voos da companhia atual
	flights, err := dao.GetFlightDAO().FindByCompany(s.ServerName)
	if err != nil {
		log.Printf("Error retrieving flights from database: %v", err)
		return
	}

	requestMsg, err := models.CreateMessage(s.ServerName, id, s.VectorClock, flights)

	if err != nil {
		log.Printf("Error creating database request message: %v", err)
		return
	}
	jsonData, err := json.Marshal(requestMsg)
	if err != nil {
		log.Printf("Error encoding flights to JSON: %v", err)
		return
	}

	// Cria uma requisição HTTP PUT com o payload JSON
	req, err := http.NewRequest("PUT", url, bytes.NewBuffer(jsonData))
	if err != nil {
		log.Printf("Error creating database send request: %v", err)
		return
	}
	req.Header.Set("Content-Type", "application/json")

	// Envia a requisição para o servidor de destino
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		log.Printf("Error sending database to %s: %v", url, err)
		return
	}
	defer resp.Body.Close()

	// Verifica o status da resposta para garantir que a operação foi bem-sucedida
	if resp.StatusCode != http.StatusOK {
		log.Printf("Failed to send database to %s - status: %s", url, resp.Status)
		return
	}

	log.Printf("Database successfully sent to server at %s", url)
}

func (s *System) RemoveDatabase(company string) {
	s.Lock.Lock()
	defer s.Lock.Unlock()

	RemoveFlights(company)
}

func (s *System) RequestDatabaseRemoval(id string, address string, port string) {
	s.Lock.Lock()
	defer s.Lock.Unlock()

	url := URL_PREFIX + address + ":" + port + "/server/database"

	requestMsg, err := models.CreateMessage(s.ServerName, id, s.VectorClock, "")
	if err != nil {
		log.Printf("Error creating database request message: %v", err)
		return
	}

	jsonData, err := json.Marshal(requestMsg)
	if err != nil {
		log.Printf("Error encoding flights to JSON: %v", err)
		return
	}

	req, err := http.NewRequest("DELETE", url, bytes.NewBuffer(jsonData))
	if err != nil {
		log.Printf("Error creating database send request: %v", err)
		return
	}

	// Envia a solicitação ao servidor remoto
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		log.Printf("Error requesting database removal from %s", url)
		return
	}
	defer resp.Body.Close()

	// Verifica o status da resposta
	if resp.StatusCode != http.StatusOK {
		log.Printf("Failed to remove database from %s", url)
		return
	}
}
