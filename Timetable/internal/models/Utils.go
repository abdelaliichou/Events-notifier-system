package models

import (
	"fmt"
	"github.com/gofrs/uuid"
	"time"
)

// Constant static values
const (
	AppName            = "ICHOU_GoApp"
	Version            = "1.0.0"
	M1_GROUPE_1_lANGUE = "13295"
	M1_GROUPE_2_lANGUE = "13345"
	M1_GROUPE_3_lANGUE = "13397"
	M1_GROUPE_1_OPTION = "7224"
	M1_GROUPE_2_OPTION = "7225"
	M1_GROUPE_3_OPTION = "62962"
	M1_GROUPE_OPTION   = "62090"
	M1_TUTORAT_L2      = "56529"
	DB_NAME            = "file:timetable.db"
	CREAT_EVENT        = `CREATE TABLE IF NOT EXISTS events (
						  id TEXT PRIMARY KEY NOT NULL UNIQUE,
						  uid TEXT NOT NULL UNIQUE, -- UID is unique
						  description TEXT NOT NULL,
						  name TEXT NOT NULL,
						  start DATETIME NOT NULL,
						  end DATETIME NOT NULL,
						  location TEXT NOT NULL,
						  last_update DATETIME NOT NULL
						 );`
	// this one if to represent list of resources blongs to the same event
	CREAT_RESOURCE = `CREATE TABLE IF NOT EXISTS event_resources (
					      event_uid TEXT NOT NULL,
					      resource_id TEXT NOT NULL,
					      PRIMARY KEY (event_uid, resource_id),
					      FOREIGN KEY (event_uid) REFERENCES events(uid) ON DELETE CASCADE,
					      FOREIGN KEY (resource_id) REFERENCES resources(id) ON DELETE CASCADE
				         );`
	INSERT_RESOURCE_IDS = `INSERT INTO event_resources (event_uid, resource_id) VALUES (?, ?)
             					  ON CONFLICT(event_uid, resource_id) DO NOTHING`
	DELETE_RESOURCE_IDS        = `DELETE FROM event_resources WHERE event_uid = ?`
	GET_ALL_EVENTS             = "SELECT * FROM events"
	GET_EVENT_BY_ID            = "SELECT * FROM events WHERE id = ?"
	GET_EVENT_BY_UID           = "SELECT * FROM events WHERE uid = ?"
	GET_ALL_RESOURCES_OF_EVENT = "SELECT resource_id FROM event_resources WHERE event_uid = ?"
	UPDATE_EVENT               = `UPDATE events 
								  SET description = ?, name = ?, start = ?, end = ?, location = ?, last_update = ?
								  WHERE uid = ?`
	POST_EVENT = `INSERT INTO events (id, uid, description, name, start, end, location, last_update) 
				  VALUES (?, ?, ?, ?, ?, ?, ?, ?)`
	DELETE_EVENT = "DELETE FROM events WHERE id = ?"
)

func IsEventModified(existing *Event, newEvent *Event) bool {
	return IsResourcesModified(existing.ResourceIDs, newEvent.ResourceIDs) ||
		existing.Description != newEvent.Description ||
		existing.Name != newEvent.Name ||
		!existing.Start.Equal(newEvent.Start) ||
		!existing.End.Equal(newEvent.End) ||
		existing.Location != newEvent.Location ||
		!existing.LastUpdate.Equal(newEvent.LastUpdate)
}

// IsResourcesModified checks if two slices of UUIDs are equal (order does not matter)
func IsResourcesModified(existingResources []*uuid.UUID, newResources []*uuid.UUID) bool {
	if len(existingResources) != len(newResources) {
		return true
	}

	// Convert both lists to UUID value sets for comparison
	existingSet := make(map[uuid.UUID]bool)
	for _, id := range existingResources {
		if id != nil {
			existingSet[*id] = true
		}
	}

	for _, id := range newResources {
		if id == nil || !existingSet[*id] {
			return true
		}
	}

	return false
}

func DisplayEvents(event Event, i int, all bool) {
	if all {
		fmt.Printf("Event %d:\n", i+1)
	}
	fmt.Printf("  Description: %s\n", event.Description)
	fmt.Printf("  Location: %s\n", event.Location)
	fmt.Printf("  RESOURCES ID: %s\n", event.ResourceIDs)
	fmt.Printf("  UID: %s\n", event.UID)
	fmt.Printf("  NAME: %s\n", event.Name)
	fmt.Printf("  Start: %s\n", event.Start.Format(time.RFC3339))
	fmt.Printf("  End: %s\n", event.End.Format(time.RFC3339))
	fmt.Printf("  Last Update: %s\n", event.LastUpdate.Format(time.RFC3339))
	fmt.Println("-----")
}
