package mq

import (
	"encoding/json"
	"fmt"
	"github.com/nats-io/nats.go"
	"github.com/sirupsen/logrus"
	"log"
	"middleware/example/internal/models"
	repository "middleware/example/internal/repositories/events"
	Events "middleware/example/internal/services/events"
	"net/http"
)

var jsc nats.JetStreamContext

// StartStreamConsumer Starts Consumer for the Event Stream coming from the scheduler
func StartStreamConsumer() {

	// Connect to NATS server
	var err error
	nc, err := nats.Connect(nats.DefaultURL)
	if err != nil {
		log.Fatal("❌ Failed to connect to Scheduler NATS:", err)
	}

	// Initialize JetStream context
	jsc, err = nc.JetStream()
	if err != nil {
		log.Fatal("❌ Failed to initialize Scheduler JetStream:", err)
	}

	fmt.Println("✅ Connected to Scheduler NATS and JetStream initialized")

	// Subscribe to the EVENTS.stream using JetStream
	_, err = jsc.Subscribe("EVENTS.stream", func(m *nats.Msg) {

		// Expecting and jsonify the list of events
		var events []models.Event
		if err := json.Unmarshal(m.Data, &events); err != nil {
			log.Println("❌ Failed to decode event list:", err)
			return
		}

		// Process each event
		var CreatedModifiedEvents []string
		for _, event := range events {

			fmt.Println("Received Event from Scheduler:", event)

			// Saving the event in the db
			eventUID, err := UpsertEvent(event)
			if err != nil {
				log.Println("❌ Failed to save event:", err)
			} else {
				if eventUID != "" {
					CreatedModifiedEvents = append(CreatedModifiedEvents, eventUID)
				}
			}
		}

		// Acknowledge the message after processing
		m.Ack()

		// save the updated event so we send it to the Alerter using nats
		sendToAlerter(CreatedModifiedEvents)

	}, nats.Durable("event-consumer"), nats.ManualAck())

	if err != nil {
		log.Fatal("❌ Failed to subscribe to EVENTS.stream:", err)
	}

	log.Println("✅ Event Stream Consumer is running...")

	// Prevent program from exiting, keeping the consumer alive forever
	select {}
}

// UpsertEvent is for checking if there's any modification and to see if we add it to db or not
// and after this we need to add nats to send to alerter about anything changed in this db
func UpsertEvent(event models.Event) (string, error) {
	// Check if event exists using UID
	existingEvent, err := repository.GetEventByUID(event)

	if err != nil {
		if customErr, ok := err.(*models.CustomError); ok && customErr.Code == http.StatusNotFound {
			logrus.Infof("No existing event found for UID: %s. Proceeding to create it.", event.UID)
			existingEvent = nil // Ensure nil so we proceed to create it
		} else {
			logrus.Errorf("Error fetching event by UID: %s", err.Error())
			return "", err
		}
	}

	// If event does not exist, create it
	if existingEvent == nil {

		_, err = Events.CreatEvent(event)
		if err != nil {
			logrus.Errorf("Error creating event: %s", err.Error())
			return "", err
		}

		fmt.Printf("New event created successfully: %+v\n", event)

		return event.UID, nil
	}

	// Event exists, check for modifications
	if !models.IsEventModified(existingEvent, &event) {
		fmt.Printf("Event with UID %s exists but has no changes, skipping update.\n", event.UID)
		return "", nil
	}

	// Update the existing event
	fmt.Printf("Event with UID %s has modifications, updating it.\n", event.UID)

	err = Events.UpdateEvent(event)
	if err != nil {
		logrus.Errorf("Error updating event: %s", err.Error())
		return "", err
	}

	fmt.Printf("Event updated successfully: %+v\n", event)

	return event.UID, nil
}

// sendModificationsToAlerter will send created/modified events to the Alerter as a stream to MQ
func sendToAlerter(eventUIDs []string) {

	err := PublishEventsAsStream(eventUIDs)
	if err != nil {
		log.Println("Error sending events to MQ:", err)
		return
	}

	fmt.Println("✅ All events sent successfully to MQ")
}
