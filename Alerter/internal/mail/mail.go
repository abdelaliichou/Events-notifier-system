package mail

import "middleware/example/internal/models"

import (
	"fmt"
	"log"
	"strings"
)

func PreparingMail(mail string, events []models.Event, eventChanges map[string]map[string]interface{}, all bool) {
	// Loop through each event and format the changes
	for _, event := range events {

		var mailBody strings.Builder

		mailBody.WriteString(fmt.Sprintf("Bonjour.\n\nÉvénement %s a été modifié.\n\n", event.Name))

		eventData, exists := eventChanges[event.UID]
		if !exists {
			continue
		}

		changes, hasChanges := eventData["changes"].(map[string]interface{})
		if !hasChanges {
			continue
		}

		mailBody.WriteString("Les changements apportés : \n\n")
		for field, change := range changes {
			changeMap, ok := change.(map[string]interface{})
			if !ok {
				continue
			}

			oldValue, oldExists := changeMap["old"].(string)
			newValue, newExists := changeMap["new"].(string)

			if oldExists && newExists {
				mailBody.WriteString(fmt.Sprintf("Field \"%s\": \n    Avant: %s\n    Après: %s\n\n", field, oldValue, newValue))
			}
		}

		// Add event details
		mailBody.WriteString("\nDétails : \n")
		mailBody.WriteString(fmt.Sprintf("Nom : %s\n", event.Name))
		mailBody.WriteString(fmt.Sprintf("Description : %s\n", event.Description))
		mailBody.WriteString(fmt.Sprintf("Début : %s\n", event.Start.Format("2006-01-02 15:04:05")))
		mailBody.WriteString(fmt.Sprintf("Fin : %s\n", event.End.Format("2006-01-02 15:04:05")))
		mailBody.WriteString(fmt.Sprintf("Lieu : %s\n", event.Location))
		mailBody.WriteString(fmt.Sprintf("Dernière mise à jour : %s\n\n", event.LastUpdate.Format("2006-01-02 15:04:05")))

		mailBody.WriteString("Cordialement,\nAbdelali ichou du L'équipe P&G Innovations\n")
		//mailBody.WriteString("------------------------------------------------------\n")

		// Send mail
		// fmt.Printf("📧 Sending mail to %s:\n%s\n", mail, mailBody.String())
		sendMail(mail, mailBody.String())
	}
}

func sendMail(mail string, content string) {
	// Token required for the API
	token := "PueiQkxDnrLjMHlFzfVVUCojDPTlZchQeRWecXTk"

	// Example event data
	event := struct {
		EventContent string
	}{
		EventContent: content,
	}

	// Get email html shape from template
	emailContent, err := models.GetEmailContent("mail.html", event)
	if err != nil {
		log.Fatalf("Failed to generate email content: %s", err)
	}

	// Send email
	// here im working with abdelali.ichou@etu.uca.fr instead of the mail argument just for testing purposes
	err = models.SendEmail("abdelali.ichou@etu.uca.fr", emailContent.Subject, emailContent.Body, token)
	if err != nil {
		log.Fatalf("Failed to send email: %s", err)
	}
}
