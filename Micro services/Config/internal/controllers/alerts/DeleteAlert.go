package collections

import (
	"encoding/json"
	"github.com/gofrs/uuid"
	"github.com/sirupsen/logrus"
	"middleware/example/internal/models"
	services "middleware/example/internal/services/alerts"
	"net/http"
)

// DeleteAlert handles the deletion of an alert
func DeleteAlert(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	alertID, _ := ctx.Value("alertID").(uuid.UUID)

	// Call the service layer to delete the alert
	err := services.DeleteAlert(alertID)
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
	body, _ := json.Marshal(map[string]string{"message": "Alert deleted successfully"})
	_, _ = w.Write(body)
}
