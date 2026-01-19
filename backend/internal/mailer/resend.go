package mailer

import "github.com/resend/resend-go/v2"

// TODO: aggiungi dominio a Resend (sulla dashboard)

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

func (mailer *ResendMailer) Send(templateFile, name, email string, data any, isSandbox bool) error {
	params := &resend.SendEmailRequest{
		From:    mailer.fromEmail,
		To:      []string{email},
		Subject: "Test",
		Html:    "<h2>Siamo felici di averti qui</h2><p>Per verificare il tuo indirizzo email clicca sul link qua sotto.</p>",
	}

	// TODO: Retry?
	_, err := mailer.client.Emails.Send(params)
	if err != nil {
		return err
	}

	return nil
}
