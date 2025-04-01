package collections

import (
	"encoding/json"
	"github.com/sirupsen/logrus"
	"middleware/example/internal/models"
	"middleware/example/internal/services/events"
	"net/http"
)

// GetEvents
// @Tags         events
// @Summary      Get all events
// @Description  This endpoint returns a list of all events
// @Success      200 {array} models.Event
// @Failure      500 "Error fetching events"
// @Router       /events [get]
func GetEvents(w http.ResponseWriter, r *http.Request) {
	// Fetch the events from the service layer
	// we don't use context because we don't need IDs
	events, err := events.GetAllEvents()
	if err != nil {
		logrus.Errorf("Error: %s", err.Error())
		customError, isCustom := err.(*models.CustomError)
		if isCustom {
			w.WriteHeader(customError.Code)
			body, _ := json.Marshal(customError)
			_, _ = w.Write(body)
		} else {
			w.WriteHeader(http.StatusInternalServerError)
		}
		return
	}

	// Always return JSON, even if no events found
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	if len(events) == 0 {
		_, _ = w.Write([]byte("[]")) // Return an empty JSON array if no events found
		return
	}

	json.NewEncoder(w).Encode(events)
}
