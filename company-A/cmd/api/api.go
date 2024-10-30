package main

import (
	"encoding/json"
	"log"
	"net"
	"net/http"
	"vendepass/internal/models"
)

const (
	port      = ":9999"
	CONN_PORT = "8888"
	CONN_HOST = "app"
	CONN_TYPE = "tcp"
)

// main initializes and starts the HTTP server for the application.
// It sets up the routes for various API endpoints and listens on the specified port.
func main() {
	http.HandleFunc("/login", handleLogin)
	http.HandleFunc("/logout", handleLogout)
	http.HandleFunc("/user", handleGetUser)
	http.HandleFunc("/route", handleGetRoute)
	http.HandleFunc("/flights", handleGetFlights)
	http.HandleFunc("/reservation", handleReservation)
	http.HandleFunc("/cart", handleGetCart)
	http.HandleFunc("/ticket", handleTicket)
	http.HandleFunc("/tickets", handleGetTickets)
	log.Fatal(http.ListenAndServe(port, nil))
}

// allowCrossOrigin is a middleware function that handles Cross-Origin Resource Sharing (CORS)
// for HTTP requests. It sets the necessary headers to allow cross-origin requests and
// handles the preflight OPTIONS request.
//
// Parameters:
//   - w: http.ResponseWriter to write the response headers.
//   - r: *http.Request to read the request method.
func allowCrossOrigin(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")

	if r.Method == http.MethodOptions {
		w.WriteHeader(http.StatusNoContent)
		return
	}
}

// handleGetTickets handles HTTP GET requests to retrieve a list of tickets for the authenticated user.
// It checks the request method to ensure it's a GET request and retrieves the user's authorization token from the request headers.
// It then constructs a Request object with the appropriate action and authorization token, and sends it to the server.
// The server's response is then decoded and returned as a JSON object in the HTTP response.
//
// Parameters:
//   - w: http.ResponseWriter to write the HTTP response.
//   - r: *http.Request to read the HTTP request.
func handleGetTickets(w http.ResponseWriter, r *http.Request) {
	allowCrossOrigin(w, r)

	if r.Method != http.MethodGet {
		http.Error(w, "only GET allowed", http.StatusMethodNotAllowed)
		return
	}

	token := r.Header.Get("Authorization")

	writeAndReturnResponse(w, models.Request{
		Action: "tickets",
		Auth:   token,
	})
}

// handleTicket is a HTTP handler function that handles requests for buying and canceling tickets.
// It checks the HTTP method of the request and calls the appropriate handler function based on the method.
// If the method is neither POST nor DELETE, it returns a 405 Method Not Allowed status with an error message.
//
// Parameters:
//   - w: http.ResponseWriter to write the HTTP response.
//   - r: *http.Request to read the HTTP request.
func handleTicket(w http.ResponseWriter, r *http.Request) {
	allowCrossOrigin(w, r)
	switch r.Method {
	case http.MethodPost:
		handleBuyTicket(w, r)
	case http.MethodDelete:
		handleCancelTicket(w, r)
	default:
		http.Error(w, "only POST or DELETE allowed", http.StatusMethodNotAllowed)
		return
	}
}

// handleBuyTicket is a HTTP handler function that handles requests for buying tickets.
// It extracts the user's authorization token from the request headers and decodes the request body into a BuyTicket struct.
// If the decoding fails, it returns a 400 Bad Request status.
// It then constructs a Request object with the appropriate action, authorization token, and buy ticket data,
// and sends it to the server using the writeAndReturnResponse function.
//
// Parameters:
//   - w: http.ResponseWriter to write the HTTP response.
//   - r: *http.Request to read the HTTP request.
func handleBuyTicket(w http.ResponseWriter, r *http.Request) {
	token := r.Header.Get("Authorization")
	var buyTicket models.BuyTicket

	err := json.NewDecoder(r.Body).Decode(&buyTicket)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	writeAndReturnResponse(w, models.Request{
		Action: "buy",
		Auth:   token,
		Data:   buyTicket,
	})
}

// handleCancelTicket is a HTTP handler function that handles requests for canceling tickets.
// It extracts the user's authorization token from the request headers and decodes the request body into a CancelBuyRequest struct.
// If the decoding fails, it returns a 400 Bad Request status.
// It then constructs a Request object with the appropriate action, authorization token, and cancel ticket data,
// and sends it to the server using the writeAndReturnResponse function.
//
// Parameters:
//   - w: http.ResponseWriter to write the HTTP response.
//   - r: *http.Request to read the HTTP request.
func handleCancelTicket(w http.ResponseWriter, r *http.Request) {

	token := r.Header.Get("Authorization")

	var ticketId models.CancelBuyRequest

	err := json.NewDecoder(r.Body).Decode(&ticketId)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	writeAndReturnResponse(w, models.Request{
		Action: "cancel-buy",
		Auth:   token,
		Data:   ticketId,
	})
}

// handleGetCart is an HTTP handler function that retrieves the user's shopping cart.
// It checks the HTTP method of the request to ensure it's a GET request.
// If the method is not GET, it returns a 405 Method Not Allowed status with an error message.
// It extracts the user's authorization token from the request headers and constructs a Request object
// with the appropriate action and authorization token.
// The constructed Request object is then sent to the server using the writeAndReturnResponse function.
//
// Parameters:
//   - w: http.ResponseWriter to write the HTTP response.
//   - r: *http.Request to read the HTTP request.
func handleGetCart(w http.ResponseWriter, r *http.Request) {
	allowCrossOrigin(w, r)

	if r.Method != http.MethodGet {
		http.Error(w, "only GET allowed", http.StatusMethodNotAllowed)
		return
	}

	token := r.Header.Get("Authorization")

	writeAndReturnResponse(w, models.Request{
		Action: "cart",
		Auth:   token,
	})
}

// handleReservation is an HTTP handler function that handles requests for making and canceling reservations.
// It checks the HTTP method of the request and calls the appropriate handler function based on the method.
// If the method is neither POST nor DELETE, it returns a 405 Method Not Allowed status with an error message.
//
// Parameters:
//   - w: http.ResponseWriter to write the HTTP response.
//   - r: *http.Request to read the HTTP request.
func handleReservation(w http.ResponseWriter, r *http.Request) {
	allowCrossOrigin(w, r)
	switch r.Method {
	case http.MethodPost:
		handleMakeReservations(w, r)
	case http.MethodDelete:
		handleCancelReservation(w, r)
	default:
		http.Error(w, "only POST or DELETE allowed", http.StatusMethodNotAllowed)
		return
	}
}

// handleMakeReservations is an HTTP handler function that processes requests for making reservations.
// It extracts the user's authorization token from the request headers and decodes the request body into a FlightsRequest struct.
// If the decoding fails, it returns a 400 Bad Request status.
// It then constructs a Request object with the appropriate action, authorization token, and flight IDs,
// and sends it to the server using the writeAndReturnResponse function.
//
// Parameters:
//   - w: http.ResponseWriter to write the HTTP response.
//   - r: *http.Request to read the HTTP request.
func handleMakeReservations(w http.ResponseWriter, r *http.Request) {
	// Extract the authorization token from the request headers
	token := r.Header.Get("Authorization")

	var flightIds models.FlightsRequest

	err := json.NewDecoder(r.Body).Decode(&flightIds)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	request := models.Request{
		Action: "reservation",
		Auth:   token,
		Data:   flightIds,
	}

	writeAndReturnResponse(w, request)
}

// handleCancelReservation is an HTTP handler function that processes requests for canceling reservations.
// It extracts the user's authorization token from the request headers and decodes the request body into a CancelReservationRequest struct.
// If the decoding fails, it returns a 400 Bad Request status.
// It then constructs a Request object with the appropriate action, authorization token, and reservation ID,
// and sends it to the server using the writeAndReturnResponse function.
//
// Parameters:
//   - w: http.ResponseWriter to write the HTTP response.
//   - r: *http.Request to read the HTTP request.
func handleCancelReservation(w http.ResponseWriter, r *http.Request) {

	token := r.Header.Get("Authorization")

	var reservationId models.CancelReservationRequest

	err := json.NewDecoder(r.Body).Decode(&reservationId)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	writeAndReturnResponse(w, models.Request{
		Action: "cancel-reservation",
		Auth:   token,
		Data:   reservationId,
	})

}

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
		http.Error(w, "only GET allowed", http.StatusMethodNotAllowed)
		return
	}

	token := r.Header.Get("Authorization")

	var flightIds models.FlightsRequest

	err := json.NewDecoder(r.Body).Decode(&flightIds)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	writeAndReturnResponse(w, models.Request{
		Action: "flights",
		Auth:   token,
		Data:   flightIds,
	})
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
	writeAndReturnResponse(w, models.Request{
		Action: "route",
		Auth:   token,
		Data: models.RouteRequest{
			Source: src,
			Dest:   dest,
		},
	})
}

// handleGetUser is an HTTP handler function that retrieves user information.
// It checks the HTTP method of the request to ensure it's a GET request.
// If the method is not GET, it returns a 405 Method Not Allowed status with an error message.
// It extracts the user's authorization token from the request headers and constructs a Request object
// with the appropriate action and authorization token.
// The constructed Request object is then sent to the server using the writeAndReturnResponse function.
//
// Parameters:
//   - w: http.ResponseWriter to write the HTTP response.
//   - r: *http.Request to read the HTTP request.
func handleGetUser(w http.ResponseWriter, r *http.Request) {
	allowCrossOrigin(w, r)
	if r.Method != http.MethodGet {
		http.Error(w, "only GET allowed", http.StatusMethodNotAllowed)
		return
	}

	token := r.Header.Get("Authorization")

	writeAndReturnResponse(w, models.Request{Action: "get-user", Auth: token})
}

// handleLogout handles HTTP GET requests to log out the authenticated user.
// It checks the request method to ensure it's a GET request and retrieves the user's authorization token from the request headers.
// It then constructs a Request object with the appropriate action and authorization token, and sends it to the server.
// The server's response is then returned as a JSON object in the HTTP response.
//
// Parameters:
//   - w: http.ResponseWriter to write the HTTP response.
//   - r: *http.Request to read the HTTP request.
func handleLogout(w http.ResponseWriter, r *http.Request) {
	allowCrossOrigin(w, r)
	if r.Method != http.MethodGet {
		http.Error(w, "only GET allowed", http.StatusMethodNotAllowed)
		return
	}

	token := r.Header.Get("Authorization")

	writeAndReturnResponse(w, models.Request{Action: "logout", Auth: token})
}

// handleLogin handles HTTP POST requests to log in the authenticated user.
// It checks the request method to ensure it's a POST request and retrieves the user's login credentials from the request body.
// If the method is not POST, it returns a 405 Method Not Allowed status with an error message.
// If the decoding of the login credentials fails, it returns a 400 Bad Request status.
// It then constructs a Request object with the appropriate action and login credentials, and sends it to the server using the writeAndReturnResponse function.
//
// Parameters:
//   - w: http.ResponseWriter to write the HTTP response.
//   - r: *http.Request to read the HTTP request.
func handleLogin(w http.ResponseWriter, r *http.Request) {
	allowCrossOrigin(w, r)

	if r.Method != http.MethodPost {
		http.Error(w, "only POST allowed", http.StatusMethodNotAllowed)
		return
	}

	var logCred models.LoginCredentials
	err := json.NewDecoder(r.Body).Decode(&logCred)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	writeAndReturnResponse(w, models.Request{
		Action: "login",
		Data:   logCred,
	})
}

// writeAndReturnResponse is a function that establishes a connection to the server,
// sends a request to the server, receives the server's response, and writes the response to the HTTP response writer.
//
// Parameters:
// - w: http.ResponseWriter to write the HTTP response.
// - req: models.Request representing the request to be sent to the server.
func writeAndReturnResponse(w http.ResponseWriter, req models.Request) {
	conn, err := net.Dial(CONN_TYPE, CONN_HOST+":"+CONN_PORT)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer conn.Close()

	buffer, _ := json.Marshal(req)
	_, writeErr := conn.Write(buffer)
	if writeErr != nil {
		http.Error(w, writeErr.Error(), http.StatusInternalServerError)
		return
	}

	receive := make([]byte, 2048)
	n, readErr := conn.Read(receive)
	if readErr != nil {
		http.Error(w, readErr.Error(), http.StatusInternalServerError)
		return
	}

	var responseData models.Response
	err = json.Unmarshal(receive[:n], &responseData)
	if err != nil {
		http.Error(w, "Failed to decode response from server", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(responseData)
}
