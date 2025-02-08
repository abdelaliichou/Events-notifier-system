package collections

import (
	"encoding/json"
	"github.com/sirupsen/logrus"
	"middleware/example/internal/models"
	services "middleware/example/internal/services/alerts"
	"net/http"
)

// CreateAlert handles creating a new alert
func CreateAlert(w http.ResponseWriter, r *http.Request) {

	var newAlert models.Alert
	if err := json.NewDecoder(r.Body).Decode(&newAlert); err != nil {
		logrus.Errorf("Failed to decode request body: %s", err.Error())
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	alert, err := services.CreateAlert(newAlert)

	if err != nil {
		logrus.Errorf("Error creating alert: %s", err.Error())
		http.Error(w, "Error creating alert", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(alert)
}
