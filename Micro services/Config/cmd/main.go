package main

import (
	"encoding/json"
	"fmt"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/gofrs/uuid"
	"github.com/sirupsen/logrus"
	alertsHandler "middleware/example/internal/controllers/alerts"
	resourcesHandler "middleware/example/internal/controllers/resources"
	"middleware/example/internal/helpers"
	"middleware/example/internal/models"
	"net/http"
)

func main() {

	routes()

}

func routes() {

	r := chi.NewRouter()

	// Middleware for logging and recovery
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)

	// ---------------------- ALERT ROUTES ----------------------
	r.Route("/alerts", func(r chi.Router) {
		r.Post("/", alertsHandler.CreateAlert) // Create an alert
		r.Get("/", alertsHandler.GetAlerts)    // Get all alerts

		r.Route("/{id}", func(r chi.Router) {
			r.Use(alertsHandler.CtxAlert)            // Middleware to extract alert ID
			r.Get("/", alertsHandler.GetAlert)       // Get alert by ID
			r.Put("/", alertsHandler.UpdateAlert)    // Update alert by ID
			r.Delete("/", alertsHandler.DeleteAlert) // Delete alert by ID
		})
	})

	// ---------------------- RESOURCE ROUTES ----------------------
	r.Route("/resources", func(r chi.Router) {
		r.Post("/", resourcesHandler.CreateResource) // Create a resource
		r.Get("/", resourcesHandler.GetAllResources) // Get all resources

		r.Route("/{id}", func(r chi.Router) {
			r.Use(resourcesHandler.CtxResource)            // Middleware to extract resource ID
			r.Get("/", resourcesHandler.GetResource)       // Get resource by ID
			r.Put("/", resourcesHandler.UpdateResource)    // Update resource by ID
			r.Delete("/", resourcesHandler.DeleteResource) // Delete resource by ID
		})
	})

	// Start the server
	logrus.Info("[INFO] Web server started. Now listening on *:8080")
	logrus.Fatalln(http.ListenAndServe(":8080", r))

}

func testingModels() {

	// Création d'un UUID pour la ressource
	resourceID, _ := uuid.NewV4()

	// Création d'une ressource
	resource := models.Resource{
		Id:    &resourceID,
		UcaID: 12345,
		Name:  "Salle Informatique",
	}

	// Création d'une alerte associée à la ressource
	alertWithResource := models.Alert{
		Id:         &resourceID,
		Email:      "user@example.com",
		IsAll:      false,
		ResourceID: &resourceID,
	}

	// Création d'une alerte sans ressource
	alertWithoutResource := models.Alert{
		Id:         &resourceID,
		Email:      "user@example.com",
		IsAll:      true,
		ResourceID: nil,
	}

	// Sérialisation JSON
	resourceJSON, _ := json.Marshal(resource)
	alertWithResourceJSON, _ := json.Marshal(alertWithResource)
	alertWithoutResourceJSON, _ := json.Marshal(alertWithoutResource)

	fmt.Println("Resource JSON:", string(resourceJSON))
	fmt.Println("Alert (with resource) JSON:", string(alertWithResourceJSON))
	fmt.Println("Alert (without resource) JSON:", string(alertWithoutResourceJSON))

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

	// Define table schemas for Resource and Alert
	schemes := []string{
		models.RESOURCES_TABLE,
		models.ALERTS_TABLE,
	}

	// Execute each table creation query
	for _, scheme := range schemes {
		if _, err := db.Exec(scheme); err != nil {
			logrus.Fatalln("Could not create table! Error: " + err.Error())
		}
	}

	helpers.CloseDB(db)

}
