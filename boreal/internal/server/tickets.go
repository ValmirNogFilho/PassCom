package server

import (
	"boreal/internal/dao"
	"boreal/internal/models"
	"encoding/json"
	"net/http"
	"strconv"
)

// handleGetTickets handles HTTP GET requests to retrieve a list of tickets for the authenticated user.
// It checks the request method to ensure it's a GET request and retrieves the user's authorization token from the request headers.
// It then constructs a Request object with the appropriate action and authorization token, and sends it to the
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
	response := GetTickets(
		models.Request{
			Auth: token,
		},
	)

	returnResponse(w, r, response)
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

	response := BuyTicket(models.Request{
		Auth: token,
		Data: buyTicket,
	})

	returnResponse(w, r, response)

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

	queryParams := r.URL.Query()

	id := queryParams.Get("id")

	idUint, err := strconv.ParseUint(id, 10, 32)

	if err != nil {
		returnResponse(w, r, models.Response{
			Error:  err.Error(),
			Status: http.StatusBadRequest,
		})
	}

	response := CancelBuy(uint(idUint), models.Request{
		Auth: token,
	})
	returnResponse(w, r, response)
}

// GetTickets retrieves all tickets associated with the authenticated client.
// It sends a response containing a list of tickets with their respective source, destination, and ID.
//
// Parameters:
//   - auth: A string representing the authentication token.
//   - conn: A net.Conn object representing the connection to the client.
//
// Return:
//   - No return value.
func GetTickets(request models.Request) models.Response {
	session, exists := SessionIfExists(request.Auth)

	if !exists {
		return models.Response{
			Error:  "not authorized",
			Status: http.StatusUnauthorized,
		}
	}
	responseData := make([]map[string]interface{}, 0)

	client, _ := dao.GetClientDAO().FindById(session.ClientID)

	for _, ticket := range client.ClientFlights {
		flight := ticket.Flight

		flightresponse := make(map[string]interface{})
		flightresponse["Src"] = flight.OriginAirport.City
		flightresponse["Dest"] = flight.DestinationAirport.City
		flightresponse["ID"] = ticket.ID
		flightresponse["Company"] = flight.Company
		responseData = append(responseData, flightresponse)
	}

	return models.Response{
		Data: map[string]interface{}{
			"Tickets": responseData,
		},
		Status: http.StatusOK,
	}
}

// BuyTicket handles the process of purchasing a ticket for an authenticated client.
// It checks if the client is authorized, validates the reservation, updates the flight and client data,
// and sends a response indicating success or failure.
//
// Parameters:
//   - auth: A string representing the authentication token.
//   - data: An interface containing the necessary data for purchasing a ticket.
//   - conn: A net.Conn object representing the connection to the client.
//
// Return:
//   - No return value.
func BuyTicket(request models.Request) models.Response {
	session, exists := SessionIfExists(request.Auth)

	if !exists {

		return models.Response{
			Error: "not authorized",
		}
	}

	var buyTicket models.BuyTicket

	jsonData, _ := json.Marshal(request.Data)
	json.Unmarshal(jsonData, &buyTicket)

	flight, _ := dao.GetFlightDAO().FindById(buyTicket.FlightId)

	if flight.Seats > 0 {
		var ticket models.Ticket
		success := false
		id, conn := instance.FindConnectionByName(flight.Company)
		if flight.Company != instance.ServerName && (id != "" && conn.IsOnline) {
			success = instance.initiateBuy(flight.Company, flight.UniqueId)
		} else if flight.Company == instance.ServerName {
			flight.Seats--
			dao.GetFlightDAO().Update(*flight)
			success = true
			instance.broadcast(*flight)
		}

		ticket = models.Ticket{
			ClientId: session.ClientID,
			FlightId: buyTicket.FlightId,
		}

		if success {
			dao.GetTicketDAO().Insert(ticket)
			return models.Response{
				Data: map[string]interface{}{
					"msg": "success",
				},
				Status: http.StatusOK,
			}
		}
	}
	return models.Response{
		Data: map[string]interface{}{
			"Error": "not available seats",
		},
		Status: http.StatusNotAcceptable,
	}

}

// CancelBuy handles the cancellation of a ticket for an authenticated client.
// It checks if the client is authorized, finds the ticket to be canceled, updates the flight and client data,
// and sends a response indicating success or failure.
//
// Parameters:
//   - auth: A string representing the authentication token. This is used to identify the client.
//   - data: An interface containing the necessary data for canceling a ticket.
//   - conn: A net.Conn object representing the connection to the client. This is used to send a response.
//
// Return:
//   - No return value.
func CancelBuy(id uint, request models.Request) models.Response {
	_, exists := SessionIfExists(request.Auth)

	if !exists {
		return models.Response{
			Error:  "not authorized",
			Status: http.StatusUnauthorized,
		}

	}

	ticket, err := dao.GetTicketDAO().FindById(id)

	if err != nil {
		return models.Response{
			Error:  "ticket not found",
			Status: http.StatusNotFound,
		}
	}

	flight := ticket.Flight

	success := false
	connId, conn := instance.FindConnectionByName(flight.Company)
	if flight.Company != instance.ServerName && (connId != "" && conn.IsOnline) {
		success = instance.initiateCancel(flight.Company, flight.UniqueId)
	} else {
		flight.Seats++
		dao.GetFlightDAO().Update(flight)
		success = true
		instance.broadcast(flight)
	}

	if success {
		dao.GetTicketDAO().Delete(*ticket)
	}

	return models.Response{
		Data: map[string]interface{}{
			"msg": "success",
		},
		Status: http.StatusOK,
	}
}
