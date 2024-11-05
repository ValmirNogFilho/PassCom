package server

import (
	"net/http"
	"rumos/internal/models"
)

func handleGetAirports(w http.ResponseWriter, r *http.Request) {
	allowCrossOrigin(w, r)

	if r.Method != http.MethodGet {
		http.Error(w, "only GET allowed", http.StatusMethodNotAllowed)
		return
	}

	token := r.Header.Get("Authorization")
	response := GetAirports(
		models.Request{
			Auth: token,
		},
	)

	returnResponse(w, r, response)
}
