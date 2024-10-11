package main

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"os/signal"
	"syscall"
)

func main() {
	srv := &http.Server{Addr: ":9877"}

	go func() {
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			fmt.Printf("Error: %v\n", err)
		}
	}()

	http.HandleFunc("/A", handleA)
	http.HandleFunc("/C", handleC)

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

func allowCrossOrigin(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")

	if r.Method == http.MethodOptions {
		w.WriteHeader(http.StatusNoContent)
		return
	}
}

func handleA(w http.ResponseWriter, r *http.Request) {
	allowCrossOrigin(w, r)
	resp, err := http.Get("http://localhost:9876/callA")
	if err != nil {
		fmt.Println("problem", err)
		return
	}
	fmt.Println(resp)
	b, _ := io.ReadAll(resp.Body)

	w.Write(b)
}

func handleC(w http.ResponseWriter, r *http.Request) {

}
