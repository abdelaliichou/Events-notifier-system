package collections

import (
	"encoding/json"
	"github.com/gofrs/uuid"
	"github.com/sirupsen/logrus"
	"middleware/example/internal/models"
	"middleware/example/internal/services/events"
	"net/http"
)

// UpdateEvent updates an existing event by its ID
func UpdateEvent(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	eventID, _ := ctx.Value("eventID").(uuid.UUID)

	var updatedEvent models.Event
	if err := json.NewDecoder(r.Body).Decode(&updatedEvent); err != nil {
		logrus.Errorf("Failed to decode request body: %s", err.Error())
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	updatedEvent.Id = &eventID // Ensure the ID matches the request param

	err := events.UpdateEvent(updatedEvent)
	if err != nil {
		logrus.Errorf("Error updating event: %s", err.Error())
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

	w.WriteHeader(http.StatusOK)
	body, _ := json.Marshal(map[string]string{"message": "Event updated successfully"})
	_, _ = w.Write(body)
}
