package events

import (
	"database/sql"
	"github.com/gofrs/uuid"
	"middleware/example/internal/helpers"
	"middleware/example/internal/models"
	"net/http"
)

// GetEvents retrieves all events from the database
func GetEvents() ([]*models.Event, error) {
	db, err := helpers.OpenDB()
	if err != nil {
		return nil, err
	}
	defer helpers.CloseDB(db)

	// Fetching events
	rows, err := db.Query(models.GET_ALL_EVENTS)
	if err != nil {
		// If no rows exist, return sql.ErrNoRows so service can handle it
		if err == sql.ErrNoRows {
			return nil, sql.ErrNoRows
		}
		return nil, err
	}
	defer rows.Close()

	var events []*models.Event
	for rows.Next() {
		var event models.Event
		err := rows.Scan(
			&event.Id,
			&event.UID,
			&event.Description,
			&event.Name,
			&event.Start,
			&event.End,
			&event.Location,
			&event.LastUpdate,
		)
		if err != nil {
			return nil, err
		}

		// Fetch resource IDs related to this event
		resourceRows, err := db.Query(models.GET_ALL_RESOURCES_OF_EVENT, event.Id)
		if err != nil {
			return nil, err
		}

		var resourceIDs []uuid.UUID
		for resourceRows.Next() {
			var resourceID uuid.UUID
			if err := resourceRows.Scan(&resourceID); err != nil {
				return nil, err
			}
			resourceIDs = append(resourceIDs, resourceID)
		}
		resourceRows.Close()

		event.ResourceIDs = resourceIDs
		events = append(events, &event)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	// If no events found, return sql.ErrNoRows
	if len(events) == 0 {
		return nil, sql.ErrNoRows
	}

	return events, nil
}

// GetEventById fetches an event by its ID from the database
func GetEventById(eventID uuid.UUID) (*models.Event, error) {
	db, err := helpers.OpenDB()
	if err != nil {
		return nil, err
	}
	defer helpers.CloseDB(db)

	// Query to fetch the event details
	query := models.GET_EVENT_BY_ID
	row := db.QueryRow(query, eventID)

	var event models.Event
	err = row.Scan(
		&event.Id,
		&event.UID,
		&event.Description,
		&event.Name,
		&event.Start,
		&event.End,
		&event.Location,
		&event.LastUpdate,
	)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, &models.CustomError{
				Message: "Event not found",
				Code:    http.StatusNotFound,
			}
		}
		return nil, err
	}

	// Fetch resource IDs related to this event
	resourceQuery := models.GET_ALL_RESOURCES_OF_EVENT
	resourceRows, err := db.Query(resourceQuery, event.Id)
	if err != nil {
		return nil, err
	}
	defer resourceRows.Close()

	var resourceIDs []uuid.UUID
	for resourceRows.Next() {
		var resourceID uuid.UUID
		if err := resourceRows.Scan(&resourceID); err != nil {
			return nil, err
		}
		resourceIDs = append(resourceIDs, resourceID)
	}

	event.ResourceIDs = resourceIDs
	return &event, nil
}

// UpdateEvent updates an existing alert in the database
func UpdateEvent(event models.Event) error {
	db, err := helpers.OpenDB()
	if err != nil {
		return err
	}
	defer helpers.CloseDB(db)

	result, err := db.Exec(models.UPDATE_EVENT,
		event.Description, event.Name, event.Start,
		event.End, event.Location, event.LastUpdate, event.Id.String(),
	)

	// Check if any row was affected (if no rows were deleted, the event wasn't found)
	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}
	if rowsAffected == 0 {
		return &models.CustomError{
			Message: "Event not found",
			Code:    401,
		}
	}

	return nil
}

// CreatEvent inserts a new event into the database
func CreatEvent(event models.Event) error {
	db, err := helpers.OpenDB()
	if err != nil {
		return err
	}
	defer helpers.CloseDB(db)

	_, err = db.Exec(models.POST_EVENT,
		event.Id.String(), event.UID, event.Description,
		event.Name, event.Start, event.End, event.Location,
		event.LastUpdate,
	)

	return err
}

// DeleteEventById deletes an event from the database by its ID
func DeleteEventById(id uuid.UUID) error {
	db, err := helpers.OpenDB()
	if err != nil {
		return err
	}
	defer helpers.CloseDB(db)

	// Delete the event based on the ID
	result, err := db.Exec(models.DELETE_EVENT, id.String())
	if err != nil {
		return err
	}

	// Check if any row was affected (if no rows were deleted, the event wasn't found)
	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}
	if rowsAffected == 0 {
		return &models.CustomError{
			Message: "Event not found",
			Code:    401,
		}
	}

	return nil
}
