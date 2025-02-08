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

// UpdateAlert updates an existing alert
func UpdateAlert(alert models.Alert) error {
	err := repository.UpdateAlert(alert)
	if err != nil {
		logrus.Errorf("Error updating alert: %s", err.Error())
		return &models.CustomError{
			Message: "Error updating the alert",
			Code:    http.StatusInternalServerError,
		}
	}
	return nil
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

// DeleteAlert deletes an alert by its ID
func DeleteAlert(id uuid.UUID) error {
	err := repository.DeleteAlertById(id)
	if err != nil {
		logrus.Errorf("Error deleting alert: %s", err.Error())
		return &models.CustomError{
			Message: "Error deleting the alert",
			Code:    http.StatusInternalServerError,
		}
	}
	return nil
}
