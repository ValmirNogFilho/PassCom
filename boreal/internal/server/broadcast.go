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

func (s *System) HandleBroadcast(w http.ResponseWriter, r *http.Request) {
	allowCrossOrigin(w, r)
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

	jsonData, err := json.Marshal(msg.Body)
	if err != nil {
		http.Error(w, "Error marshalling data: "+err.Error(), http.StatusInternalServerError)
		return
	}

	var flight models.Flight
	err = json.Unmarshal(jsonData, &flight)
	if err != nil {
		http.Error(w, "Error unmarshalling data: "+err.Error(), http.StatusInternalServerError)
		return
	}

	prevFlight, err := dao.GetFlightDAO().FindByUniqueId(flight.UniqueId)
	if err != nil {
		http.Error(w, "Flight not found", http.StatusNotFound)
		return
	}

	prevFlight.Seats = flight.Seats
	dao.GetFlightDAO().Update(*prevFlight)

	responseMsg, err := models.CreateMessage(s.ServerId.String(), to, s.VectorClock, "")
	if err != nil {
		http.Error(w, "Failed to create response message", http.StatusInternalServerError)
		return
	}

	utils.SendJSONResponse(w, responseMsg, http.StatusOK)
}

func (s *System) broadcast(flight models.Flight) {
	s.IncrementClock()

	for id, conn := range s.Connections {
		// Cria a mensagem para cada conexão
		newMsg, err := models.CreateMessage(s.ServerId.String(), id, s.VectorClock, flight)
		if err != nil {
			log.Printf("Error creating message for flight %s: %v", flight.UniqueId, err)
			continue
		}
		url := URL_PREFIX + conn.Address + ":" + conn.Port + "/server/broadcast"

		// Adiciona uma nova goroutine ao WaitGroup para envio assíncrono
		s.wg.Add(1)
		go s.sendFlight(url, flight, *newMsg)
	}

	// Aguarda o término de todas as goroutines de envio
	s.wg.Wait()
}

func (s *System) sendFlight(url string, flight models.Flight, message models.Message) {
	defer s.wg.Done()

	// Serializa a mensagem para JSON
	jsonData, err := json.Marshal(message)
	if err != nil {
		log.Printf("Error marshalling message for flight %s: %v", flight.UniqueId, err)
		return
	}

	// Envia a requisição HTTP POST ao servidor de destino
	resp, err := http.Post(url, "application/json", bytes.NewBuffer(jsonData))
	if err != nil {
		log.Printf("Error sending flight %s to %s: %v", flight.UniqueId, url, err)
		return
	}
	defer resp.Body.Close()

	// Verifica o status da resposta
	if resp.StatusCode != http.StatusOK {
		log.Printf("Failed to broadcast flight %s to %s, status: %s", flight.UniqueId, url, resp.Status)
	} else {
		log.Printf("Successfully broadcasted flight %s to %s", flight.UniqueId, url)
	}
}
