package mq

import (
	"encoding/json"
	"fmt"
	"github.com/nats-io/nats.go"
	"log"
	"middleware/example/internal/mail"
	"middleware/example/internal/models"
	"middleware/example/internal/webservice"
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

		// Expecting a JSON slice of changes
		var eventChanges []map[string]interface{}
		if err := json.Unmarshal(m.Data, &eventChanges); err != nil {
			log.Println("❌ Failed to decode event changes:", err)
			return
		}

		if len(eventChanges) <= 0 {
			return
		}

		var events []models.Event
		var event models.Event
		eventChangesMap := make(map[string]map[string]interface{})
		// Process each event
		for _, changeEntry := range eventChanges {

			eventUID, ok := changeEntry["uid"].(string)
			if !ok {
				log.Println("❌ Missing or invalid event UID in change entry")
				continue
			}

			// fmt.Println("Received Event from Timetable:", eventUID)
			response := webservice.HttpRequest("http://localhost:8090/events/search?uid="+eventUID, false)

			err := json.Unmarshal(response, &event)
			if err != nil {
				log.Println("❌ Failed to decode event JSON:", err)
				continue
			}

			// Remove "uid" from the changeEntry and store the rest as changes
			delete(changeEntry, "uid")
			eventChangesMap[eventUID] = changeEntry

			events = append(events, event)
		}

		models.DisplayEvents(events)

		// Handling sending alerts to users subscribed in a resourceID  with changes
		getAlerts(events, eventChangesMap)

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

func getAlerts(events []models.Event, eventChanges map[string]map[string]interface{}) {
	if len(events) <= 0 {
		//	fmt.Println("Everything is updated!")
		return
	}

	// Make the HTTP GET FROM CONFIG request
	var body []byte
	body = webservice.HttpRequest(models.CONFIG_ALERT_URL, false)

	// Parse the JSON data into a slice of alerts structs
	var alerts []models.Alert
	err := json.Unmarshal(body, &alerts)
	if err != nil {
		fmt.Println("Error parsing JSON:", err)
		return
	}

	// Print the parsed data
	fmt.Println("\nGetting all alerts from Config : ", models.CONFIG_ALERT_URL)
	// models.DisplayAlerts(alerts)

	handlingAlerts(alerts, events, eventChanges)
}

func handlingAlerts(alerts []models.Alert, events []models.Event, eventChanges map[string]map[string]interface{}) {
	for _, alert := range alerts {
		if alert.IsAll {
			mail.PreparingMail(alert.Email, events, eventChanges, true)
			continue
		}

		// checking if the resourceID of this alert exists in the resourceIDs of the events received
		subscribeEvents := models.CheckAlertResourceInEventResources(alert.ResourceID, events)
		if len(subscribeEvents) <= 0 {
			fmt.Printf("User with email %s has no Alert for his events\n", alert.Email)
			continue
		}
		mail.PreparingMail(alert.Email, subscribeEvents, eventChanges, false)
	}
}
