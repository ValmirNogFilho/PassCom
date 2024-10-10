package main

import (
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"
)

func main() {
	srv := &http.Server{Addr: ":9878"}

	go func() {
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			fmt.Printf("Error: %v\n", err)
		}
	}()

	// Captura sinais de interrupção
	c := make(chan os.Signal, 1)
	signal.Notify(c, syscall.SIGINT, syscall.SIGTERM)

	// Espera pelo sinal
	<-c
	fmt.Println("Shutting down server...")

	if err := srv.Close(); err != nil {
		fmt.Printf("Error closing server: %v\n", err)
	}
}
