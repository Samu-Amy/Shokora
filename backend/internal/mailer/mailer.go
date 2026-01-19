package mailer

import (
	"context"
	"embed"
)

// - Constants -

const (
	FromName                  = "Shokora"
	MaxRetries                = 3
	EmailVerificationTemplate = "email_verification.tmpl"
)

// - Templates -

// Embed template files into go binaries
//
//go:embed "templates"
var FS embed.FS

// - Interface -

type Client interface {
	SendEmail(ctx context.Context, templateFile, name, email string, data any, isSandbox bool) error
}
