package main

import (
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/sirupsen/logrus"
	httpSwagger "github.com/swaggo/http-swagger"
	_ "middleware/example/docs"
	eventsHandler "middleware/example/internal/controllers/events"
	"middleware/example/internal/helpers"
	"middleware/example/internal/models"
	"middleware/example/internal/mq"
	"net/http"
)

// @title          Events Timetable API
// @version        1.0
// @description    API for managing events in the timetable
// @host           localhost:8090
// @BasePath       /events

func main() {

	// Starting consumer
	go mq.StartStreamConsumer()

	// Starting producer
	go mq.InitStream()

	routes()

	// swagger documentation : http://localhost:8090/swagger/index.html/index.html
 	// github repo : https://github.com/abdelaliichou/Events-notifier-system

}

func routes() {

	r := chi.NewRouter()

	// Middleware for logging and recovery
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)

	// ---------------------- SWAGGER ROUTES ----------------------
	r.Route("/swagger", func(r chi.Router) {
		r.Get("/*", httpSwagger.WrapHandler)
	})

	// ---------------------- EVENTS ROUTES ----------------------
	r.Route("/events", func(r chi.Router) {
		r.Get("/", eventsHandler.GetEvents)              // Get all event
		r.Get("/search", eventsHandler.SearchEventByUID) // Search event by UID

		r.Route("/{id}", func(r chi.Router) {
			//r.Use(eventsHandler.CtxAlert)      // Middleware to extract event ID
			//r.Get("/", eventsHandler.GetEvent) // Get event by UID
		})
	})

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
