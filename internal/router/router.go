package router

import (
	configuration "github.com/esponges/initial-setup/internal"
	"github.com/esponges/initial-setup/internal/handlers"
	"github.com/gorilla/mux"
)

type Application struct {
	Router        *mux.Router
	Configuration configuration.Configuration
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
		Router:        r,
		Configuration: *configuration.NewConfiguration(),
	}
}
