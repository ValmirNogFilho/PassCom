package main

import (
	"log"
	"passcom/boreal/src/server"
)

func main() {
	var System = server.GetInstance()
	err := System.StartServer()
	if err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}
