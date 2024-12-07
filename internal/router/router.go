package router

import (
	"github.com/gorilla/mux"

	"github.com/esponges/initial-setup/internal"
	common_validator "github.com/esponges/initial-setup/internal/common/middleware"
	"github.com/esponges/initial-setup/internal/handlers"
)

func NewRoutes(r *mux.Router, config *internal.Configuration) {
	// Example routes
	r.HandleFunc("/", handlers.HomeHandler).Methods("GET")
	r.HandleFunc("/health", handlers.HealthCheckHandler).Methods("GET")
	r.HandleFunc("/sample_post_request", config.API.SamplePostRequestHandler.SamplePostRequestHandler).Methods("POST")
	r.HandleFunc("/create_singers", config.API.CreateSingersHandler.CreateSingersHandler).Methods("POST")

	// Middleware
	headersValidator := &common_validator.CommonValidator{}
	r.Use(headersValidator.Validate)
}

func SetupRouter() internal.Application {
	r := mux.NewRouter()

	config := internal.NewConfiguration()

	NewRoutes(r, config)

	return internal.Application{
		Router:        r,
		Configuration: *config,
	}
}
