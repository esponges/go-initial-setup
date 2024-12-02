package handlers

import (
	"fmt"
	"log"
	"net/http"
	"time"

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

type SamplePostRequestHandler struct {
	validator *validator.Validate
}

func NewSamplePostRequestHandler(validator *validator.Validate) *SamplePostRequestHandler {
	return &SamplePostRequestHandler{
		validator: validator,
	}
}

func (s *SamplePostRequestHandler) samplePostRequestHandler(w http.ResponseWriter, r *http.Request) {
	// todo: impl validation

	w.WriteHeader(http.StatusOK)
	w.Write([]byte("OK"))
}
