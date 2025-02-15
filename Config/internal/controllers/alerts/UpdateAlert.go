package collections

import (
	"encoding/json"
	"github.com/gofrs/uuid"
	"github.com/sirupsen/logrus"
	"middleware/example/internal/models"
	services "middleware/example/internal/services/alerts"
	"net/http"
)

// UpdateAlert updates an existing alert by its ID
func UpdateAlert(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	alertID, _ := ctx.Value("alertID").(uuid.UUID)

	var updatedAlert models.Alert
	if err := json.NewDecoder(r.Body).Decode(&updatedAlert); err != nil {
		logrus.Errorf("Failed to decode request body: %s", err.Error())
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	updatedAlert.Id = &alertID // Ensure the ID matches the request param

	err := services.UpdateAlert(updatedAlert)
	if err != nil {
		logrus.Errorf("Error updating alert: %s", err.Error())
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
	body, _ := json.Marshal(map[string]string{"message": "Alert updated successfully"})
	_, _ = w.Write(body)
}
