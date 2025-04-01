package main

import (
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/sirupsen/logrus"
	httpSwagger "github.com/swaggo/http-swagger"
	_ "middleware/example/docs"
	alertsHandler "middleware/example/internal/controllers/alerts"
	resourcesHandler "middleware/example/internal/controllers/resources"
	"middleware/example/internal/helpers"
	"middleware/example/internal/models"
	"net/http"
)

// @title          Alerts & Resources API
// @version        1.0
// @description    API for managing alerts and resources
// @host           localhost:8080
// @BasePath       /alerts & /resources

func main() {

	routes()
	// swagger documentation : http://localhost:8080/swagger/index.html

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

	// ---------------------- ALERT ROUTES ----------------------
	r.Route("/alerts", func(r chi.Router) {
		r.Post("/", alertsHandler.CreateAlert) // Create an alert
		r.Get("/", alertsHandler.GetAlerts)    // Get all alerts

		r.Route("/{id}", func(r chi.Router) {
			r.Use(alertsHandler.CtxAlert)      // Middleware to extract alert ID
			r.Get("/", alertsHandler.GetAlert) // Get alert by ID
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
