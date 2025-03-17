package mq

import (
	"encoding/json"
	"fmt"
	"github.com/gofrs/uuid"
	"github.com/nats-io/nats.go"
	"io"
	"log"
	"middleware/example/internal/models"
	repository "middleware/example/internal/repositories/alerts"
	"net/http"
	"strings"
	"time"
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
			response := HttpRequest("http://localhost:8090/events/"+eventID, false)

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
			preparingMail(alert.Email, events, true)
			continue
		}

		// checking if the resourceID of this alert exists in the resourceIDs of the events received
		subscribeEvents := checkAlertResourceInEventResources(alert.ResourceID, events)
		if len(subscribeEvents) <= 0 {
			fmt.Printf("User with email %s has no Alert for his events\n", alert.Email)
			continue
		}
		preparingMail(alert.Email, subscribeEvents, false)
	}
}

func checkAlertResourceInEventResources(alertResourceID *uuid.UUID, events []models.Event) []models.Event {
	var subscribeEvents []models.Event
	for _, event := range events {
		for _, res := range event.ResourceIDs {
			if res != nil && res.String() == alertResourceID.String() {
				subscribeEvents = append(subscribeEvents, event)
				break
			}
		}
	}
	return subscribeEvents
}

func preparingMail(mail string, events []models.Event, all bool) {
	var eventsNames []string
	for _, event := range events {
		eventsNames = append(eventsNames, event.Name)
	}
	mailBody := strings.Join(eventsNames, ", ")

	if all {
		// send mail about all the events
		fmt.Printf("Sending mail concerning : %s \n", mailBody)
		sendMail(mail, mailBody)
		return
	}

	// send mail only about some subscribe events
	fmt.Printf("Sending mail concerning : %s \n", mailBody)
	sendMail(mail, mailBody)
}

func sendMail(mail string, content string) {
	// Token required for the API
	token := "PueiQkxDnrLjMHlFzfVVUCojDPTlZchQeRWecXTk"

	// Example event data
	event := struct {
		EventName string
		Start     string
		End       string
		Location  string
	}{
		EventName: content,
		Start:     time.Now().Format("2006-01-02 15:04"),
		End:       time.Now().Add(2 * time.Hour).Format("2006-01-02 15:04"),
		Location:  "ISIMA",
	}

	// Get email html shape from template
	emailContent, err := models.GetEmailContent("mail.html", event)
	if err != nil {
		log.Fatalf("Failed to generate email content: %s", err)
	}

	// Send email
	err = models.SendEmail("abdelali.ichou@etu.uca.fr", emailContent.Subject, emailContent.Body, token)
	if err != nil {
		log.Fatalf("Failed to send email: %s", err)
	}
}
