package main

import (
	"bufio"
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"
)

// to run this file => go run .

func main() {

	FetchingFromConfig()

}

func FetchingFromConfig() {

	// Make the HTTP GET request
	var body []byte
	body = HttpRequest(CONFIG_PATH, false)

	// Parse the JSON data into a slice of Resource structs
	var resources []Resource
	err := json.Unmarshal(body, &resources)
	if err != nil {
		fmt.Println("Error parsing JSON:", err)
		return
	}

	// Print the parsed data
	fmt.Println("\nGetting all resources from Config : ", CONFIG_PATH)
	for _, resource := range resources {
		fmt.Printf("ID: %s, UcaID: %d, Name: %s\n", resource.Id, resource.UcaID, resource.Name)
	}

	// fetching the data from the uca server
	FetchingFromUCA(resources)

}

func FetchingFromUCA(resources []Resource) {

	fmt.Println("\nURL with all resourceIDs : ", UCA_URL("8", resources))

	for _, resource := range resources {

		// filter resources based on the Name field
		// filter()

		var customResources []Resource
		customResources = append(customResources, resource)
		url := UCA_URL("8", customResources)
		fmt.Printf("\nURL of resource %d id : %s\n", resource.UcaID, url)

		// doing the request with this particular resourceID
		uca_resp := HttpRequest(url, true)
		ParsingEvents(uca_resp)

	}
}

func filter() {
	// Filtering recourses
	/*if resource.Name == "M1 Groupe option" {
	//	// || resource.Name == "M1 Groupe 2 langue" {
	}*/
}

func HttpRequest(url string, show bool) []byte {

	// Make the HTTP GET request
	resp, err := http.Get(url)
	if err != nil {
		fmt.Println("Error making the request:", err)
		return []byte("No body exists because of error!")
	}
	defer resp.Body.Close()

	// Check if the response status code is OK (200)
	if resp.StatusCode != http.StatusOK {
		fmt.Println("Error: Received status code", resp.StatusCode)
		return []byte("No body exists because of error!")
	}

	// Read the response body
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		fmt.Println("Error reading the response body:", err)
		return []byte("No body exists because of error!")
	}

	// Print the raw iCalendar data
	if show {
		fmt.Println("\niCalendar response from : ", url)
		fmt.Println(string(body))
	}

	return body
}

func ParsingEvents(data []byte) {
	// create line reader from data
	scanner := bufio.NewScanner(bytes.NewReader(data))

	// create vars
	var eventArray []map[string]string
	currentEvent := map[string]string{}

	currentKey := ""
	currentValue := ""

	inEvent := false

	// inspecting each line
	fmt.Printf("----------------  STARTING PREPARING DATA TO BE PARSED USING SCANNER :")
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
			fmt.Printf("%s\n", scanner.Text())
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
			Start:       startTime,
			End:         endTime,
			LastUpdate:  lastModified,
		})
	}

	// Print the structured events
	fmt.Printf("\n----------------  THE PARSED EVENTS FROM THE CALENDAR RESPONSE : \n")
	for i, event := range structuredEvents {
		fmt.Printf("Event %d:\n", i+1)
		fmt.Printf("  Description: %s\n", event.Description)
		fmt.Printf("  Location: %s\n", event.Location)
		fmt.Printf("  Start: %s\n", event.Start.Format(time.RFC3339))
		fmt.Printf("  End: %s\n", event.End.Format(time.RFC3339))
		fmt.Printf("  Last Update: %s\n", event.LastUpdate.Format(time.RFC3339))
		fmt.Println("-----")
	}
}
