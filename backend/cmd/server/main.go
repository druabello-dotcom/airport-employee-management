package main

import (
	"log"
	"net/http"

	"github.com/druabello-dotcom/airport-employee-management/internal/handlers"
)

func main() {
	log.Println("Setting up handler")

	mux := http.NewServeMux()
	mux.HandleFunc("GET /helloworld", handlers.HandleHelloWorld)
	mux.HandleFunc("POST /checkpoints", handlers.HandleCheckpoints)

	handler := handlers.WithCORS(mux)

	log.Println("Server listening on port 8080")
	log.Fatalln(http.ListenAndServe(":8080", handler))
}
