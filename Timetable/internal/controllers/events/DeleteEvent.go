package collections

import (
	"encoding/json"
	"github.com/gofrs/uuid"
	"github.com/sirupsen/logrus"
	"middleware/example/internal/models"
	"middleware/example/internal/services/events"
	"net/http"
)

// DeleteEvent handles the deletion of an event
func DeleteEvent(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	eventID, _ := ctx.Value("eventID").(uuid.UUID)

	err := events.DeleteEvent(eventID)
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

	// Return success response (200 OK)
	w.WriteHeader(http.StatusOK)
	body, _ := json.Marshal(map[string]string{"message": "Event deleted successfully"})
	_, _ = w.Write(body)
}
