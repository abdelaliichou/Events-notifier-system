package models

import "strings"

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
	DELETE_ALERT      = "DELETE FROM alerts WHERE id = ?"
	GET_ALL_ALERTS    = "SELECT * FROM alerts"
	GET_ALERT         = "SELECT * FROM alerts WHERE id = ?"
	GET_ALL_RESOURCES = "SELECT * FROM resources"
	GET_RESOURCE      = "SELECT * FROM resources WHERE id=?"
	CREAT_RESOURCE    = "INSERT INTO resources (id, ucaID, name) VALUES (?, ?, ?)"
	UPDATE_RESOURCE   = "UPDATE resources SET ucaID=?, name=? WHERE id=?"
	DELETE_RESOURCE   = "DELETE FROM resources WHERE id=?"
)

// Function to generate calendar URL with multiple resource IDs

func CalendarURL(nbWeeks string, RESOURCE_ID ...string) string {
	// Join multiple resource IDs with ","
	joinedResources := strings.Join(RESOURCE_ID, ",")

	return "https://edt.uca.fr/jsp/custom/modules/plannings/anonymous_cal.jsp?resources=" + joinedResources +
		"&projectId=2&calType=ical&" + nbWeeks + "=8&displayConfigId=128"
}
