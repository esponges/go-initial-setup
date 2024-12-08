package internal

import (
	"github.com/esponges/initial-setup/internal/handlers"
	"github.com/esponges/initial-setup/internal/handlers/create_singer_handler"
	"github.com/go-playground/validator/v10"
)

type Configuration struct {
	Port string
	API  API
}

type API struct {
	SamplePostRequestHandler handlers.SamplePostRequestHandlerImpl
	CreateSingersHandler     create_singer_handler.CreateSingersHandlerImpl
}

func NewConfiguration() *Configuration {
	// Register validators
	validate := validator.New()

	// Register handlers
	samplePostRequestHandler := handlers.NewSamplePostRequestHandler(validate)
	createSingersRequestHandler := create_singer_handler.NewCreateSingersHandler(validate)

	return &Configuration{
		Port: "8080",
		API: API{
			SamplePostRequestHandler: *samplePostRequestHandler,
			CreateSingersHandler:     *createSingersRequestHandler,
		},
	}
}
