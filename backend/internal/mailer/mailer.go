package mailer

import (
	"embed"
)

// - Constants -

const (
	FromName            = "Shokora"
	MaxRetries          = 3
	UserWelcomeTemplate = "user_invitation.tmpl"
)

// - Templates -

// Embed template files into go binaries
//
//go:embed "templates"
var FS embed.FS

// - Errors -

// TODO: togli (?)
// var (
// 	SendEmailErr = errors.New("impossibile inviare l'email")
// )

// - Interface -

type Client interface {
	Send(templateFile, name, email string, data any, isSandbox bool) error
}
