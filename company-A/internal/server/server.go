package server

import (
	"encoding/json"
	"fmt"
	"net"
	"time"
	"vendepass/internal/dao"
	"vendepass/internal/models"

	"github.com/google/uuid"
)

// HandleConn manages a single client connection. It reads incoming requests, processes them,
// and sends responses back to the client. If an error occurs during the handling of a request,
// an appropriate error response is sent back to the client.
//
// Parameters:
//   - conn: The net.Conn object representing the client connection.
func HandleConn(conn net.Conn) {
	defer conn.Close()

	buffer := make([]byte, 2048)
	n, err := conn.Read(buffer)

	if err != nil {
		WriteNewResponse(models.Response{
			Error: "error when reading the buffer",
		}, conn)
		return
	}

	if n > 0 {
		var request models.Request
		err = json.Unmarshal(buffer[:n], &request)
		if err != nil {
			WriteNewResponse(models.Response{
				Error: "error on request format",
			}, conn)
			return
		}

		handleRequest(request, conn)
	}
}

// CleanupSessions periodically checks for inactive sessions and reservations, and cleans them up.
// It runs every minute and checks each session and its reservations against the provided timeout.
// If a session or a reservation is inactive (i.e., its last activity time is older than the timeout),
// it is deleted from the system.
//
// Parameters:
//   - timeout: The duration after which a session or a reservation is considered inactive.
func CleanupSessions(timeout time.Duration) {
	ticker := time.NewTicker(time.Minute)
	defer ticker.Stop()

	for range ticker.C {
		for _, session := range dao.GetSessionDAO().FindAll() {
			if time.Since(session.LastTimeActive) > timeout {
				fmt.Printf("Encerrando sessão %s por inatividade\n", session.ID)
				dao.GetSessionDAO().Delete(session)
			}
		}
	}
}

// handleRequest processes incoming requests and dispatches them to the appropriate handler function.
// It reads a request from the provided net.Conn, unmarshals it into a models.Request struct, and then
// performs an action based on the request's Action field.
//
// Parameters:
//   - request: A models.Request struct containing the incoming request data.
//   - conn: A net.Conn object representing the client connection.
func handleRequest(request models.Request, conn net.Conn) {
	// switch request.Action {
	// case "login":
	// 	Login(request.Data)
	// case "get-user":
	// 	GetUserBySessionToken(request)
	// case "logout":
	// 	Logout(request)
	// case "all-routes":
	// 	AllRoutes(request.Auth, conn)
	// case "route":
	// 	Route(request.Auth, request.Data, conn)
	// case "flights":
	// 	Flights(request)
	// case "buy":
	// 	BuyTicket(request.Auth, request.Data, conn)
	// case "cancel-buy":
	// 	CancelBuy(request.Auth, request.Data, conn)
	// case "tickets":
	// 	GetTickets(request.Auth, conn)
	// }
}

// WriteNewResponse sends a response to the client over the provided net.Conn connection.
// It first checks if the response's Data field is nil. If it is, it initializes it as an empty map.
// Then, it marshals the response into JSON format. If the marshalling process encounters an error,
// it logs the error and returns without sending any response.
// After successfully marshalling the response, it writes the JSON data to the connection.
// If there is an error while writing the data, it logs the error and returns without sending any response.
//
// Parameters:
//   - response: A models.Response struct containing the response data to be sent to the client.
//   - conn: A net.Conn object representing the client connection.
func WriteNewResponse(response models.Response, conn net.Conn) {
	if response.Data == nil {
		response.Data = make(map[string]interface{})
	}

	jsonData, err := json.Marshal(response)
	if err != nil {
		fmt.Println("Error marshalling response:", err)
		return
	}
	_, err = conn.Write(jsonData)
	if err != nil {
		fmt.Println("Error writing response:", err)
	}
}

// SessionIfExists checks if a session exists for the given token.
// If a session is found, it updates the session's last activity time and returns the session along with a boolean value of true.
// If no session is found or an error occurs during the process, it returns nil and false.
//
// Parameters:
//   - token: A string representing the session token to be checked.
//
// Return:
//   - *models.Session: A pointer to the found session if it exists, or nil if no session is found or an error occurs.
//   - bool: A boolean value indicating whether a session was found (true) or not (false).
func SessionIfExists(token string) (*models.Session, bool) {
	uuid, err := uuid.Parse(token)
	if err != nil {
		return nil, false
	}
	session, err := dao.GetSessionDAO().FindById(uuid)
	if err != nil {
		return nil, false
	}
	session.LastTimeActive = time.Now()
	dao.GetSessionDAO().Update(session)
	return session, true
}
