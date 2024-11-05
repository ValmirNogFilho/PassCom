package server

import (
	"encoding/json"
	"giro/internal/dao"
	"giro/internal/models"
	"log"
	"net/http"
	"strconv"
)

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

func handleWishlist(w http.ResponseWriter, r *http.Request) {
	allowCrossOrigin(w, r)
	switch r.Method {
	case http.MethodGet:
		handleGetWishlist(w, r)
	case http.MethodPost:
		handleAddToWishlist(w, r)
	case http.MethodDelete:
		handleRemoveFromWishlist(w, r)
	default:
		http.Error(w, "only GET, POST, DELETE allowed", http.StatusMethodNotAllowed)
		return
	}
}

func handleRemoveFromWishlist(w http.ResponseWriter, r *http.Request) {
	allowCrossOrigin(w, r)
	token := r.Header.Get("Authorization")

	queryParams := r.URL.Query()

	id := queryParams.Get("id")

	idUint, err := strconv.ParseUint(id, 10, 32)
	if err != nil {
		returnResponse(w, r, models.Response{
			Status: http.StatusBadRequest,
		})
		return
	}

	response := DeleteFromWishlist(uint(idUint),
		models.Request{
			Auth: token,
		},
	)

	returnResponse(w, r, response)
}

func handleAddToWishlist(w http.ResponseWriter, r *http.Request) {
	allowCrossOrigin(w, r)
	token := r.Header.Get("Authorization")

	var addWish models.WishlistOperation

	err := json.NewDecoder(r.Body).Decode(&addWish)

	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	response := AddToWishlist(
		models.Request{
			Auth: token,
			Data: addWish,
		},
	)
	returnResponse(w, r, response)
}

func handleGetWishlist(w http.ResponseWriter, r *http.Request) {
	token := r.Header.Get("Authorization")
	response := GetWishlist(
		models.Request{
			Auth: token,
		},
	)
	returnResponse(w, r, response)
}

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

func returnResponse(w http.ResponseWriter, r *http.Request, responseData models.Response) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(responseData.Status)
	json.NewEncoder(w).Encode(responseData)
}

func (s *System) handleConnect(w http.ResponseWriter, r *http.Request) {
	allowCrossOrigin(w, r)

	// Verifica o método da solicitação
	switch r.Method {
	case http.MethodPost:
		// Processa a solicitação para adicionar uma nova conexão
		var message models.Message
		err := json.NewDecoder(r.Body).Decode(&message)
		if err != nil {
			http.Error(w, "Invalid request payload", http.StatusBadRequest)
			return
		}

		// Decodifica o Body para map[string]interface{} para acessar Address e Port
		body, ok := message.Body.(map[string]interface{})
		if !ok {
			http.Error(w, "Invalid body format", http.StatusBadRequest)
			return
		}

		name, ok := body["Name"].(string)
		if !ok {
			http.Error(w, "Invalid name format", http.StatusBadRequest)
			return
		}

		address, ok := body["Address"].(string)
		if !ok {
			http.Error(w, "Invalid Address format", http.StatusBadRequest)
			return
		}

		port, ok := body["Port"].(string)
		if !ok {
			http.Error(w, "Invalid Port format", http.StatusBadRequest)
			return
		}

		// Cria uma nova conexão com os dados extraídos
		newConnection := models.Connection{
			Name:     name,
			Address:  address,
			Port:     port,
			IsOnline: true,
		}

		// Adiciona a nova conexão ao sistema
		err = s.AddConnection(message.From, newConnection)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		log.Printf("New connection from %s: %s:%s", message.From, address, port)

		// Monta a resposta como models.Message contendo o novo models.Connection
		responseMessage, err := models.CreateMessage(s.ServerId.String(), message.From, s.VectorClock, map[string]interface{}{"Name": s.ServerName})

		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		// Define o cabeçalho Content-Type e envia o JSON da resposta
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusCreated)
		if err := json.NewEncoder(w).Encode(responseMessage); err != nil {
			http.Error(w, "Failed to encode response message", http.StatusInternalServerError)
		}

	case http.MethodDelete:
		// Processa a solicitação para remover uma conexão existente
		var message models.Message
		err := json.NewDecoder(r.Body).Decode(&message)
		if err != nil {
			http.Error(w, "Invalid request payload", http.StatusBadRequest)
			return
		}

		s.RemoveConnection(message.From)
		log.Printf("Connection removed for server %s", message.From)

		w.WriteHeader(http.StatusOK)
		w.Write([]byte("Connection removed successfully"))

	default:
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
	}
}

func (s *System) handleHeartbeat(w http.ResponseWriter, r *http.Request) {
	// Define o tipo de resposta e status OK
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	// Decodifica o *heartbeat* recebido
	var receivedMessage models.Message
	if err := json.NewDecoder(r.Body).Decode(&receivedMessage); err != nil {
		http.Error(w, "Invalid request payload", http.StatusBadRequest)
		return
	}

	// Bloqueia o mutex para manipular o VectorClock de maneira segura
	s.Lock.Lock()
	defer s.Lock.Unlock()

	// Incrementa o relógio do sistema para indicar que o heartbeat foi recebido
	s.IncrementClock()

	log.Print("Received heartbeat from ", s.Connections[receivedMessage.From].Name)

	// Atualiza o VectorClock com base no *heartbeat* recebido
	s.UpdateClock(receivedMessage.VectorClock)

	// Cria uma nova mensagem de resposta com o VectorClock atualizado
	responseMessage, err := models.CreateMessage(s.ServerId.String(), receivedMessage.From, s.VectorClock, "Healthy")

	if err != nil {
		log.Printf("Error creating heartbeat response message: %v", err)
		return
	}

	// Codifica a mensagem de resposta como JSON
	if err := json.NewEncoder(w).Encode(responseMessage); err != nil {
		log.Printf("Error encoding heartbeat response: %v", err)
	}

	log.Print("Sent heartbeat response to ", s.Connections[receivedMessage.From].Name)
}

func (s *System) handleDatabase(w http.ResponseWriter, r *http.Request) {
	// Define o tipo de resposta como JSON e status OK
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	// Usa o DAO para buscar todos os voos cujo campo Company é igual a s.ServerName
	flights, err := dao.GetFlightDAO().FindByCompany(s.ServerName)
	if err != nil {
		http.Error(w, "Failed to retrieve flights", http.StatusInternalServerError)
		log.Printf("Error retrieving flights from database: %v", err)
		return
	}

	// Codifica a lista de voos como JSON e envia na resposta
	if err := json.NewEncoder(w).Encode(flights); err != nil {
		http.Error(w, "Failed to encode response", http.StatusInternalServerError)
		log.Printf("Error encoding flights to JSON: %v", err)
	}
}
