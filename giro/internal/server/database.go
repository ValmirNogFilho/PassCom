package server

import (
	"encoding/json"
	"giro/internal/dao"
	"giro/internal/models"
	"log"
	"net/http"
)

func (s *System) RequestDatabase(address string, port string) {
	s.Lock.Lock()
	defer s.Lock.Unlock()

	url := URL_PREFIX + address + ":" + port + "/server/database"

	req, err := http.NewRequest("GET", url, nil)
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

	// Decodifica a resposta JSON contendo os dados do banco de dados
	var flights []models.Flight
	if err := json.NewDecoder(resp.Body).Decode(&flights); err != nil {
		log.Printf("Error decoding database response: %v", err)
		return
	}

	// Insere ou atualiza cada registro de voo recebido no banco de dados local
	for _, flight := range flights {
		flight.ID = 0
		dao.GetFlightDAO().Insert(flight)
	}
}

func (s *System) RemoveDatabase(company string) {
	s.Lock.Lock()
	defer s.Lock.Unlock()

	dao.GetFlightDAO().DeleteByCompany(company)
}
