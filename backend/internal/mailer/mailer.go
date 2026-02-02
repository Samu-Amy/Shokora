package mailer

import (
	"context"
	"embed"
)

// - Constants -

// General
const (
	FromName   string = "Shokora"
	MaxRetries uint8  = 3
)

// Template files
type TemplateFile = string

const (
	EmailVerificationTemplate TemplateFile = "email_verification.tmpl"
	PasswordResetTemplate     TemplateFile = "password_reset.tmpl"
	TwoFactorAuthTemplate     TemplateFile = "two_factor_auth.tmpl"
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
