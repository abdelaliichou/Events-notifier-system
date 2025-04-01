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
		resourceRows, err := db.Query(models.GET_ALL_RESOURCES_OF_EVENT, event.UID)
		if err != nil {
			return nil, err
		}

		var resourceIDs []*uuid.UUID
		for resourceRows.Next() {
			var resourceID uuid.UUID
			if err := resourceRows.Scan(&resourceID); err != nil {
				return nil, err
			}
			resourceIDs = append(resourceIDs, &resourceID)
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

// GetEventByUID helps us to see if Event has been modified
func GetEventByUID(eventUID string) (*models.Event, error) {
	db, err := helpers.OpenDB()
	if err != nil {
		return nil, err
	}
	defer helpers.CloseDB(db)

	// Check if event exists by UID
	query := models.GET_EVENT_BY_UID
	row := db.QueryRow(query, eventUID)

	var newEvent models.Event
	err = row.Scan(
		&newEvent.Id,
		&newEvent.UID,
		&newEvent.Description,
		&newEvent.Name,
		&newEvent.Start,
		&newEvent.End,
		&newEvent.Location,
		&newEvent.LastUpdate,
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

	// Query resource IDs associated with the event
	resourceQuery := models.GET_ALL_RESOURCES_OF_EVENT
	rows, err := db.Query(resourceQuery, eventUID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	// Collect resource IDs
	var resourceIDs []*uuid.UUID
	for rows.Next() {
		var resourceID uuid.UUID
		if err := rows.Scan(&resourceID); err != nil {
			return nil, err
		}
		resourceIDs = append(resourceIDs, &resourceID)
	}

	newEvent.ResourceIDs = resourceIDs
	return &newEvent, nil
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
		event.End, event.Location, event.LastUpdate, event.UID,
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

// UpdateEventResources updates resource IDs linked to an event
func UpdateEventResources(eventUID string, newResourceIDs []*uuid.UUID) error {
	db, err := helpers.OpenDB()
	if err != nil {
		return err
	}
	defer helpers.CloseDB(db)

	// Remove old resource IDs
	_, err = db.Exec(models.DELETE_RESOURCE_IDS, eventUID)
	if err != nil {
		return err
	}

	// Insert new resource IDs
	return InsertEventResources(eventUID, newResourceIDs)
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

	if err != nil {
		return err
	}

	return nil
}

// InsertEventResources inserts resource IDs for a given event UID
func InsertEventResources(eventUID string, resourceIDs []*uuid.UUID) error {
	db, err := helpers.OpenDB()
	if err != nil {
		return err
	}
	defer helpers.CloseDB(db)

	// First, delete existing resource associations for this event
	_, err = db.Exec(models.DELETE_RESOURCE_IDS, eventUID)
	if err != nil {
		return err
	}

	query := models.INSERT_RESOURCE_IDS
	for _, resourceID := range resourceIDs {
		_, err := db.Exec(query, eventUID, resourceID.String())
		if err != nil {
			return err
		}
	}

	return nil
}
