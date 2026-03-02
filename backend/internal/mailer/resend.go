package mailer

import (
	"bytes"
	"context"
	"fmt"
	"html/template"
	"time"

	"github.com/Samu-Amy/Shokora/internal/auth"
	interrors "github.com/Samu-Amy/Shokora/internal/errors/int"
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

func (mailer *ResendMailer) SendVerificationEmail(ctx context.Context, templateFile TemplateFile, verificationType auth.VerificationType, name, email string, data any, isSandbox bool) error {

	if isSandbox {
		fmt.Printf("\n\nEmail Sandbox: %v\n\n", data)
		return nil
	}

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

	// Generate random verification id
	verificationId := uuid.New().String() // TODO: fare qualcosa di più efficiente e più "sicuro" (per evitare duplicati e non inviare mail - anche se re-inviando la mail se ne genera uno nuovo)?
	var idepotencyKeyPrefix string

	switch verificationType {
	case auth.EmailVerification:
		idepotencyKeyPrefix = "verify-email/"
	case auth.PasswordReset:
		idepotencyKeyPrefix = "password-reset/"
	case auth.TwoFactorAuth:
		idepotencyKeyPrefix = "two-factor-auth/"
	}

	options := &resend.SendEmailOptions{
		IdempotencyKey: idepotencyKeyPrefix + verificationId,
	}

	var retryErr error

	// Send email (with retries)
	for i := range MaxRetries {

		// TODO: va bene (aggiungere timeout - magari nell'handler?), utile per retry e altre operazioni che prendono tempo?
		select {
		case <-ctx.Done():
			return ctx.Err()
		default:
		}

		_, retryErr = mailer.client.Emails.SendWithOptions(ctx, params, options)
		if retryErr != nil {
			// TODO: controlla errore (se per mail duplicata o email sbagliata (4xx)) -> no retry

			// Linear backoff
			time.Sleep(time.Second * time.Duration(i+1))
			continue
		}

		return nil
	}

	// TODO: fai controllo errore per errori di invio duplicato della stessa email (ma comunque ricevuta)

	return interrors.IErrMaxRetriesExceeded
}
