package mailer

import "embed"

const (
	FromName            = "Shokora"
	MaxRetries          = 3
	UserWelcomeTemplate = "user_invitation.tmpl"
)

// Embed template files into go binaries
//go:embed "templates"
var FS embed.FS

type Client interface {
	Send(templateFile, name, email string, data any, isSandbox bool) error
}
