package mailer

import (
	"context"
	"embed"

	"github.com/Samu-Amy/Shokora/internal/auth"
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
	TwoFactorAuthTemplate     TemplateFile = "two_factor_auth.tmpl" // TODO: ricorda di non mettere magic link in 2FA
)

// - Templates -

// Embed template files into go binaries
//
//go:embed "templates"
var FS embed.FS

// - Interface -

type IClient interface {
	SendVerificationEmail(ctx context.Context, templateFile TemplateFile, verificationType auth.VerificationType, name, email string, data any, isSandbox bool) error
}
