package collections

import (
	"encoding/json"
	"github.com/sirupsen/logrus"
	"middleware/example/internal/models"
	services "middleware/example/internal/services/alerts"
	"net/http"
)

// GetAlerts
// @Tags         alerts
// @Summary      Get all alerts
// @Description  Retrieve a list of all alerts
// @Success      200 {array} models.Alert
// @Failure      500 "Internal server error"
// @Router       /alerts [get]
func GetAlerts(w http.ResponseWriter, r *http.Request) {
	// Fetch the alerts from the service layer
	// we don't use context because we don't need IDs, we all the alerts
	alerts, err := services.GetAlerts()
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

	// Return the list of alerts as JSON
	w.WriteHeader(http.StatusOK)
	body, _ := json.Marshal(alerts)
	_, _ = w.Write(body)
}
