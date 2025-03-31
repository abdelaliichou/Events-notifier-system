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

		if len(events) <= 0 {
			return
		}

		// Process each event
		var CreatedOrModifiedEvents []map[string]interface{}
		for i, event := range events {

			models.DisplayEvents(event, i, true)

			// Deciding if we add event to db or not
			eventChanges, err := UpsertEvent(event)
			if err != nil {
				log.Println("❌ Failed to save event:", err)
			} else {
				if eventChanges != nil {
					CreatedOrModifiedEvents = append(CreatedOrModifiedEvents, eventChanges)
				}
			}
		}

		// Acknowledge the message after processing
		m.Ack()

		// save the updated event so we send it to the Alerter using nats
		sendToAlerter(CreatedOrModifiedEvents)

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
func UpsertEvent(event models.Event) (map[string]interface{}, error) {
	// Check if event exists using UID
	existingEvent, err := repository.GetEventByUID(event.UID)

	if err != nil {
		if customErr, ok := err.(*models.CustomError); ok && customErr.Code == http.StatusNotFound {
			logrus.Infof("No existing event found for UID: %s. Proceeding to create it.", event.UID)
			existingEvent = nil // Ensure nil so we proceed to create it
		} else {
			logrus.Errorf("Error fetching event by UID: %s", err.Error())
			fmt.Println("-----")
			return nil, err
		}
	}

	// If event does not exist, create it
	if existingEvent == nil {

		_, err = Events.CreatEvent(event)
		if err != nil {
			logrus.Errorf("Error creating event: %s", err.Error())
			fmt.Println("-----")
			return nil, err
		}

		// Insert resource IDs
		err = repository.InsertEventResources(event.UID, event.ResourceIDs)
		if err != nil {
			logrus.Errorf("Error inserting resource IDs: %s", err.Error())
			return nil, err
		}

		fmt.Printf("Event with UID %s created successfully:\n", event.UID)
		models.DisplayEvents(event, 0, false)

		return map[string]interface{}{
			"uid":    event.UID,
			"status": "created",
		}, nil
	}

	// Check what has changed
	changes := models.GetEventChanges(existingEvent, &event)
	if len(changes) == 0 {
		fmt.Printf("Event with UID %s exists but has no changes, skipping update.\n", event.UID)
		fmt.Println("-----")
		return nil, nil
	}

	// Update the existing event
	logrus.Infof("Event with UID %s has modifications, updating it.", event.UID)

	err = Events.UpdateEvent(event)
	if err != nil {
		logrus.Errorf("Error updating event: %s", err.Error())
		return nil, err
	}

	// Update resource IDs
	err = repository.UpdateEventResources(event.UID, event.ResourceIDs)
	if err != nil {
		logrus.Errorf("Error updating resource IDs: %s", err.Error())
		return nil, err
	}

	//fmt.Printf("Event updated successfully: %+v\n", event)
	fmt.Printf("Event updated successfully:\n")
	models.DisplayEvents(event, 0, false)

	return map[string]interface{}{
		"uid":     event.UID,
		"changes": changes,
	}, nil
}

// sendModificationsToAlerter will send created/modified events to the Alerter as a stream to MQ
func sendToAlerter(eventChanges []map[string]interface{}) {

	if len(eventChanges) == 0 {
		log.Println("No changes to alert.")
		return
	}

	err := PublishEventsAsStream(eventChanges)
	if err != nil {
		log.Println("Error sending events to MQ:", err)
		return
	}

	fmt.Println("✅ All events sent successfully to MQ")
}
