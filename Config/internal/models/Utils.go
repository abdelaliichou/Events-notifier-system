package models

import (
	"fmt"
	"strings"
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
	DB_PATH            = "file:config.db"
	RESOURCES_TABLE    = `CREATE TABLE IF NOT EXISTS resources (
							id TEXT PRIMARY KEY NOT NULL UNIQUE,
							ucaID INTEGER NOT NULL,
							name TEXT NOT NULL
						);`
	ALERTS_TABLE = `CREATE TABLE IF NOT EXISTS alerts (
							id TEXT PRIMARY KEY NOT NULL UNIQUE,
							email TEXT NOT NULL,
							is_all BOOLEAN NOT NULL,
							resourceID TEXT NULL,
							FOREIGN KEY (resourceID) REFERENCES resources(id) ON DELETE SET NULL
						);`
	UPDATE_ALERT = `UPDATE alerts 
					SET email = ?, is_all = ?, resourceID = ? 
					WHERE id = ?`
	CREAT_ALERT = `INSERT INTO alerts (id, email, is_all, resourceID) 
					VALUES (?, ?, ?, ?)`
	GET_ALERT_FOR_EVENT = "SELECT * FROM alerts WHERE resourceID = ? OR is_all = TRUE"
	DELETE_ALERT        = "DELETE FROM alerts WHERE id = ?"
	GET_ALL_ALERTS      = "SELECT * FROM alerts"
	GET_ALERT           = "SELECT * FROM alerts WHERE id = ?"
	GET_ALL_RESOURCES   = "SELECT * FROM resources"
	GET_RESOURCE        = "SELECT * FROM resources WHERE id=?"
	CREAT_RESOURCE      = "INSERT INTO resources (id, ucaID, name) VALUES (?, ?, ?)"
	UPDATE_RESOURCE     = "UPDATE resources SET ucaID=?, name=? WHERE id=?"
	DELETE_RESOURCE     = "DELETE FROM resources WHERE id=?"
	MailAPI             = "https://mail-api.edu.forestier.re/mail"
)

// Function to generate calendar URL with multiple resource IDs

func CalendarURL(nbWeeks string, RESOURCE_ID ...string) string {
	// Join multiple resource IDs with ","
	joinedResources := strings.Join(RESOURCE_ID, ",")

	return "https://edt.uca.fr/jsp/custom/modules/plannings/anonymous_cal.jsp?resources=" + joinedResources +
		"&projectId=2&calType=ical&" + nbWeeks + "=8&displayConfigId=128"
}

func DisplayEvents(events []Event) {
	for i, event := range events {
		fmt.Printf("Event %d:\n", i+1)
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
}

func DisplayAlerts(alerts []*Alert) {
	for i, alert := range alerts {
		fmt.Printf("Alert %d:\n", i+1)
		fmt.Printf("  Id: %s\n", alert.Id)
		fmt.Printf("  isAll: %t\n", alert.IsAll)
		fmt.Printf("  Email ID: %s\n", alert.Email)
		fmt.Printf("  ResourceID: %s\n", alert.ResourceID)
		fmt.Println("-----")
	}
}
