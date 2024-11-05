package server

import (
	"log"
	"net"
	"os"
)

func getLocalIP() string {
	addrs, err := net.InterfaceAddrs()
	if err != nil {
		log.Fatalf("Error getting local IP: %v", err)
	}

	for _, addr := range addrs {
		if ipNet, ok := addr.(*net.IPNet); ok && !ipNet.IP.IsLoopback() {
			if ipNet.IP.To4() != nil {
				return ipNet.IP.String()
			}
		}
	}
	return "127.0.0.1"
}

func getPort() string {
	// Tenta obter a porta da variável de ambiente PORT
	port := os.Getenv("PORT")
	if port == "" {
		port = PORT // Porta padrão se a variável não estiver definida
		log.Printf("PORT environment variable not set. Using default port: %s", port)
	} else {
		log.Printf("Using port from environment variable PORT: %s", port)
	}
	return port
}
