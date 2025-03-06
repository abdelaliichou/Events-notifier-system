package collections

import (
	"encoding/json"
	"github.com/sirupsen/logrus"
	"middleware/example/internal/models"
	"middleware/example/internal/services/events"
	"net/http"
)

// GetEvent retrieves an event from the database using the Event ID from the context
func GetEvent(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	//eventUID, _ := ctx.Value("eventID").(uuid.UUID)
	eventUID := ctx.Value("eventID")

	// Fetch the event from the service layer
	event, err := events.GetEventByID(eventUID.(string))
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
