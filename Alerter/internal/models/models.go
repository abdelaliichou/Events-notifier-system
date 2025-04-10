package models

import (
	"github.com/gofrs/uuid"
	"time"
)

type Alert struct {
	Id         *uuid.UUID `json:"id"`
	Email      string     `json:"email"`
	IsAll      bool       `json:"all"`
	ResourceID *uuid.UUID `json:"resourceID"` // Always present, but can be null
}

type Event struct {
	Id          *uuid.UUID   `json:"id"`
	ResourceIDs []*uuid.UUID `json:"resourceIds"`
	UID         string       `json:"uid"`
	Description string       `json:"description"`
	Name        string       `json:"name"`
	Start       time.Time    `json:"start"`
	End         time.Time    `json:"end"`
	Location    string       `json:"location"`
	LastUpdate  time.Time    `json:"lastUpdate"`
}

type EmailContent struct {
	Subject string
	Body    string
}

type MailRequest struct {
	Recipient string `json:"recipient"`
	Subject   string `json:"subject"`
	Content   string `json:"content"`
}
