package handlers

import (
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/esponges/initial-setup/internal/common"
	"github.com/esponges/initial-setup/internal/handlers/create_singers_handler"
	"github.com/esponges/initial-setup/internal/handlers/sample_post_request"
	"github.com/go-playground/validator/v10"
)

type Handler struct {
	HomeHandler        func(w http.ResponseWriter, r *http.Request)
	HealthCheckHandler func(w http.ResponseWriter, r *http.Request)
}

// TODO: handlers should be moved to their own packages
func HomeHandler(w http.ResponseWriter, r *http.Request) {
	log.Println(time.Now().Date())
	fmt.Fprintf(w, "Welcome to My Project!")
}

func HealthCheckHandler(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("OK"))
}

type SamplePostRequestHandlerImpl struct {
	validator *validator.Validate
}

func NewSamplePostRequestHandler(validator *validator.Validate) *SamplePostRequestHandlerImpl {
	return &SamplePostRequestHandlerImpl{
		validator: validator,
	}
}

func (s *SamplePostRequestHandlerImpl) SamplePostRequestHandler(w http.ResponseWriter, r *http.Request) {
	// todo: impl validation
	var req sample_post_request.SamplePostRequest
	body, err := common.UnmarshalAndValidateRequest(r, &req, s.validator)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("Bad Request"))
	} else {
		log.Println("Correct Request")
		w.WriteHeader(http.StatusOK)
		w.Write(body)
	}
}

type CreateSingersHandlerImpl struct {
	validator *validator.Validate
}

func NewCreateSingersHandler(validator *validator.Validate) *CreateSingersHandlerImpl {
	return &CreateSingersHandlerImpl{
		validator: validator,
	}
}

func (c *CreateSingersHandlerImpl) CreateSingersHandler(w http.ResponseWriter, r *http.Request) {
	var req create_singers_handler.CreateSingersRequest
	body, err := common.UnmarshalAndValidateRequest(r, &req, c.validator)
	log.Println("body: " + string(body))
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("Bad Request"))
	} else {
		log.Println("Correct Request")
		w.WriteHeader(http.StatusOK)
		w.Write(body)
	}
}
