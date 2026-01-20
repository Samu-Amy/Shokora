package mailer

import (
	"bytes"
	"context"
	"fmt"
	"html/template"
	"time"

	"github.com/google/uuid"
	"github.com/resend/resend-go/v2"
)

// TODO: aggiungi dominio (e setta SMTP port per SSL/TLS) a Resend (sulla dashboard)
// TODO: aggiungi unsubscribe per email marketing

type ResendMailer struct {
	apiKey    string
	fromEmail string
	client    *resend.Client
}

func NewResendMailer(apiKey, fromEmail string) *ResendMailer {
	return &ResendMailer{
		apiKey:    apiKey,
		fromEmail: fromEmail,
		client:    resend.NewClient(apiKey),
	}
}

func (mailer *ResendMailer) SendEmail(ctx context.Context, templateFile, name, email string, data any, isSandbox bool) error {

	// Template parsing
	tmpl, err := template.ParseFS(FS, "templates/"+templateFile) // TODO: ----- crea e usa template su Resend (forse no) o genera template con MJML (meglio) -----
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
	err = tmpl.ExecuteTemplate(body, "body", data)
	if err != nil {
		return err
	}

	params := &resend.SendEmailRequest{
		From:    mailer.fromEmail,
		To:      []string{email},
		Subject: subject.String(),
		Html:    body.String(),
	}

	verificationId := uuid.New().String()

	options := &resend.SendEmailOptions{
		IdempotencyKey: "verify-email/" + verificationId,
	}

	var retryErr error

	// Send email (with retries)
	for i := 0; i < MaxRetries; i++ {

		_, retryErr = mailer.client.Emails.SendWithOptions(ctx, params, options)
		if retryErr != nil {

			// Exponential backoff
			time.Sleep(time.Second * time.Duration(i+1))
			continue
		}

		return nil
	}

	// TODO: fai controllo errore per errori di invio duplicato della stessa email (ma comunque ricevuta)

	return fmt.Errorf("failed to send email after %d attempts, error: %v", MaxRetries, retryErr)
}
