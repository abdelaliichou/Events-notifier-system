package mail

import "middleware/example/internal/models"

import (
	"fmt"
	"log"
	"strings"
	"time"
)

func PreparingMail(mail string, events []models.Event, all bool) {
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
