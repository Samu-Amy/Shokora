package payloads

import "github.com/google/uuid"

// Auth
type RegisterUserReq struct {
	UserDataReq
	EmailFieldReq
	DoublePasswordFieldReq
}

type LoginUserReq struct {
	EmailFieldReq
	PasswordFieldReq
}

type GoogleOAuthCallbackReq struct {
	State string `json:"state" validate:"required,min=43,max=43,safe-chars"`
	Code  string `json:"code" validate:"required,safe-chars"`
}

// Verification
type OTPVerificationReq struct {
	VerificationId uuid.UUID `json:"verification_id" validate:"required"`
	OTP            string    `json:"otp" validate:"required,valid-otp"`
}

type SendVerificationReq struct {
	EmailFieldReq
}

type ResetPasswordReq struct {
	PlainResetSessionToken string `json:"plain_reset_session_token" validate:"required,safe-chars"`
	PasswordFieldReq
}
