package models

import (
	"fmt"
	"github.com/gofrs/uuid"
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

// GroupEventsByUID combine different resourceIDs events in single event the resourcesID[] contains all the resources of those events
func GroupEventsByUID(events []Event) []Event {
	groupedEvents := make(map[string]Event)

	for _, event := range events {
		if len(event.ResourceIDs) == 0 {
			continue // Skip events without a resource ID
		}

		// Check if the event UID already exists
		if existingEvent, found := groupedEvents[event.UID]; found {
			// Add his resourceID to the existing resources
			existingEvent.ResourceIDs = MergeUnique(existingEvent.ResourceIDs, event.ResourceIDs)
			groupedEvents[event.UID] = existingEvent
		} else {
			// Add new event to the map
			groupedEvents[event.UID] = event
		}
	}

	// Convert the map to a slice of events
	var mergedEvents []Event
	for _, event := range groupedEvents {
		mergedEvents = append(mergedEvents, event)
	}

	return mergedEvents
}

// MergeUnique merge unique resource IDs
func MergeUnique(existing []*uuid.UUID, newIDs []*uuid.UUID) []*uuid.UUID {
	resourceSet := make(map[uuid.UUID]struct{})

	// Add existing resource IDs to the set
	for _, id := range existing {
		resourceSet[*id] = struct{}{}
	}

	// Add new resource ID if not already present
	for _, id := range newIDs {
		if _, exists := resourceSet[*id]; !exists {
			existing = append(existing, id)
			resourceSet[*id] = struct{}{}
		}
	}

	return existing
}
