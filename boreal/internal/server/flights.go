package server

import (
	"boreal/internal/dao"
	"boreal/internal/models"
	"encoding/json"
	"net/http"
)

// handleGetFlights is an HTTP handler function that retrieves flight information based on the provided flight IDs.
// It checks the HTTP method of the request to ensure it's a POST request.
// If the method is not POST, it returns a 405 Method Not Allowed status with an error message.
// It extracts the user's authorization token from the request headers and decodes the request body into a FlightsRequest struct.
// If the decoding fails, it returns a 400 Bad Request status.
// It then constructs a Request object with the appropriate action, authorization token, and flight IDs,
// and sends it to the server using the writeAndReturnResponse function.
//
// Parameters:
//   - w: http.ResponseWriter to write the HTTP response.
//   - r: *http.Request to read the HTTP request.
func handleGetFlights(w http.ResponseWriter, r *http.Request) {
	allowCrossOrigin(w, r)

	if r.Method != http.MethodPost {
		http.Error(w, "only POST allowed", http.StatusMethodNotAllowed)
		return
	}

	token := r.Header.Get("Authorization")

	var flightIds models.FlightsRequest
	err := json.NewDecoder(r.Body).Decode(&flightIds)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	response := Flights(models.Request{
		Auth: token,
		Data: flightIds,
	})

	returnResponse(w, r, response)
}

// handleGetRoute is an HTTP handler function that retrieves route information based on the provided source and destination.
// It checks the HTTP method of the request to ensure it's a GET request.
// If the method is not GET, it returns a 405 Method Not Allowed status with an error message.
// It extracts the source and destination from the request query parameters and the user's authorization token from the request headers.
// It then constructs a Request object with the appropriate action, authorization token, and route request data,
// and sends it to the server using the writeAndReturnResponse function.
//
// Parameters:
//   - w: http.ResponseWriter to write the HTTP response.
//   - r: *http.Request to read the HTTP request.
func handleGetRoute(w http.ResponseWriter, r *http.Request) {
	allowCrossOrigin(w, r)
	if r.Method != http.MethodGet {
		http.Error(w, "only GET allowed", http.StatusMethodNotAllowed)
		return
	}
	queryParams := r.URL.Query()

	src := queryParams.Get("src")
	dest := queryParams.Get("dest")

	token := r.Header.Get("Authorization")
	response := Route(models.Request{
		Auth: token,
		Data: models.RouteRequest{
			Source: src,
			Dest:   dest,
		}})
	returnResponse(w, r, response)
}

func AddFlights(flights []models.Flight) {
	for _, flight := range flights {
		flight.ID = 0
		dao.GetFlightDAO().Insert(flight)
	}
}

func RemoveFlights(company string) {
	dao.GetFlightDAO().DeleteByCompany(company)
}
