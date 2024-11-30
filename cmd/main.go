package main

import (
	"log"

	// "net/http"

	"github.com/esponges/initial-setup/internal/router"
)

func main() {
	// Initialize router
	r := router.SetupRouter()

	log.Println("run!")

	// methodsOk := handlers.AllowedMethods([]string{"GET"})
	// server := &http.Server{
	// 	Addr:    ":" + r.Configuration.Port,
	// 	Handler: handlers.CORS(methodsOk)(r.Router),
	// }
	r.Run(r.Router)

	// Start server
	log.Println("Starting server on", r.Configuration.Port)
	// log.Fatal(server.ListenAndServe())
}
