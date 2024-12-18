package server

import (
	"encoding/json"
	"fmt"
	"giro/internal/dao"
	"giro/internal/models"
	"net/http"
	"time"
)

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

	response := GetUserBySessionToken(models.Request{Auth: token})

	returnResponse(w, r, response)

}

// handleLogout handles HTTP GET requests to log out the authenticated user.
// It checks the request method to ensure it's a GET request and retrieves the user's authorization token from the request headers.
// It then constructs a Request object with the appropriate action and authorization token, and sends it to the
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
	req := models.Request{Auth: token}
	response := Logout(req)

	returnResponse(w, r, response)
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

	responseData := Login(logCred)
	returnResponse(w, r, responseData)
}

// passwordMatches checks if the provided password matches the password stored in the client's record.
//
// Parameters:
// - client: A pointer to a models.Client representing the user for whom the password needs to be checked.
// - password: A string representing the password provided by the user.
//
// Return:
//   - A boolean value indicating whether the provided password matches the password stored in the client's record.
//     Returns true if the passwords match, false otherwise.
func passwordMatches(client *models.Client, password string) bool {
	return client.Password == password
}

// login handles the login process for a user.
// It receives a data interface and a connection to the client.
// It first unmarshals the data into a LoginCredentials struct, then retrieves the client from the database using the provided username.
// If the client is not found, it sends an error response to the client and returns.
// If the client is found, it checks if the provided password matches the client's password.
// If the passwords match, it checks if the client is already logged in by searching for an active session.
// If the client is already logged in, it sends an error response to the client and returns.
// If the client is not logged in, it creates a new session for the client, stores it in the database, and sends a success response with the session token to the client.
// If the passwords do not match, it sends an error response to the client.
func Login(data interface{}) models.Response {
	var logCred models.LoginCredentials

	response := models.Response{Data: make(map[string]interface{})}

	jsonData, _ := json.Marshal(data)
	json.Unmarshal(jsonData, &logCred)

	login, err := dao.GetClientDAO().FindByUsername(logCred.Username)

	if err != nil {
		return models.Response{
			Error:  "client not found",
			Status: http.StatusUnauthorized,
		}
	}

	var session *models.Session

	if passwordMatches(login, logCred.Password) {

		if s := findUser(login); s != nil {
			return models.Response{
				Error:  "more than one user logged",
				Status: http.StatusUnauthorized,
			}

		} else {
			session = &models.Session{ClientID: login.ID, LastTimeActive: time.Now()}
			dao.GetSessionDAO().Insert(session)
		}

		token := session.ID.String()

		response.Data["token"] = token
		response.Status = http.StatusOK

	} else {
		response.Error = "invalid credentials"
		response.Status = http.StatusUnauthorized
	}
	return response
}

// findUser searches for an active session associated with a given client.
// It iterates through all sessions stored in the database and checks if the client's ID matches the session's ClientID.
// If a matching session is found, it is returned. Otherwise, nil is returned.
//
// Parameters:
// - login: A pointer to a models.Client representing the client for which the session needs to be found.
//
// Return:
//   - A pointer to a models.Session representing the active session associated with the given client.
//     If no matching session is found, nil is returned.
func findUser(login *models.Client) *models.Session {
	fmt.Println(dao.GetSessionDAO().FindAll())
	for _, s := range dao.GetSessionDAO().FindAll() {
		if s.ClientID == login.ID {
			return s
		}
	}
	return nil
}

// logout handles the logout process for a user.
// It closes the connection, prepares a response object, and checks if a session exists for the given authentication token.
// If the session is found, it deletes the session from the database and calls the removeReservations function to release reserved seats.
// Finally, it sends a success message in the response and writes it to the provided connection.
//
// Parameters:
// - auth: A string representing the authentication token.
// - conn: A net.Conn representing the connection to the client.
//
// Return:
// - None. The function writes the response directly to the connection.
func Logout(req models.Request) models.Response {
	response := models.Response{Data: make(map[string]interface{})}

	session, exists := SessionIfExists(req.Auth)

	if !exists {
		response.Error = "session not found"
		response.Status = http.StatusNotFound
		return response
	}

	dao.GetSessionDAO().Delete(session)

	response.Data["msg"] = "logout successfully made"
	response.Status = http.StatusOK
	return response
}

// getUserBySessionToken retrieves the user associated with a given session token.
// It uses the provided authentication token to find the corresponding session in the database.
// If the session is found, it retrieves the user's ID from the session and uses it to fetch the user from the database.
// If the session or user is not found, appropriate error messages are set in the response.
// The function then writes the response to the provided connection.
//
// Parameters:
// - auth: A string representing the authentication token.
// - conn: A net.Conn representing the connection to the client.
//
// Return:
// - None. The function writes the response directly to the connection.
func GetUserBySessionToken(request models.Request) models.Response {
	response := models.Response{Data: make(map[string]interface{})}

	session, exists := SessionIfExists(request.Auth)

	if !exists {
		response.Error = "session not found"
		response.Status = http.StatusNotFound
		return response
	}

	id := session.ClientID

	client, err := dao.GetClientDAO().FindById(id)

	if err != nil {
		response.Error = "client not found"
		response.Status = http.StatusNotFound
		return response
	}

	response.Data["user"] = map[string]interface{}{
		"Name":          client.Name,
		"ClientFlights": client.ClientFlights,
		"Username":      client.Username,
	}
	response.Status = http.StatusOK
	return response
}
