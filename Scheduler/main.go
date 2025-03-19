package main

import (
	"Scheduler/models"
	"Scheduler/mq"
	"Scheduler/webservice"
	"bufio"
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"github.com/gofrs/uuid"
	"github.com/zhashkevych/scheduler"
	"log"
	"os"
	"os/signal"
	"strings"
	"time"
)

// to run this file => go run .

func main() {

	// Starting NATS JetStream
	mq.InitStream()

	schedulingFunctionCall()

}

func schedulingFunctionCall() {

	ctx := context.Background()
	sc := scheduler.NewScheduler()
	sc.Add(ctx, FetchingFromConfig, time.Second*10)

	// cette partie sert à maintenir le programme en vie, tant qu'il n'a pas reçu de signal os.Interrupt
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt)

	<-quit
	sc.Stop()
}

func FetchingFromConfig(_ context.Context) {

	// Make the HTTP GET request
	var body []byte
	body = webservice.HttpRequest(models.CONFIG_PATH, false)

	// Parse the JSON data into a slice of Resource structs
	var resources []models.Resource
	err := json.Unmarshal(body, &resources)
	if err != nil {
		fmt.Println("Error parsing JSON:", err)
		return
	}

	// Print the parsed data
	fmt.Println("\nGetting all resources from Config : ", models.CONFIG_PATH)
	for _, resource := range resources {
		fmt.Printf("ID: %s, UcaID: %d, Name: %s\n", resource.Id, resource.UcaID, resource.Name)
	}

	// fetching the data from the uca server
	FetchingFromUCA(resources, false)

}

func FetchingFromUCA(resources []models.Resource, show bool) {

	if show {
		fmt.Println("\nURL with all resourceIDs : ", models.UCA_URL("8", resources))
	}

	var allEvents []models.Event
	for _, resource := range resources {

		var customResources []models.Resource
		customResources = append(customResources, resource)
		url := models.UCA_URL("8", customResources)
		if show {
			fmt.Printf("\nURL of resource %d id : %s\n", resource.UcaID, url)
		}

		// doing the request with this particular resourceID
		ucaResp := webservice.HttpRequest(url, false)
		specificResourceEvents := ParsingEvents(ucaResp, resource.Id, false)
		allEvents = append(allEvents, specificResourceEvents...)

	}

	fmt.Println("ALL EVENTS FROM ALL RESOURCES :")
	models.DisplayEvents(allEvents)

	// Group events by UID & merge resource IDs
	groupedEvents := models.GroupEventsByUID(allEvents)

	fmt.Println("ALL EVENTS FROM ALL RESOURCES AFTER GROUPING :")
	models.DisplayEvents(groupedEvents)

	// Send into MQ
	sendEventsToMQ(groupedEvents)
}

func ParsingEvents(data []byte, ResourceID *uuid.UUID, show bool) []models.Event {
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

	var structuredEvents []models.Event
	for _, event := range eventArray {

		startTime, _ := time.Parse("20060102T150405Z", event["DTSTART"])
		endTime, _ := time.Parse("20060102T150405Z", event["DTEND"])
		lastModified, _ := time.Parse("20060102T150405Z", event["LAST-MODIFIED"])

		structuredEvents = append(structuredEvents, models.Event{
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
		models.DisplayEvents(structuredEvents)
	}

	return structuredEvents
}

// sendEventsToMQ will send structured events to our producer as a stream to MQ
func sendEventsToMQ(structuredEvents []models.Event) {

	err := mq.PublishEventsAsStream(structuredEvents)
	if err != nil {
		log.Println("Error sending events to MQ:", err)
		return
	}

	fmt.Println("All events sent successfully to MQ")
}
