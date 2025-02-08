package models

import (
	"github.com/gofrs/uuid"
)

type Collection struct {
	Id      *uuid.UUID `json:"id"`
	Content string     `json:"content"`
}

// Resource represent un employ du temps
type Resource struct {
	Id    *uuid.UUID `json:"id"`
	UcaID int        `json:"ucaID"`
	Name  string     `json:"name"`
}
type Alert struct {
	Id    *uuid.UUID `json:"id"`
	Email string     `json:"email"`
	All   bool       `json:"all"`
	//	ResourceID *uuid.UUID `json:"resourceID,omitempty"` // UUID ou NULL (nullable)
	ResourceID *uuid.UUID `json:"resourceID"` // Always present, but can be null
}
