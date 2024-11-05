package server

import (
	"bufio"
	"encoding/json"
	"log"
	"net"
	"net/http"
	"rumos/internal/models"
	"strings"
	"syscall"
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

func returnResponse(w http.ResponseWriter, r *http.Request, responseData models.Response) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(responseData.Status)
	json.NewEncoder(w).Encode(responseData)
}

// HandleCLIServer starts a TCP server listening on the specified CLIPORT.
// This server accepts incoming connections and handles them using the handleCLIConnection function.
// The server logs any errors during initialization, listening, or accepting connections.
//
// The function performs the following steps:
// 1. Listens for incoming TCP connections on the specified CLIPORT.
// 2. If an error occurs during listening, logs the error and terminates the program.
// 3. Accepts incoming connections and starts a new goroutine to handle each connection using the handleCLIConnection function.
// 4. Logs any errors that occur during connection acceptance.
func (s *System) HandleCLIServer() {
	listener, err := net.Listen("tcp", CLIPORT)
	if err != nil {
		log.Fatal("Error initiating server:", err)
	}
	defer listener.Close()

	log.Println("See your server working on the TCP CLI server on port", CLIPORT)

	for {
		conn, err := listener.Accept()
		if err != nil {
			log.Println("Error accepting connection:", err)
			continue
		}

		go s.handleCLIConnection(conn)
	}
}

// handleCLIConnection handles incoming connections from the CLI server.
// It reads commands from the connection, processes them, and responds accordingly.
//
// Parameters:
// - conn: The net.Conn object representing the connection from the client.
//
// Return:
// This function does not return any value.
func (s *System) handleCLIConnection(conn net.Conn) {
	defer conn.Close()
	conn.Write([]byte("Welcome to the CLI server!\nType 'help' to see the commands.\n"))

	scanner := bufio.NewScanner(conn)
	for scanner.Scan() {
		input := strings.TrimSpace(scanner.Text())
		parts := strings.Fields(input)

		if len(parts) == 0 {
			conn.Write([]byte("Empty command.\n"))
			continue
		}

		command := parts[0]
		args := parts[1:]

		switch command {
		case "help":
			conn.Write(
				[]byte("Available commands:" +
					"\n- help: to see commands" +
					"\n- info: to see server informations" +
					"\n- addconn <address> <port>: to add a new connection" +
					"\n- quit: to close the connection" +
					"\n- shutdown: to shut down the server\n"))

		case "info":
			conn.Write([]byte(s.getServerInfo()))

		case "addconn":
			if len(args) < 2 {
				conn.Write([]byte("Error: 'addconn' requires two arguments (address, port).\n"))
			} else {
				address := args[0]
				connPort := args[1]
				conn.Write([]byte("Requesting connection to " + address + ":" + connPort + "...\n"))
				s.RequestConnection(address, connPort)
				id, serverConn := s.FindConnectionByName(address)
				if serverConn == nil {
					conn.Write([]byte("Connection not found.\n"))
				} else {
					conn.Write([]byte("Requesting database from " + address + ":" + connPort + "...\n"))
					s.RequestDatabase(id, address, connPort)
					conn.Write([]byte("Sending database to " + address + ":" + connPort + "...\n"))
					s.SendDatabase(id, address, connPort)
				}
			}

		case "rmconn":
			if len(args) < 1 {
				conn.Write([]byte("Error: 'rmconn' requires one argument (connection ID).\n"))
			} else {
				connId := args[0]
				id, serverConn := s.FindConnectionByName(connId)
				if serverConn == nil {
					conn.Write([]byte("Connection not found.\n"))
				} else {
					name := serverConn.Name
					conn.Write([]byte("Requesting database removal from " + name + "...\n"))
					s.RequestDatabaseRemoval(id, serverConn.Address, serverConn.Port)
					conn.Write([]byte("Requested disconnection from " + serverConn.Address + ":" + serverConn.Port + "...\n"))
					s.RequestDisconnection(serverConn.Address, serverConn.Port)
					conn.Write([]byte("Removing connection from " + name + "...\n"))
					s.RemoveConnection(id)
					conn.Write([]byte("Removing database from " + name + "...\n"))
					s.RemoveDatabase(connId)
				}
			}

		case "quit":
			conn.Write([]byte("Closing CLI...\n"))
			return

		case "shutdown":
			conn.Write([]byte("Shutting down the server...\n"))
			go func() {
				s.shutdown <- syscall.SIGTERM
			}()
			return

		default:
			conn.Write([]byte("Command not found.\n"))
		}
	}
}
