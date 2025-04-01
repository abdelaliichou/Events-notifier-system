package collections

import (
	"encoding/json"
	"github.com/sirupsen/logrus"
	"middleware/example/internal/models"
	"middleware/example/internal/services/events"
	"net/http"
)

// SearchEventByUID
// @Tags         events
// @Summary      Search for an event by UID
// @Description  This endpoint returns an event by its unique ID
// @Param        uid query string true "Event UID"
// @Success      200 {object} models.Event
// @Failure      400 "Missing UID query parameter"
// @Failure      404 "Event not found"
// @Failure      500 "Error fetching event"
// @Router       /events [get]
func SearchEventByUID(w http.ResponseWriter, r *http.Request) {
	eventUID := r.URL.Query().Get("uid")
	if eventUID == "" {
		http.Error(w, "Missing uid query parameter", http.StatusBadRequest)
		return
	}

	// Fetch the event from the service layer
	event, err := events.GetEventByID(eventUID)
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

	// Return the event details as JSON
	w.WriteHeader(http.StatusOK)
	body, _ := json.Marshal(event)
	_, _ = w.Write(body)
}
