package collections

import (
	"encoding/json"
	"github.com/gofrs/uuid"
	"github.com/sirupsen/logrus"
	"middleware/example/internal/models"
	services "middleware/example/internal/services/alerts"
	"net/http"
)

// GetAlert
// @Tags         alerts
// @Summary      Get an alert by ID
// @Description  Retrieve a specific alert using its unique ID
// @Param        id path string true "Alert UUID"
// @Success      200 {object} models.Alert
// @Failure      404 "Alert not found"
// @Failure      500 "Internal server error"
// @Router       /alerts/{id} [get]
func GetAlert(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	alertID, _ := ctx.Value("alertID").(uuid.UUID)

	// Fetch the alert from the service layer
	alert, err := services.GetAlertById(alertID)
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

	// Return the alert details as JSON
	w.WriteHeader(http.StatusOK)
	body, _ := json.Marshal(alert)
	_, _ = w.Write(body)
}
