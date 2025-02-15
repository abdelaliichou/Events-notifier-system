package main

import (
	"bufio"
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"strings"
)

// to run this file => go run .

func main() {
	FetchingFromConfig()
}

func FetchingFromConfig() {

	// Make the HTTP GET request
	resp, err := http.Get(CONFIG_PATH)
	if err != nil {
		fmt.Println("Error making the request:", err)
		return
	}
	defer resp.Body.Close() // Ensure the response body is closed

	// Read the response body
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		fmt.Println("Error reading the response body:", err)
		return
	}

	// Parse the JSON data into a slice of Resource structs
	var resources []Resource
	err = json.Unmarshal(body, &resources)
	if err != nil {
		fmt.Println("Error parsing JSON:", err)
		return
	}

	// Print the parsed data
	for _, resource := range resources {
		fmt.Printf("ID: %s, UcaID: %d, Name: %s\n", resource.Id, resource.UcaID, resource.Name)
	}

	// fetching the data from the uca server
	FetchingFromUCA(resources)
}

func FetchingFromUCA(resources []Resource) {

	fmt.Println(UCA_URL("8", resources))

	url := UCA_URL("8", resources)

	// Retrieve data from Resources
	resp, err := http.Get(url)
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println(resp)
}

func test2() {

	// Retrieve data from Resources
	resp, err := http.Get("http://localhost:8080/resources/")

	// resp, err := http.Get("https://edt.uca.fr/jsp/custom/modules/plannings/anonymous_cal.jsp?resources=13295&projectId=2&calType=ical&nbWeeks=8&displayConfigId=128")
	// TODO manage error
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println(resp)

	// Read all and store in value
	rawData, err := io.ReadAll(resp.Body)
	// TODO manage error
	fmt.Println(resp)

	// Create a line-reader from data
	scanner := bufio.NewScanner(bytes.NewReader(rawData))

	// Create vars
	var eventArray []map[string]string
	fmt.Println(eventArray)
	currentEvent := map[string]string{}

	currentKey := ""
	currentValue := ""

	inEvent := false

	// Inspect each line
	for scanner.Scan() {
		// Ignore calendar lines
		if !inEvent && scanner.Text() != "BEGIN:VEVENT" {
			continue
		}
		// If new event, go to next line
		if scanner.Text() == "BEGIN:VEVENT" {
			inEvent = true
			continue
		}

		// TODO if end event

		// TODO if multi-line data

		// Split scan
		fmt.Println(scanner.Text())
		splitted := strings.SplitN(scanner.Text(), ":", 2)
		currentKey = splitted[0]
		currentValue = splitted[1]

		// Store current event attribute
		currentEvent[currentKey] = currentValue
	}

	// TODO Transform to proper custom object

	// TODO parse to JSON and display
}
