package helpers

import (
	"database/sql"
	_ "github.com/mattn/go-sqlite3"
	"github.com/sirupsen/logrus"
	"middleware/example/internal/models"
)

func OpenDB() (*sql.DB, error) {
	db, err := sql.Open("sqlite3", models.DB_NAME)
	if err != nil {
		db.SetMaxOpenConns(1)
	}
	return db, err
}
func CloseDB(db *sql.DB) {
	err := db.Close()
	if err != nil {
		logrus.Errorf("error closing db : %s", err.Error())
	}
}
