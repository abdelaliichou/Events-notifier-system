package models

import (
	"github.com/gofrs/uuid"
	"time"
)

// Timetable represent un employ du temps
type Event struct {
	Id          *uuid.UUID  `json:"id"`
	ResourceIDs []uuid.UUID `json:"resourceIds"`
	UID         string      `json:"uid"`
	Description string      `json:"description"`
	Name        string      `json:"name"`
	Start       time.Time   `json:"start"`
	End         time.Time   `json:"end"`
	Location    string      `json:"location"`
	LastUpdate  time.Time   `json:"lastUpdate"`
}
