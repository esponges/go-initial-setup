package main

import (
	"log"
	"net/http"

	"github.com/esponges/initial-setup/internal/router"
)

func main() {
	// Initialize router
	r := router.SetupRouter()

	// Start server
	log.Println("Starting server on :8080")
	log.Fatal(http.ListenAndServe(":8080", r))
}
