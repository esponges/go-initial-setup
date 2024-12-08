package main

import (
	"log"

	"github.com/joho/godotenv"

	"github.com/esponges/initial-setup/internal/router"
)

func main() {
	// add env variables
	err := godotenv.Load()
	if err != nil {
		log.Fatal(err)
	}

	// Initialize router
	r := router.SetupRouter()

	log.Println("run!")

	r.Run(r.Router)
}
