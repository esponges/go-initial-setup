package router

import (
	"log"
	"net/http"
	"time"

	gorillaHandlers "github.com/gorilla/handlers"
	"github.com/gorilla/mux"

	configuration "github.com/esponges/initial-setup/internal"
	"github.com/esponges/initial-setup/internal/handlers"
)

type Application struct {
	Router        *mux.Router
	Configuration configuration.Configuration
}

func (app *Application) Run(router *mux.Router) {
	methodsOk := gorillaHandlers.AllowedMethods([]string{"GET"})

	log.Println(app.Router)
	log.Println(router)
	log.Println("starting server")
	server := &http.Server{
		Addr:    ":" + app.Configuration.Port,
		Handler: gorillaHandlers.CORS(methodsOk)(app.Router),
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
