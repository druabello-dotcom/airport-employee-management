package handlers

import (
	"log"
	"net/http"
)

func HandleHelloWorld(w http.ResponseWriter, r *http.Request) {
	if _, err := w.Write([]byte("Hello World!\n")); err != nil {
		log.Printf("Error writing hello world resp: %v", err)
	}
}
