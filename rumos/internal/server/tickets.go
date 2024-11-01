package server

import (
	"encoding/json"
	"net/http"
	"vendepass/internal/dao"
	"vendepass/internal/models"
)

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
		flightresponse["Id"] = ticket.ID
		flightresponse["Seats"] = flight.Seats
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
		ticket := models.Ticket{
			ClientId: session.ClientID,
			FlightId: buyTicket.FlightId,
		}

		flight.Seats--
		dao.GetFlightDAO().Update(*flight)
		dao.GetTicketDAO().Insert(ticket)
		return models.Response{
			Data: map[string]interface{}{
				"msg": "success",
			},
			Status: http.StatusOK,
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
func CancelBuy(request models.Request) models.Response {
	_, exists := SessionIfExists(request.Auth)

	if !exists {
		return models.Response{
			Error:  "not authorized",
			Status: http.StatusUnauthorized,
		}

	}

	var cancelBuy models.CancelBuyRequest

	jsonData, _ := json.Marshal(request.Data)
	json.Unmarshal(jsonData, &cancelBuy)

	ticket, err := dao.GetTicketDAO().FindById(cancelBuy.TicketId)

	if err != nil {
		return models.Response{
			Error:  "ticket not found",
			Status: http.StatusNotFound,
		}
	}

	flight := ticket.Flight

	flight.Seats++
	dao.GetFlightDAO().Update(flight)
	dao.GetTicketDAO().Delete(*ticket)

	return models.Response{
		Data: map[string]interface{}{
			"msg": "success",
		},
		Status: http.StatusOK,
	}
}
