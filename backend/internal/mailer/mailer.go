package mailer

import (
	"context"
	"embed"
)

// - Constants -

const (
	FromName                  string = "Shokora"
	MaxRetries                uint8  = 3
	EmailVerificationTemplate string = "email_verification.tmpl"
)

// - Templates -

// Embed template files into go binaries
//
//go:embed "templates"
var FS embed.FS

// - Interface -

type ClientI interface {
	SendEmail(ctx context.Context, templateFile, name, email string, data any, isSandbox bool) error
}
