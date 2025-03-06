package models

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
)

// SendEmail Sends an email using the external API
func SendEmail(to string, subject string, body string, token string) error {
	payload := MailRequest{
		To:      to,
		Subject: subject,
		Body:    body,
		Token:   token,
	}

	jsonData, err := json.Marshal(payload)
	if err != nil {
		return fmt.Errorf("error marshaling request: %w", err)
	}

	req, err := http.NewRequest("POST", MailAPI, bytes.NewBuffer(jsonData))
	if err != nil {
		return fmt.Errorf("error creating request: %w", err)
	}

	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return fmt.Errorf("error sending request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("failed to send email, status: %d", resp.StatusCode)
	}

	fmt.Println("Email sent successfully to", to)
	return nil
}
