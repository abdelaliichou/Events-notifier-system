package models

// Constant static values
const (
	AppName         = "ICHOU_GoApp"
	Version         = "1.0.0"
	DB_PATH         = "file:config.db"
	RESOURCES_TABLE = `CREATE TABLE IF NOT EXISTS resources (
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
)
