package main

import (
	"log"
	"os"

	"github.com/joho/godotenv"

	"github.com/esponges/initial-setup/internal/router"
)

func main() {
	// add env variables
	err := godotenv.Load()
	if err != nil {
		log.Fatal(err)
	}

	// Now you can access your environment secrets using os.Getenv()
	dbUsername := os.Getenv("DB_USERNAME")
	dbPassword := os.Getenv("DB_PASSWORD")
	apiKey := os.Getenv("API_KEY")

	log.Printf("DB_USERNAME: %s\n", dbUsername)
	log.Printf("DB_PASSWORD: %s\n", dbPassword)
	log.Printf("API_KEY: %s\n", apiKey)

	// Initialize router
	r := router.SetupRouter()

	log.Println("run!")

	r.Run(r.Router)
}
