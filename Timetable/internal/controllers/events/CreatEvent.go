package collections

import (
	"encoding/json"
	"github.com/sirupsen/logrus"
	"middleware/example/internal/models"
	"middleware/example/internal/services/events"
	"net/http"
)

// CreatEvent handles creating a new event
func CreatEvent(w http.ResponseWriter, r *http.Request) {

	var newEvent models.Event
	if err := json.NewDecoder(r.Body).Decode(&newEvent); err != nil {
		logrus.Errorf("Failed to decode request body: %s", err.Error())
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	event, err := events.CreatEvent(newEvent)

	if err != nil {
		logrus.Errorf("Error creating event: %s", err.Error())
		http.Error(w, "Error creating event", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(event)
}
