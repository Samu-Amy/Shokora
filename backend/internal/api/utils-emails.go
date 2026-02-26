package api

import (
	"context"
	"time"

	"github.com/Samu-Amy/Shokora/internal/auth"
	"github.com/Samu-Amy/Shokora/internal/errorcodes"
	"github.com/Samu-Amy/Shokora/internal/mailer"
)

type VerificationEmailData struct {
	Name          string
	ActivationURL string
	MagicLinkExp  string
	OTPToken      string
	OTPExp        string
}

// TODO: fai creazione ed invio email in un service (con opzione per tipo di verifica -> generazione url e scelta template adatti)
// TODO: sistema le vars (anche OTP e scadenze (?)) - fai utils apposta per email verification, password reset e 2FA

/*
This method is used to send an email for "email verification", "password reset" or "2 Factor Auth". Parameters:

context: context from the request

verificationType: enum in auth package (TokenEmailVerification, TokenPasswordReset, TokenTwoFactorAuth)

... user and token data

return: error (from SendEmail method in mailer Client)
*/
func (app *App) SendVerificationEmail(ctx context.Context, verificationType auth.VerificationType, user_name, email string, plainMagicLinkToken *string, plainOTP string, magicLinkTokenExp, otpExp time.Duration) error {

	isSandbox := app.config.Mail.SandboxEnv

	// 2FA (only OTP)
	if verificationType == auth.TwoFactorAuth {
		vars := struct {
			Name     string
			OTPToken string
		}{
			Name:     user_name,
			OTPToken: plainOTP,
		}

		return app.mailer.SendEmail(ctx, mailer.TwoFactorAuthTemplate, user_name, email, vars, isSandbox)
	}

	// Email Verification and Password Reset (magic link + OTP)
	if plainMagicLinkToken == nil {
		app.logger.Warnf("plainMagicLinkToken required for verification type: %v", verificationType)
		return errorcodes.InternalErrInvalidEmailVars
	}

	var templateFile mailer.TemplateFile
	activationURL := app.config.FrontEndURL

	switch verificationType {
	case auth.EmailVerification:
		templateFile = mailer.EmailVerificationTemplate
		activationURL += "/verify-email/"

	case auth.PasswordReset:
		templateFile = mailer.PasswordResetTemplate
		activationURL += "/reset-password/"
	}

	activationURL += *plainMagicLinkToken

	vars := struct {
		Name          string
		ActivationURL string
		OTPToken      string
	}{
		Name:          user_name,
		ActivationURL: activationURL,
		OTPToken:      plainOTP,
	}

	return app.mailer.SendEmail(ctx, templateFile, user_name, email, vars, isSandbox)
}
