package mq

import (
	"Scheduler/models"
	"encoding/json"
	"fmt"
	"github.com/nats-io/nats.go"
	"log"
)

var jsc nats.JetStreamContext

// InitStream initializes the NATS JetStream connection
func InitStream() {

	// Connect to NATS server
	var err error
	//nc, err := nats.Connect(nats.DefaultURL)
	nc, err := nats.Connect("nats://nats-server:4222")
	if err != nil {
		log.Fatal("❌ Failed to connect to Scheduler NATS:", err)
	}

	// Initialize JetStream context
	jsc, err = nc.JetStream()
	if err != nil {
		log.Fatal("❌ Failed to initialize Scheduler JetStream:", err)
	}

	// Stream doesn't exist, create it
	_, err = jsc.AddStream(&nats.StreamConfig{
		Name:     "EVENTS",             // Name of the stream
		Subjects: []string{"EVENTS.>"}, // Allow any subject under EVENTS
	})
	if err != nil {
		log.Fatal("❌ Failed to create stream:", err)
	}

	log.Println("✅ Stream EVENTS created successfully")
}

// SendEventsToMQ will send structured events to our producer as a stream to MQ
func SendEventsToMQ(structuredEvents []models.Event) {

	err := publishEventsAsStream(structuredEvents)
	if err != nil {
		log.Println("Error sending events to MQ:", err)
		return
	}

	fmt.Println("All events sent successfully to MQ")
}

// publishEventsAsStream Publish a list of events as a stream
func publishEventsAsStream(events []models.Event) error {

	// Convert the whole list to JSON
	data, err := json.Marshal(events)
	if err != nil {
		log.Println("Failed to encode event list:", err)
		return err
	}

	// Publish the entire list as one message
	pubAckFuture, err := jsc.PublishAsync("EVENTS.stream", data)
	if err != nil {
		log.Println("Failed to publish event list from the start:", err)
		return err
	}

	// Wait for acknowledgment
	select {
	case <-pubAckFuture.Ok():
		fmt.Println("✅ Event list published successfully!")
		return nil
	case err := <-pubAckFuture.Err():
		log.Println("❌ Failed to publish event list:", err)
		return err
	}
}
