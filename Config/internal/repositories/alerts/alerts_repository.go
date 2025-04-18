package alerts

import (
	"github.com/gofrs/uuid"
	"middleware/example/internal/helpers"
	"middleware/example/internal/models"
)

// GetAllAlerts retrieves all alerts from the database
func GetAllAlerts() ([]*models.Alert, error) {
	db, err := helpers.OpenDB()
	if err != nil {
		return nil, err
	}
	// Ensure the db connection is closed when the function returns
	defer helpers.CloseDB(db)

	rows, err := db.Query(models.GET_ALL_ALERTS)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var alerts []*models.Alert

	for rows.Next() {
		var alert models.Alert
		err := rows.Scan(&alert.Id, &alert.Email, &alert.IsAll, &alert.ResourceID)
		if err != nil {
			return nil, err
		}
		alerts = append(alerts, &alert)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return alerts, nil
}

// GetAlertById fetches an alert by its ID from the database
func GetAlertById(id uuid.UUID) (*models.Alert, error) {
	db, err := helpers.OpenDB()
	if err != nil {
		return nil, err
	}
	defer helpers.CloseDB(db)

	var alert models.Alert
	query := models.GET_ALERT
	row := db.QueryRow(query, id.String())

	err = row.Scan(&alert.Id, &alert.Email, &alert.IsAll, &alert.ResourceID)
	if err != nil {
		return nil, err // Let service handle the error
	}

	return &alert, nil
}

// CreateAlert inserts a new alert into the database
func CreateAlert(alert models.Alert) error {
	db, err := helpers.OpenDB()
	if err != nil {
		return err
	}
	defer helpers.CloseDB(db)

	_, err = db.Exec(models.CREAT_ALERT,
		alert.Id.String(), alert.Email, alert.IsAll, alert.ResourceID)

	return err
}
