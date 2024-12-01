package main

import (
	"log"

	"github.com/esponges/initial-setup/internal/router"
)

func main() {
	// Initialize router
	r := router.SetupRouter()

	log.Println("run!")

	r.Run(r.Router)
}
