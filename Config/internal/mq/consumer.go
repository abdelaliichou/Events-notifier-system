package mq

import (
	"encoding/json"
	"fmt"
	"github.com/nats-io/nats.go"
	"io"
	"log"
	"middleware/example/internal/models"
	"net/http"
)

var jsc nats.JetStreamContext

func ensureStreamExists() error {
	_, err := jsc.AddStream(&nats.StreamConfig{
		Name:     "ALERTS",             // Stream Name
		Subjects: []string{"ALERTS.>"}, // Capture all subjects under ALERTS.
	})
	if err != nil && err != nats.ErrStreamNameAlreadyInUse {
		return err
	}
	return nil
}

// StartStreamConsumer Starts Consumer for the Event Stream coming from the scheduler
func StartStreamConsumer() {
	// Connect to NATS server
	var err error
	nc, err := nats.Connect(nats.DefaultURL)
	if err != nil {
		log.Fatal("❌ Failed to connect to ALERTS NATS:", err)
	}

	// Initialize JetStream context
	jsc, err = nc.JetStream()
	if err != nil {
		log.Fatal("❌ Failed to initialize ALERTS JetStream:", err)
	}

	fmt.Println("✅ Connected to ALERTS NATS and JetStream initialized")

	// Ensure the ALERTS stream exists
	if err := ensureStreamExists(); err != nil {
		log.Fatal("❌ Failed to create or verify ALERTS stream:", err)
	}

	// Subscribe to ALERTS.stream using JetStream
	_, err = jsc.Subscribe("ALERTS.stream", func(m *nats.Msg) {

		// Expecting and parsing the list of event IDs
		var eventIDs []string
		if err := json.Unmarshal(m.Data, &eventIDs); err != nil {
			log.Println("❌ Failed to decode event list:", err)
			return
		}

		var events []models.Event
		var event models.Event
		// Process each event
		for _, eventID := range eventIDs {
			fmt.Println("Received Event from Timetable:", eventID)
			response := HttpRequest("http://localhost:8090/events/"+eventID, false)

			err := json.Unmarshal(response, &event)
			if err != nil {
				return
			}

			events = append(events, event)
		}

		models.DisplayEvents(events)

		// Acknowledge the message after processing
		m.Ack()

	}, nats.Durable("event-consumer"), nats.ManualAck())

	if err != nil {
		log.Fatal("❌ Failed to subscribe to ALERTS.stream:", err)
	}

	log.Println("✅ ALERTS Stream Consumer is running...")

	// Prevent program from exiting
	select {}
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
		fmt.Println("TimeTable Http response from : ", url)
		fmt.Println(string(body))
	}

	return body
}
