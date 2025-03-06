package models

import (
	"bytes"
	"embed"
	"fmt"
	"github.com/adrg/frontmatter"
	"strings"
	"text/template"
)

//go:embed mail.html
var embeddedTemplates embed.FS

// GetEmailContent Extracts subject and formatted email content
func GetEmailContent(templatePath string, data interface{}) (EmailContent, error) {
	var content EmailContent

	// Ensure the correct template file path
	temp, err := template.ParseFS(embeddedTemplates, templatePath)
	if err != nil {
		return content, fmt.Errorf("error parsing template: %w", err)
	}

	var tpl bytes.Buffer
	if err = temp.Execute(&tpl, data); err != nil {
		return content, fmt.Errorf("error executing template: %w", err)
	}

	// Extract subject using frontmatter
	var matter struct {
		Subject string `yaml:"subject"`
	}
	mailContent, err := frontmatter.Parse(strings.NewReader(tpl.String()), &matter)
	if err != nil {
		return content, fmt.Errorf("error parsing frontmatter: %w", err)
	}

	content.Subject = matter.Subject
	content.Body = string(mailContent)
	return content, nil
}
