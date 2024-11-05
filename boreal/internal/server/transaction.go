package server

import (
	"boreal/internal/dao"
	"boreal/internal/models"
	"boreal/internal/utils"
	"bytes"
	"encoding/json"
	"log"
	"net/http"
)

func (s *System) HandleServerTicketPurchase(w http.ResponseWriter, r *http.Request) {
	allowCrossOrigin(w, r)

	s.Lock.Lock()
	defer s.Lock.Unlock()

	s.IncrementClock()

	if r.Method != http.MethodPost {
		http.Error(w, "Only POST allowed", http.StatusMethodNotAllowed)
		return
	}

	var msg models.Message
	err := json.NewDecoder(r.Body).Decode(&msg)

	if err != nil {
		http.Error(w, "Invalid request payload", http.StatusBadRequest)
		return
	}

	to := msg.To
	body, ok := msg.Body.(string)

	if !ok || body == "" {
		http.Error(w, "Invalid purchase data", http.StatusBadRequest)
		return
	}

	flight, err := dao.GetFlightDAO().FindByUniqueId(body)

	if err != nil {
		http.Error(w, "Flight not found", http.StatusNotFound)
		return
	}

	flight.Seats--
	dao.GetFlightDAO().Update(*flight)

	responseMsg, err := models.CreateMessage(s.ServerId.String(), to, s.VectorClock, "")
	if err != nil {
		http.Error(w, "Failed to create response message", http.StatusInternalServerError)
		return
	}
	utils.SendJSONResponse(w, responseMsg, http.StatusOK)

	s.broadcast(*flight)
}

func (s *System) HandleServerTicketCancel(w http.ResponseWriter, r *http.Request) {
	allowCrossOrigin(w, r)

	s.Lock.Lock()
	defer s.Lock.Unlock()

	s.IncrementClock()

	if r.Method != http.MethodDelete {
		http.Error(w, "Only DELETE allowed", http.StatusMethodNotAllowed)
		return
	}

	var msg models.Message
	err := json.NewDecoder(r.Body).Decode(&msg)

	if err != nil {
		http.Error(w, "Invalid request payload", http.StatusBadRequest)
		return
	}

	to := msg.To
	body, ok := msg.Body.(string)

	if !ok || body == "" {
		http.Error(w, "Invalid cancellation data", http.StatusBadRequest)
		return
	}

	flight, err := dao.GetFlightDAO().FindByUniqueId(body)
	if err != nil {
		http.Error(w, "Flight not found", http.StatusNotFound)
		return
	}

	flight.Seats++
	dao.GetFlightDAO().Update(*flight)

	responseMsg, err := models.CreateMessage(s.ServerId.String(), to, s.VectorClock, "")
	if err != nil {
		http.Error(w, "Failed to create response message", http.StatusInternalServerError)
		return
	}

	utils.SendJSONResponse(w, responseMsg, http.StatusOK)

	s.broadcast(*flight)
}

func (s *System) initiateBuy(company, uniqueId string) bool {
	id, conn := s.FindConnectionByName(company)
	if id == "" {
		log.Printf("Connection not found for company: %s", company)
		return false
	}

	url := URL_PREFIX + conn.Address + ":" + conn.Port + "/server/ticket/purchase"

	// Cria a mensagem de compra com UniqueId do voo
	requestMsg, err := models.CreateMessage(s.ServerId.String(), company, s.VectorClock, uniqueId)
	if err != nil {
		log.Printf("Error creating request message for purchase: %v", err)
		return false
	}

	// Converte a mensagem para JSON
	jsonData, err := json.Marshal(requestMsg)
	if err != nil {
		log.Printf("Error encoding purchase message to JSON: %v", err)
		return false
	}

	// Envia a solicitação de compra ao servidor da companhia
	resp, err := http.Post(url, "application/json", bytes.NewBuffer(jsonData))
	if err != nil {
		log.Printf("Error sending purchase request: %v", err)
		return false
	}
	defer resp.Body.Close()

	// Verifica o status da resposta
	if resp.StatusCode != http.StatusOK {
		log.Printf("Purchase request failed with status: %s", resp.Status)
		return false
	}

	log.Printf("Purchase request successful for flight %s on company %s", uniqueId, company)
	return true
}

func (s *System) initiateCancel(company, uniqueId string) bool {
	// Localiza o endereço do servidor da companhia responsável
	id, conn := s.FindConnectionByName(company)
	if id == "" {
		log.Printf("Connection not found for company: %s", company)
		return false
	}

	url := URL_PREFIX + conn.Address + ":" + conn.Port + "/server/ticket/cancel"

	// Cria a mensagem de cancelamento com UniqueId do voo
	requestMsg, err := models.CreateMessage(s.ServerId.String(), id, s.VectorClock, uniqueId)
	if err != nil {
		log.Printf("Error creating request message for cancellation: %v", err)
		return false
	}

	// Converte a mensagem para JSON
	jsonData, err := json.Marshal(requestMsg)
	if err != nil {
		log.Printf("Error encoding cancellation message to JSON: %v", err)
		return false
	}

	// Cria uma requisição HTTP DELETE
	req, err := http.NewRequest(http.MethodDelete, url, bytes.NewBuffer(jsonData))
	if err != nil {
		log.Printf("Error creating cancel request: %v", err)
		return false
	}
	req.Header.Set("Content-Type", "application/json")

	// Envia a solicitação de cancelamento ao servidor da companhia
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		log.Printf("Error sending cancellation request: %v", err)
		return false
	}
	defer resp.Body.Close()

	// Verifica o status da resposta
	if resp.StatusCode != http.StatusOK {
		log.Printf("Cancellation request failed with status: %s", resp.Status)
		return false
	}

	log.Printf("Cancellation request successful for flight %s on company %s", uniqueId, company)
	return true
}
