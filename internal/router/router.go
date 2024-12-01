package router

import (
	"github.com/gorilla/mux"

	"github.com/esponges/initial-setup/internal"
	common_validator "github.com/esponges/initial-setup/internal/common/middleware"
	"github.com/esponges/initial-setup/internal/handlers"
)

func NewRoutes(r *mux.Router) {
	// Example routes
	r.HandleFunc("/", handlers.HomeHandler).Methods("GET")
	r.HandleFunc("/health", handlers.HealthCheckHandler).Methods("GET")

	// Middleware
	headersValidator := &common_validator.CommonValidator{}
	r.Use(headersValidator.Validate)
}

func SetupRouter() internal.Application {
	r := mux.NewRouter()

	NewRoutes(r)

	return internal.Application{
		Router:        r,
		Configuration: *internal.NewConfiguration(),
	}
}
