package internal

import (
	"github.com/esponges/initial-setup/internal/handlers"
	"github.com/go-playground/validator/v10"
)

type Configuration struct {
	Port string
	API  API
}

type API struct {
	SamplePostRequestHandler handlers.SamplePostRequestHandlerImpl
}

func NewConfiguration() *Configuration {
	// Register validators
	validate := validator.New()

	// Register handlers
	samplePostRequestHandler := handlers.NewSamplePostRequestHandler(validate)

	return &Configuration{
		Port: "8080",
		API: API{
			SamplePostRequestHandler: *samplePostRequestHandler,
		},
	}
}
