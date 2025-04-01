package models

import (
	"bufio"
	"bytes"
	"fmt"
	"github.com/gofrs/uuid"
	"strconv"
	"strings"
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

func ParsingEvents(data []byte, ResourceID *uuid.UUID, show bool) []Event {
	// create line reader from data
	scanner := bufio.NewScanner(bytes.NewReader(data))

	// create vars
	var eventArray []map[string]string
	currentEvent := map[string]string{}

	currentKey := ""
	currentValue := ""

	inEvent := false

	// inspecting each line
	if show {
		fmt.Printf("----------------  STARTING PREPARING DATA TO BE PARSED USING SCANNER :")
	}
	fmt.Printf("\n")
	for scanner.Scan() {

		// ignore calendar lines
		if !inEvent && scanner.Text() != "BEGIN:VEVENT" {
			continue
		}

		// if new event, go to next line
		if scanner.Text() == "BEGIN:VEVENT" {
			inEvent = true
			currentEvent = map[string]string{}
			continue
		}

		if scanner.Text() == "END:VEVENT" {
			inEvent = false
			eventArray = append(eventArray, currentEvent)
			continue
		}

		if strings.HasPrefix(scanner.Text(), " ") {
			currentEvent[currentKey] += scanner.Text()
		} else {
			// split scan
			if show {
				fmt.Printf("%s\n", scanner.Text())
			}
			splitted := strings.SplitN(scanner.Text(), ":", 2)
			currentKey = splitted[0]
			currentValue = splitted[1]

			// store current event attribute
			currentEvent[currentKey] = currentValue
		}
	}

	var structuredEvents []Event
	for _, event := range eventArray {

		startTime, _ := time.Parse("20060102T150405Z", event["DTSTART"])
		endTime, _ := time.Parse("20060102T150405Z", event["DTEND"])
		lastModified, _ := time.Parse("20060102T150405Z", event["LAST-MODIFIED"])

		structuredEvents = append(structuredEvents, Event{
			Description: event["DESCRIPTION"],
			Location:    event["LOCATION"],
			UID:         event["UID"],
			ResourceIDs: []*uuid.UUID{ResourceID},
			Start:       startTime,
			Name:        event["SUMMARY"],
			End:         endTime,
			LastUpdate:  lastModified,
		})
	}

	// Print the structured events
	if show {
		fmt.Printf("\n----------------  THE PARSED EVENTS FROM THE CALENDAR RESPONSE : \n")
		DisplayEvents(structuredEvents)
	}

	return structuredEvents
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
