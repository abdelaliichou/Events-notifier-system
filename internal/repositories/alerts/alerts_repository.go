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

	rows, err := db.Query("SELECT * FROM alerts")
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
	query := "SELECT * FROM alerts WHERE id = ?"
	row := db.QueryRow(query, id.String())

	err = row.Scan(&alert.Id, &alert.Email, &alert.IsAll, &alert.ResourceID)
	if err != nil {
		return nil, err // Let service handle the error
	}

	return &alert, nil
}

// UpdateAlert updates an existing alert in the database
func UpdateAlert(alert models.Alert) error {
	db, err := helpers.OpenDB()
	if err != nil {
		return err
	}
	defer helpers.CloseDB(db)

	result, err := db.Exec(`
		UPDATE alerts 
		SET email = ?, is_all = ?, resourceID = ? 
		WHERE id = ?`,
		alert.Email, alert.IsAll, alert.ResourceID, alert.Id.String(),
	)

	// Check if any row was affected (if no rows were deleted, the alert wasn't found)
	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}
	if rowsAffected == 0 {
		return &models.CustomError{
			Message: "Resource not found",
			Code:    401,
		}
	}

	return err
}

// CreateAlert inserts a new alert into the database
func CreateAlert(alert models.Alert) error {
	db, err := helpers.OpenDB()
	if err != nil {
		return err
	}
	defer helpers.CloseDB(db)

	_, err = db.Exec(`
		INSERT INTO alerts (id, email, is_all, resourceID) 
		VALUES (?, ?, ?, ?)`,
		alert.Id.String(), alert.Email, alert.IsAll, alert.ResourceID,
	)

	return err
}

// DeleteAlertById deletes an alert from the database by its ID
func DeleteAlertById(id uuid.UUID) error {
	db, err := helpers.OpenDB()
	if err != nil {
		return err
	}
	defer helpers.CloseDB(db)

	// Delete the alert based on the ID
	result, err := db.Exec("DELETE FROM alerts WHERE id = ?", id.String())
	if err != nil {
		return err
	}

	// Check if any row was affected (if no rows were deleted, the alert wasn't found)
	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}
	if rowsAffected == 0 {
		return &models.CustomError{
			Message: "Alert not found",
			Code:    401,
		}
	}

	return nil
}
