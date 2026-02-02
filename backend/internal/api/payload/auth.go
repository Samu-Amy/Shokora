package payload

// - Payloads -
type VerificationPayload struct {
	OTP string `json:"otp" validate:"required,min=4,max=10"`
}
