package main

import (
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/sirupsen/logrus"
	eventsHandler "middleware/example/internal/controllers/events"
	"middleware/example/internal/helpers"
	"middleware/example/internal/models"
	"middleware/example/internal/mq"
	"net/http"
)

func main() {

	// Starting consumer
	go mq.StartStreamConsumer()

	// Starting producer
	go mq.InitStream()

	routes()

}

func routes() {

	r := chi.NewRouter()

	// Middleware for logging and recovery
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)

	// ---------------------- EVENTS ROUTES ----------------------
	r.Route("/events", func(r chi.Router) {
		//r.Post("/", eventsHandler.CreatEvent) // Create an event
		r.Get("/", eventsHandler.GetEvents)              // Get all event
		r.Get("/search", eventsHandler.SearchEventByUID) // Search event by UID

		r.Route("/{id}", func(r chi.Router) {
			//r.Use(eventsHandler.CtxAlert)      // Middleware to extract event ID
			//r.Get("/", eventsHandler.GetEvent) // Get event by UID
			//r.Put("/", eventsHandler.UpdateEvent) // Update event by ID
			//r.Delete("/", eventsHandler.DeleteEvent) // Delete event by ID
		})
	})

	// readme
	// events?uid=		-- DONE
	// swagger
	// docker
	// function to send only what have changed on the event -- DONE
	// handle how to send the changes only in the mq from the consumer in the Timetable --DONE
	// handle how to receive the changes from the consumer in config Timetable -- not yet ( done in chatgpt concerns consumer of config )
	// handle how to send email structure with the changes
	// send multiple emails to same person in case there where multiple changed events
	// alerter dans un autre code
	// separer le code MVC

	// Start the server
	logrus.Info("[INFO] Web server started. Now listening on *:8090")
	logrus.Fatalln(http.ListenAndServe(":8090", r))

}

func init() {
	db, err := helpers.OpenDB()
	if err != nil {
		logrus.Fatalf("Error while opening database: %s", err.Error())
	}

	// Activer explicitement les FOREIGN KEYS ( Ne faites pas trop confiance aux constraints avec SQLite )
	if _, err := db.Exec("PRAGMA foreign_keys = ON;"); err != nil {
		logrus.Fatalln("Could not enable foreign key support: " + err.Error())
	}

	// Define table schemas for Events and Resource ( Resource is to represent resourceIds[] of an event )
	schemes := []string{
		models.CREAT_EVENT,
		models.CREAT_RESOURCE,
	}

	// Execute each table creation query
	for _, scheme := range schemes {
		if _, err := db.Exec(scheme); err != nil {
			logrus.Fatalln("Could not create table! Error: " + err.Error())
		}
	}

	helpers.CloseDB(db)

}
