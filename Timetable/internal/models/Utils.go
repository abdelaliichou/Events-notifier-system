package models

import (
	"fmt"
	"time"
)

// Constant static values
const (
	AppName     = "ICHOU_GoApp"
	Version     = "1.0.0"
	DB_NAME     = "file:timetable.db"
	CREAT_EVENT = `CREATE TABLE IF NOT EXISTS events (
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

func GetEventChanges(existing *Event, newEvent *Event) map[string]interface{} {
	changes := make(map[string]interface{})

	if existing.Description != newEvent.Description {
		changes["Description"] = map[string]string{
			"old": existing.Description,
			"new": newEvent.Description,
		}
	}
	if existing.Name != newEvent.Name {
		changes["Name"] = map[string]string{
			"old": existing.Name,
			"new": newEvent.Name,
		}
	}
	if !existing.Start.Equal(newEvent.Start) {
		changes["Start"] = map[string]string{
			"old": existing.Start.Format(time.RFC3339),
			"new": newEvent.Start.Format(time.RFC3339),
		}
	}
	if !existing.End.Equal(newEvent.End) {
		changes["End"] = map[string]string{
			"old": existing.End.Format(time.RFC3339),
			"new": newEvent.End.Format(time.RFC3339),
		}
	}
	if existing.Location != newEvent.Location {
		changes["Location"] = map[string]string{
			"old": existing.Location,
			"new": newEvent.Location,
		}
	}
	if !existing.LastUpdate.Equal(newEvent.LastUpdate) {
		changes["LastUpdate"] = map[string]string{
			"old": existing.LastUpdate.Format(time.RFC3339),
			"new": newEvent.LastUpdate.Format(time.RFC3339),
		}
	}

	return changes
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

// TestFunction generates test events
func TestFunction() []Event {
	return []Event{
		{
			UID:         "12345",
			Name:        "Updated Exam ICHOU 10",
			Description: "New description 10",
			Start:       time.Now(),
			End:         time.Now().Add(2 * time.Hour),
			Location:    "New Location 10",
			LastUpdate:  time.Now(),
		},
		{
			UID:         "67890",
			Name:        "New Seminar ICHOU 10",
			Description: "Seminar on AI 10",
			Start:       time.Now().Add(24 * time.Hour),
			End:         time.Now().Add(26 * time.Hour),
			Location:    "HallC 10",
			LastUpdate:  time.Now(),
		},
	}
}
