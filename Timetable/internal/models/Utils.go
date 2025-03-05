package models

import (
	"github.com/gofrs/uuid"
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
							uid TEXT NOT NULL,
							description TEXT NOT NULL,
							name TEXT NOT NULL,
							start DATETIME NOT NULL,
							end DATETIME NOT NULL,
							location TEXT NOT NULL,
							last_update DATETIME NOT NULL
						);`
	// this one if to represent list of resources blongs to the same event
	CREAT_RESOURCE = `CREATE TABLE IF NOT EXISTS event_resources (
							event_id TEXT NOT NULL,
							resource_id TEXT NOT NULL,
							PRIMARY KEY (event_id, resource_id),
							FOREIGN KEY (event_id) REFERENCES events(id) ON DELETE CASCADE,
							FOREIGN KEY (resource_id) REFERENCES resources(id) ON DELETE CASCADE
						);`
	GET_ALL_EVENTS             = "SELECT * FROM events"
	GET_EVENT_BY_ID            = "SELECT * FROM events WHERE id = ?"
	GET_EVENT_BY_UID           = "SELECT * FROM events WHERE uid = ?"
	GET_ALL_RESOURCES_OF_EVENT = "SELECT resource_id FROM event_resources WHERE event_id = ?"
	UPDATE_EVENT               = `UPDATE events 
								  SET description = ?, name = ?, start = ?, end = ?, location = ?, last_update = ?
								  WHERE uid = ?`
	POST_EVENT = `INSERT INTO events (id, uid, description, name, start, end, location, last_update) 
				  VALUES (?, ?, ?, ?, ?, ?, ?, ?)`
	DELETE_EVENT = "DELETE FROM events WHERE id = ?"
)

func IsEventModified(existing *Event, newEvent *Event) bool {
	return existing.Description != newEvent.Description ||
		existing.Name != newEvent.Name ||
		!existing.Start.Equal(newEvent.Start) ||
		!existing.End.Equal(newEvent.End) ||
		existing.Location != newEvent.Location ||
		!existing.LastUpdate.Equal(newEvent.LastUpdate)
	// !IsResourcesEqual(existing.ResourceIDs, newEvent.ResourceIDs)
}

func IsResourcesEqual(existing []*uuid.UUID, new []*uuid.UUID) bool {
	if len(existing) != len(new) {
		return false
	}

	// Maps for quick comparison, charge ids in the map and check if all other ids exits in this map
	existingMap := make(map[uuid.UUID]bool)
	for _, id := range existing {
		existingMap[*id] = true
	}
	for _, id := range new {
		if !existingMap[*id] {
			return false
		}
	}

	return true
}
