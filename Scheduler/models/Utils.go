package models

import (
	"fmt"
	"strconv"
	"time"
)

const (
	CONFIG_PATH = "http://localhost:8080/resources/"
)

func UCA_URL(nbWeeks string, resources []Resource) string {
	// Join multiple resource IDs with ","
	IDs := ""
	for _, resource := range resources {
		if IDs != "" {
			IDs += ","
		}
		IDs += strconv.Itoa(resource.UcaID)
	}

	return "https://edt.uca.fr/jsp/custom/modules/plannings/anonymous_cal.jsp?resources=" + IDs +
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
