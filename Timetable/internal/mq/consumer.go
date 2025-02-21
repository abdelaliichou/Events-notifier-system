package mq

import (
	"encoding/json"
	"fmt"
	"github.com/nats-io/nats.go"
	"log"
	"middleware/example/internal/models"
)

var jsc nats.JetStreamContext

// StartStreamConsumer Starts Consumer for the Event Stream coming from the scheduler
func StartStreamConsumer() {

	// Connect to NATS server
	var err error
	nc, err := nats.Connect(nats.DefaultURL)
	if err != nil {
		log.Fatal("❌ Failed to connect to NATS:", err)
	}

	// Initialize JetStream context
	jsc, err = nc.JetStream()
	if err != nil {
		log.Fatal("❌ Failed to initialize JetStream:", err)
	}

	fmt.Println("✅ Connected to NATS and JetStream initialized")

	// Subscribe to the EVENTS.stream using JetStream
	_, err = jsc.Subscribe("EVENTS.stream", func(m *nats.Msg) {

		fmt.Printf("Received a message: %s\n", m.Subject)

		// Expecting and jsonify the list of events
		var events []models.Event
		if err := json.Unmarshal(m.Data, &events); err != nil {
			log.Println("❌ Failed to decode event list:", err)
			return
		}

		// Process each event
		for _, event := range events {
			fmt.Println("Received Event from Timetable:", event)

			// Here you would call your storage function (currently commented out)
			// err = storage.SaveEvent(event)
			// if err != nil {
			// 	log.Println("❌ Failed to save event:", err)
			// } else {
			// 	fmt.Println("✅ Event Received & Stored:", event)
			// }
		}

		// Acknowledge the message after processing
		m.Ack()

		fmt.Println("✅ Event processed and acknowledged:", m.Subject)

	}, nats.Durable("event-consumer"), nats.ManualAck())

	// Error handling for subscription
	if err != nil {
		log.Fatal("❌ Failed to subscribe to EVENTS.stream:", err)
	}

	log.Println("✅ Event Stream Consumer is running...")

	// Prevent program from exiting, keeping the consumer alive forever
	select {}
}
