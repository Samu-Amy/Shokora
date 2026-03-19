package payloads

import "github.com/google/uuid"

// Auth
type RegisterUserReq struct {
	UserDataReq
	EmailFieldReq
	PasswordFieldReq
}

type LoginUserReq struct {
	EmailFieldReq
	PasswordFieldReq
}

// Verification
type OTPVerificationReq struct {
	VerificationId uuid.UUID `json:"verification_id" validate:"gte=0"`
	OTP            string    `json:"otp" validate:"required,min=4,max=10"`
}

type RequestPasswordResetReq struct {
	EmailFieldReq
}

type ResetPasswordReq struct {
	PlainResetSessionToken string
	PasswordFieldReq
}
