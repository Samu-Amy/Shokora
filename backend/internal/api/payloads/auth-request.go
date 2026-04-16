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
	State string `json:"state" validate:"required,valid-base64-rawurl-32"` // base64
	Code  string `json:"code" validate:"required,max=512,safe-chars"`      // 512 is just to avoid strings too long
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
	PlainResetSessionToken string `json:"plain_reset_session_token" validate:"required,valid-base64-rawurl-32"` // base64
	PasswordFieldReq
}
