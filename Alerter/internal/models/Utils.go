package models

import (
	"fmt"
	"github.com/gofrs/uuid"
	"time"
)

// Constant static values
const (
	MailAPI          = "https://mail-api.edu.forestier.re/mail"
	CONFIG_ALERT_URL = "http://localhost:8080/alerts/"
)

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

func DisplayAlerts(alerts []Alert) {
	for i, alert := range alerts {
		fmt.Printf("Alert %d:\n", i+1)
		fmt.Printf("  Id: %s\n", alert.Id)
		fmt.Printf("  isAll: %t\n", alert.IsAll)
		fmt.Printf("  Email ID: %s\n", alert.Email)
		fmt.Printf("  ResourceID: %s\n", alert.ResourceID)
		fmt.Println("-----")
	}
}

func CheckAlertResourceInEventResources(alertResourceID *uuid.UUID, events []Event) []Event {
	var subscribeEvents []Event
	for _, event := range events {
		for _, res := range event.ResourceIDs {
			if res != nil && res.String() == alertResourceID.String() {
				subscribeEvents = append(subscribeEvents, event)
				break
			}
		}
	}
	return subscribeEvents
}
