package api

import (
	"context"
	"fmt"
	"time"

	"github.com/Samu-Amy/Shokora/internal/auth"
	"github.com/Samu-Amy/Shokora/internal/mailer"
)

type VerificationEmailData struct {
	Name          string
	ActivationURL string
	MagicLinkExp  string
	OTPToken      string
	OTPExp        string
}

// TODO: fai creazione ed invio email in un metodo utils (con opzione per tipo di verifica -> generazione url e scelta template adatti)
// TODO: sistema le vars (anche OTP e scadenze (?)) - fai utils apposta per email verification, password reset e 2FA

/*
This method is used to send an email for "email verification", "password reset" or "2 Factor Auth". Parameters:

context: context from the request

verificationType: enum in auth package (TokenEmailVerification, TokenPasswordReset, TokenTwoFactorAuth)

... user and token data

return: error (from SendEmail method in mailer Client)
*/
func (app *App) SendVerificationEmail(ctx context.Context, verificationType auth.VerificationType, user_name, email string, plainMagicLinkToken *string, plainOTP string, magicLinkTokenExp, otpExp time.Duration) error {
	var templateFile mailer.TemplateFile
	var activationURL string

	vars := struct {
		Name          string
		ActivationURL string
		MagicLinkExp  string
		OTPToken      string
		OTPExp        string
	}{
		Name:         user_name,
		MagicLinkExp: FormatDurationToMinutes(magicLinkTokenExp),
		OTPToken:     plainOTP,
		OTPExp:       FormatDurationToMinutes(otpExp),
	}

	switch verificationType {
	case auth.EmailVerification:
		templateFile = mailer.EmailVerificationTemplate
		activationURL = "verify-email"

	case auth.PasswordReset:
		templateFile = mailer.PasswordResetTemplate
		activationURL = "reset-password"

	case auth.TwoFactorAuth:
		templateFile = mailer.TwoFactorAuthTemplate
	}

	if plainMagicLinkToken != nil {
		activationURL = fmt.Sprintf("%s/%s/%s", app.config.FrontEndURL, activationURL, plainMagicLinkToken)
		vars.ActivationURL = activationURL // TODO: continua
	}

	isProdEnv := app.config.Env == "prod"

	return app.mailer.SendEmail(ctx, templateFile, user_name, email, vars, !isProdEnv)
}
