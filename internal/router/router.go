package router

import (
	"github.com/esponges/initial-setup/internal/handlers"
	"github.com/gorilla/mux"
)

func SetupRouter() *mux.Router {
	r := mux.NewRouter()

	// Example routes
	r.HandleFunc("/", handlers.HomeHandler).Methods("GET")
	r.HandleFunc("/health", handlers.HealthCheckHandler).Methods("GET")

	return r
}
