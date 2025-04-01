package events

import (
	"database/sql"
	"github.com/gofrs/uuid"
	"github.com/sirupsen/logrus"
	"middleware/example/internal/models"
	repository "middleware/example/internal/repositories/events"
	"net/http"
)

// GetAllEvents retrieves all events from the repository
func GetAllEvents() ([]*models.Event, error) {
	events, err := repository.GetEvents()
	if err != nil {
		// If no events exist, return an empty list instead of nil
		if err == sql.ErrNoRows {
			return []*models.Event{}, nil
		}
		logrus.Errorf("Error retrieving events: %s", err.Error())
		return nil, &models.CustomError{
			Message: "Something went wrong",
			Code:    http.StatusInternalServerError,
		}
	}

	return events, nil
}

// GetEventByID retrieves an event by ID
func GetEventByID(id string) (*models.Event, error) {
	event, err := repository.GetEventByUID(id)
	if err != nil {
		if err.Error() == sql.ErrNoRows.Error() {
			return nil, &models.CustomError{
				Message: "Event not found",
				Code:    http.StatusNotFound,
			}
		}
		logrus.Errorf("Error retrieving event: %s", err.Error())
		return nil, &models.CustomError{
			Message: "Something went wrong",
			Code:    http.StatusInternalServerError,
		}
	}

	return event, nil
}

// UpdateEvent updates an existing event
func UpdateEvent(event models.Event) error {
	err := repository.UpdateEvent(event)
	if err != nil {
		logrus.Errorf("Error updating event: %s", err.Error())
		return &models.CustomError{
			Message: "Error updating the event",
			Code:    http.StatusInternalServerError,
		}
	}

	// Update resource IDs
	err = repository.UpdateEventResources(event.UID, event.ResourceIDs)
	if err != nil {
		logrus.Errorf("Error updating resource IDs: %s", err.Error())
		return err
	}

	return nil
}

// CreatEvent creates a new event in the database
func CreatEvent(event models.Event) (*models.Event, error) {
	// generating a new ID
	newUUID, err := uuid.NewV4()
	if err != nil {
		logrus.Errorf("Error generating UUID: %s", err.Error())
		return nil, &models.CustomError{
			Message: "Failed to generate unique ID",
			Code:    http.StatusInternalServerError,
		}
	}
	event.Id = &newUUID

	err = repository.CreatEvent(event)
	if err != nil {
		logrus.Errorf("Error creating event: %s", err.Error())
		return nil, &models.CustomError{
			Message: "Error creating event",
			Code:    http.StatusInternalServerError,
		}
	}

	return &event, nil
}
