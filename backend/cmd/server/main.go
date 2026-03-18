package main

import (
	"log"
	"net/http"

	"github.com/druabello/airport-employee-management/internal/handlers"
)

func main() {
	log.Println("Setting up handler")

	mux := http.NewServeMux()
	mux.HandleFunc("GET /helloworld", handlers.HandleHelloWorld)

	log.Println("Server listening on port 8080")
	log.Fatalln(http.ListenAndServe(":8080", mux))
}
