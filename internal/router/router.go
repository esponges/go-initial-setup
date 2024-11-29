package router

import (
	"github.com/esponges/initial-setup/internal/handlers"
	"github.com/gorilla/mux"
)

type Application struct {
	Router *mux.Router
}

func NewRoutes(r *mux.Router) *mux.Router {
	// Example routes
	r.HandleFunc("/", handlers.HomeHandler).Methods("GET")
	r.HandleFunc("/health", handlers.HealthCheckHandler).Methods("GET")

	return r
}

func SetupRouter() Application {
	r := mux.NewRouter()

	NewRoutes(r)

	return Application{
		Router: r,
	}
}
