package server

import (
	"net/http"
)

func (s *System) handleFlight(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	// Receber o Id do aeroporto de origem e aeroporto de destino
	// Procurar no DAO o Id do voo a ser
}
