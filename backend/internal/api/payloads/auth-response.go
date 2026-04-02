package payloads

import (
	"time"

	"github.com/google/uuid"
)

// ----- REGISTER -----

// The response sent to the frontend (with soft failure report)
type RegisterUserRes struct {
	User           UserRes    `json:"user"`
	VerificationId *uuid.UUID `json:"verification_id,omitempty"`
	IsEmailSent    bool       `json:"is_email_sent"` // Has the verification email been sent?
	HasAuthError   bool       `json:"has_auth_error"`
}

// Create a new RegisterUserRes with the user data and intializing the other fields
func NewRegisterUserRes(user *UserRes) *RegisterUserRes {
	return &RegisterUserRes{
		User:           *user,
		VerificationId: nil,
		IsEmailSent:    true,
		HasAuthError:   false,
	}
}

// ----- LOGIN -----

// The response sent to the frontend
type LoginUserRes struct {
	User           *UserRes   `json:"user,omitempty"`            // If present -> authenticated (already verified or to be verified, but no 2fa)
	VerificationId *uuid.UUID `json:"verification_id,omitempty"` // for 2fa or email verification (if nil: if user ok -> no verification required, if user nil -> verification error)
	IsEmailSent    bool       `json:"is_email_sent"`             // Has the verification email been sent?
}

// User != nil && VerificationID == nil -> user verified
// User != nil && VerificationID != nil -> verification required (user must be verified)
// User == nil && VerificationID != nil -> 2fa required

/*
The data required to set cookies for auth.
Should not be send this to frontend (is not serializable)
*/
type AuthTokensDto struct {
	AccessToken           string    `json:"-"`
	PlainRefreshToken     string    `json:"-"`
	AccessTokenExpiresAt  time.Time `json:"-"`
	RefreshTokenExpiresAt time.Time `json:"-"`
}

/*
Contains AuthTokensDto plus data for the middleware
*/
type AuthTokensCheckDto struct {
	// IsAccessTokenValid bool
	UserId        int64
	SessionId     int64
	AuthTokensDto AuthTokensDto
}
