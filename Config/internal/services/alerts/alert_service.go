package alerts

import (
	"database/sql"
	"github.com/gofrs/uuid"
	"github.com/sirupsen/logrus"
	"middleware/example/internal/models"
	repository "middleware/example/internal/repositories/alerts"
	"net/http"
)

// GetAlerts retrieves all alerts from the database
func GetAlerts() ([]*models.Alert, error) {
	alerts, err := repository.GetAllAlerts()
	if err != nil {
		logrus.Errorf("Error retrieving alerts: %s", err.Error())
		return nil, &models.CustomError{
			Message: "Something went wrong",
			Code:    http.StatusInternalServerError,
		}
	}

	return alerts, nil
}

// GetAlertById retrieves an alert by ID, handling errors properly
func GetAlertById(id uuid.UUID) (*models.Alert, error) {
	alert, err := repository.GetAlertById(id)
	if err != nil {
		if err.Error() == sql.ErrNoRows.Error() {
			return nil, &models.CustomError{
				Message: "Alert not found",
				Code:    http.StatusNotFound,
			}
		}
		logrus.Errorf("Error retrieving alert: %s", err.Error())
		return nil, &models.CustomError{
			Message: "Something went wrong",
			Code:    http.StatusInternalServerError,
		}
	}

	return alert, nil
}

// CreateAlert creates a new alert in the database
func CreateAlert(alert models.Alert) (*models.Alert, error) {
	// generating a new ID
	newUUID, err := uuid.NewV4()
	if err != nil {
		logrus.Errorf("Error generating UUID: %s", err.Error())
		return nil, &models.CustomError{
			Message: "Failed to generate unique ID",
			Code:    http.StatusInternalServerError,
		}
	}
	alert.Id = &newUUID

	err = repository.CreateAlert(alert)
	if err != nil {
		logrus.Errorf("Error creating alert: %s", err.Error())
		return nil, &models.CustomError{
			Message: "Error creating alert",
			Code:    http.StatusInternalServerError,
		}
	}

	return &alert, nil
}
