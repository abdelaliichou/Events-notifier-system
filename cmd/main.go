package main

import (
	"encoding/json"
	"fmt"
	"github.com/gofrs/uuid"
	"github.com/sirupsen/logrus"
	"middleware/example/internal/helpers"
	"middleware/example/internal/models"
	_ "middleware/example/internal/models"
)

func main() {
	/*
		r := chi.NewRouter()

		r.Route("/collections", func(r chi.Router) {
			r.Get("/", collections.GetCollections)
			r.Route("/{id}", func(r chi.Router) {
				r.Use(collections.Ctx)
				r.Get("/", collections.GetCollection)
			})
		})

		logrus.Info("[INFO] Web server started. Now listening on *:8080")
		logrus.Fatalln(http.ListenAndServe(":8080", r))
	*/
	testingModels()
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
		All:        false,
		ResourceID: &resourceID,
	}

	// Création d'une alerte sans ressource
	alertWithoutResource := models.Alert{
		Id:    &resourceID,
		Email: "user@example.com",
		All:   true,
	}

	// Sérialisation JSON
	resourceJSON, _ := json.Marshal(resource)
	alertWithResourceJSON, _ := json.Marshal(alertWithResource)
	alertWithoutResourceJSON, _ := json.Marshal(alertWithoutResource)

	fmt.Println("Resource JSON:", string(resourceJSON))
	fmt.Println("Alert (with resource) JSON:", string(alertWithResourceJSON))
	fmt.Println("Alert (without resource) JSON:", string(alertWithoutResourceJSON))

}

func staticVariables() {
	fmt.Println("App Name:", models.AppName)
	fmt.Println("Version:", models.Version)

	fmt.Println("CERRA Mode:", models.CalendarURL("5",
		models.M1_GROUPE_3_OPTION))

	fmt.Println("CERRA Mode:", models.CalendarURL("2",
		models.M1_GROUPE_1_lANGUE,
		models.M1_GROUPE_1_lANGUE,
		models.M1_GROUPE_1_lANGUE))
}

func init() {
	db, err := helpers.OpenDB()
	if err != nil {
		logrus.Fatalf("error while opening database : %s", err.Error())
	}
	schemes := []string{
		`CREATE TABLE IF NOT EXISTS collections (
			id VARCHAR(255) PRIMARY KEY NOT NULL UNIQUE,
			content VARCHAR(255) NOT NULL
		);`,
	}
	for _, scheme := range schemes {
		if _, err := db.Exec(scheme); err != nil {
			logrus.Fatalln("Could not generate table ! Error was : " + err.Error())
		}
	}
	helpers.CloseDB(db)
}
