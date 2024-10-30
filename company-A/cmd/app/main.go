package main

import (
	"fmt"
	"net"
	"os"
	"time"
	"vendepass/internal/server"
)

const (
	port      = ":8888"
	timeLimit = 30 * time.Minute
)

// main function is the entry point of the application.
// It sets up the server, handles incoming connections, and manages flight reservations.
func main() {
	listener, err := net.Listen("tcp", port)
	if err != nil {
		fmt.Println("Listener não criado", err)
		os.Exit(1)
	}
	defer listener.Close()

	fmt.Println("servidor ouvindo na porta :8888")

	go server.CleanupSessions(timeLimit)

	for {
		conn, err := listener.Accept()
		if err != nil {
			fmt.Println("erro ao aceitar conexão", err)
			continue
		}

		go server.HandleConn(conn)
	}

}
