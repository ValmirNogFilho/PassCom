package server

import (
	"encoding/json"
	"net"
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
func GetTickets(auth string, conn net.Conn) {
	session, exists := SessionIfExists(auth)

	if !exists {
		WriteNewResponse(models.Response{
			Error: "not authorized",
		}, conn)
		return
	}
	responseData := make([]map[string]interface{}, 0)

	client, _ := dao.GetClientDAO().FindById(session.ClientID)

	for _, ticket := range client.ClientFlights {
		flight := ticket.Flight

		flightresponse := make(map[string]interface{})

		flightresponse["Src"] = flight.OriginAirport.City
		flightresponse["Dest"] = flight.DestinationAirport.City
		flightresponse["Id"] = ticket.ID
		responseData = append(responseData, flightresponse)
	}

	WriteNewResponse(models.Response{
		Data: map[string]interface{}{
			"Tickets": responseData,
		},
	}, conn)
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
func BuyTicket(auth string, data interface{}, conn net.Conn) {
	session, exists := SessionIfExists(auth)

	if !exists {
		WriteNewResponse(models.Response{
			Error: "not authorized",
		}, conn)
		return
	}

	var buyTicket models.BuyTicket

	jsonData, _ := json.Marshal(data)
	json.Unmarshal(jsonData, &buyTicket)

	ticket := models.Ticket{
		ClientId: session.ClientID,
		FlightId: buyTicket.FlightId,
	}

	dao.GetTicketDAO().Insert(ticket)

	WriteNewResponse(models.Response{
		Data: map[string]interface{}{
			"msg": "success",
		},
	}, conn)
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
func CancelBuy(auth string, data interface{}, conn net.Conn) {
	_, exists := SessionIfExists(auth)

	if !exists {
		WriteNewResponse(models.Response{
			Error: "not authorized",
		}, conn)
		return
	}

	var cancelBuy models.CancelBuyRequest

	jsonData, _ := json.Marshal(data)
	json.Unmarshal(jsonData, &cancelBuy)

	ticket, err := dao.GetTicketDAO().FindById(cancelBuy.TicketId)

	if err != nil {
		WriteNewResponse(models.Response{
			Error: "ticket not found",
		}, conn)
		return
	}

	flight := ticket.Flight

	flight.Seats++
	dao.GetFlightDAO().Update(flight)
	dao.GetTicketDAO().Delete(*ticket)

	WriteNewResponse(models.Response{
		Data: map[string]interface{}{
			"msg": "success",
		},
	}, conn)
}
