package mailer

import (
	"bytes"
	"fmt"
	"html/template"
	"log"
	"time"

	"github.com/sendgrid/sendgrid-go"
	"github.com/sendgrid/sendgrid-go/helpers/mail"
)

type SendGridMailer struct {
	fromEmail string
	apiKey    string
	client    *sendgrid.Client
}

func NewSendGrid(apiKey, fromEmail string) *SendGridMailer {
	client := sendgrid.NewSendClient(apiKey)
	return &SendGridMailer{
		fromEmail: fromEmail,
		apiKey:    apiKey,
		client:    client,
	}
}

func (mailer *SendGridMailer) Send(templateFile, name, email string, data any, isSandbox bool) error {
	from := mail.NewEmail(FromName, mailer.fromEmail)
	to := mail.NewEmail(name, email)

	// Template parsing
	tmpl, err := template.ParseFS(FS, "templates/"+templateFile)
	if err != nil {
		return err
	}

	// Get subject from template
	subject := new(bytes.Buffer)
	err = tmpl.ExecuteTemplate(subject, "subject", data)
	if err != nil {
		return err
	}

	// Get body from template
	body := new(bytes.Buffer)
	err = tmpl.ExecuteTemplate(subject, "subject", data)
	if err != nil {
		return err
	}

	// Create mail
	message := mail.NewSingleEmail(from, subject.String(), to, "", body.String())

	// Sandbox
	message.SetMailSettings(&mail.MailSettings{
		SandboxMode: &mail.Setting{
			Enable: &isSandbox,
		},
	})

	// Send email (with retries)
	for i := 0; i < MaxRetries; i++ {
		response, err := mailer.client.Send(message)
		if err != nil {
			log.Printf("Failer to send email to %v, attempt %d od %d", email, i+1, MaxRetries)
			log.Printf("Error: %v", err.Error())

			// Exponential backoff
			time.Sleep(time.Second * time.Duration(i+1))
			continue
		}

		log.Printf("Email sent with status code %v", response.StatusCode)
		return nil
	}

	return fmt.Errorf("failed to send email after %d attempts", MaxRetries)
}
