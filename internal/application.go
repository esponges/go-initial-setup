package internal

import (
	"log"
	"net/http"
	"time"

	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"
)

type Application struct {
	Router        *mux.Router
	Configuration Configuration
}

func (app *Application) Run(router *mux.Router) {
	methodsOk := handlers.AllowedMethods([]string{"GET"})

	log.Println("starting server")
	server := &http.Server{
		Addr:    ":" + app.Configuration.Port,
		Handler: handlers.CORS(methodsOk)(app.Router),
	}

	log.Println("listen and serve!")
	errChan := make(chan error)
	go func() {
		errChan <- server.ListenAndServe()
	}()

	// The main go routine has exited by now
	// we need to start a new go routine to add logs after starting the server
	go func() {
		for {
			log.Println("Server is running...")
			time.Sleep(10 * time.Second)
		}
	}()

	err := <-errChan
	log.Println(err.Error())
}
