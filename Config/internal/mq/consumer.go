package mq

import (
	"encoding/json"
	"fmt"
	"github.com/nats-io/nats.go"
	"log"
	"middleware/example/internal/mail"
	"middleware/example/internal/models"
	repository "middleware/example/internal/repositories/alerts"
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

		// Expecting and parsing the list of event IDs
		var eventIDs []string
		if err := json.Unmarshal(m.Data, &eventIDs); err != nil {
			log.Println("❌ Failed to decode event list:", err)
			return
		}

		if len(eventIDs) <= 0 {
			return
		}

		var events []models.Event
		var event models.Event
		// Process each event
		for _, eventID := range eventIDs {

			fmt.Println("Received Event from Timetable:", eventID)
			response := webservice.HttpRequest("http://localhost:8090/events/search?uid="+eventID, false)

			err := json.Unmarshal(response, &event)
			if err != nil {
				return
			}

			events = append(events, event)
		}

		models.DisplayEvents(events)

		// Handling sending alerts to users subscribed in a resourceID
		preparingAlerts(events)

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

func preparingAlerts(events []models.Event) {
	if len(events) <= 0 {
		//	fmt.Println("Everything is updated!")
		return
	}

	alerts, err := repository.GetAllAlerts()
	if err != nil {
		return
	}

	models.DisplayAlerts(alerts)
	handlingAlerts(alerts, events)
}

func handlingAlerts(alerts []*models.Alert, events []models.Event) {
	for _, alert := range alerts {
		if alert.IsAll {
			mail.PreparingMail(alert.Email, events, true)
			continue
		}

		// checking if the resourceID of this alert exists in the resourceIDs of the events received
		subscribeEvents := models.CheckAlertResourceInEventResources(alert.ResourceID, events)
		if len(subscribeEvents) <= 0 {
			fmt.Printf("User with email %s has no Alert for his events\n", alert.Email)
			continue
		}
		mail.PreparingMail(alert.Email, subscribeEvents, false)
	}
}
