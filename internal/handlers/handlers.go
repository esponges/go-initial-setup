package handlers

import (
	"fmt"
	"log"
	"net/http"
	"time"
)

type Handler struct {
	HealthCheckHandler func(w http.ResponseWriter, r *http.Request)
}

func HomeHandler(w http.ResponseWriter, r *http.Request) {
	log.Println(time.Now().Date())
	fmt.Fprintf(w, "Welcome to My Project!")
}

func HealthCheckHandler(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("OK"))
}
