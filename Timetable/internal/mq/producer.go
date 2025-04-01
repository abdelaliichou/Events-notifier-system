package mq

import (
	"encoding/json"
	"fmt"
	"github.com/nats-io/nats.go"
	"log"
)

// InitStream initializes the NATS JetStream connection
func InitStream() {

	// Connect to NATS server
	var err error
	nc, err := nats.Connect(nats.DefaultURL)
	if err != nil {
		log.Fatal("❌ Failed to connect to Alerts NATS:", err)
	}

	// Initialize JetStream context
	jsc, err = nc.JetStream()
	if err != nil {
		log.Fatal("❌ Failed to initialize Alerts JetStream:", err)
	}

	// Stream doesn't exist, create it
	_, err = jsc.AddStream(&nats.StreamConfig{
		Name:     "ALERTS",             // Name of the stream
		Subjects: []string{"ALERTS.>"}, // Allow any subject under ALERTS
	})
	if err != nil {
		log.Fatal("❌ Failed to create stream:", err)
	}

	fmt.Println("✅ Stream ALERTS created successfully")
}

// SendToAlerter will send created/modified events to the Alerter as a stream to MQ
func SendToAlerter(eventChanges []map[string]interface{}) {

	if len(eventChanges) == 0 {
		log.Println("No changes to alert.")
		return
	}

	err := publishEventsAsStream(eventChanges)
	if err != nil {
		log.Println("Error sending events to MQ:", err)
		return
	}

	fmt.Println("✅ All events sent successfully to MQ")
}

// publishEventsAsStream Publish a list of events as a stream
func publishEventsAsStream(eventChanges []map[string]interface{}) error {

	// Convert the whole list to JSON
	data, err := json.Marshal(eventChanges)
	if err != nil {
		log.Println("Failed to encode event list:", err)
		return err
	}

	// Publish the entire list as one message
	pubAckFuture, err := jsc.PublishAsync("ALERTS.stream", data)
	if err != nil {
		log.Println("Failed to publish Alerts event list from the start:", err)
		return err
	}

	// Wait for acknowledgment
	select {
	case <-pubAckFuture.Ok():
		fmt.Println("✅ ALERTS Event list published successfully!")
		return nil
	case err := <-pubAckFuture.Err():
		log.Println("❌ Failed to publish ALERTS event list:", err)
		return err
	}
}
