package authservice

import (
	"context"
	"time"

	"github.com/Samu-Amy/Shokora/internal/auth"
	interrors "github.com/Samu-Amy/Shokora/internal/errors/int"
	"github.com/Samu-Amy/Shokora/internal/mailer"
)

type VerificationEmailData struct {
	Name          string
	ActivationURL string
	MagicLinkExp  string
	OTPToken      string
	OTPExp        string
}

/*
This method is used to send an email for "email verification", "password reset" or "2 Factor Auth". Parameters:

context: context from the request

verificationType: enum in auth package (TokenEmailVerification, TokenPasswordReset, TokenTwoFactorAuth)

... user and token data

return: error (from SendEmail method in mailer Client)
*/
func (service *AuthService) sendVerificationEmail(ctx context.Context, verificationType auth.VerificationType, userName, email string, plainMagicLinkToken string, plainOTP string, magicLinkTokenExp, otpExp time.Duration) error {

	isSandbox := service.config.Mail.IsSandboxEnv

	// 2FA (only OTP)
	if verificationType == auth.TwoFactorAuth {
		vars := struct {
			Name     string
			OTPToken string
		}{
			Name:     userName,
			OTPToken: plainOTP,
		}

		return service.mailer.SendVerificationEmail(ctx, mailer.TwoFactorAuthTemplate, verificationType, userName, email, vars, isSandbox)
	}

	// Email Verification and Password Reset (magic link + OTP)
	if plainMagicLinkToken == "" {
		service.logger.Warnf("plainMagicLinkToken required for verification type: %v", verificationType)
		return interrors.IErrInvalidEmailVars
	}

	var templateFile mailer.TemplateFile
	activationURL := service.config.Mail.FrontEndURL

	switch verificationType {
	case auth.EmailVerification:
		templateFile = mailer.EmailVerificationTemplate
		activationURL += "/verify-email/"

	case auth.PasswordReset:
		templateFile = mailer.PasswordResetTemplate
		activationURL += "/reset-password/"
	}

	activationURL += plainMagicLinkToken

	vars := struct {
		Name          string
		ActivationURL string
		OTPToken      string
	}{
		Name:          userName,
		ActivationURL: activationURL,
		OTPToken:      plainOTP,
	}

	return service.mailer.SendVerificationEmail(ctx, templateFile, verificationType, userName, email, vars, isSandbox)
}
