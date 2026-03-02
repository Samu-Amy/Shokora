package payloads

import (
	"time"
)

// The response sent to the frontend (with soft failure report)
type RegisterUserRes struct {
	User           UserRes `json:"user"`
	VerificationId *int64  `json:"verification_id,omitempty"`
	IsEmailSent    bool    `json:"is_email_sent"`
	HasAuthError   bool    `json:"has_auth_error"`
}

// Create a new RegisterUserRes with the user data and intializing the other fields
func NewRegisterUserRes(user UserRes) *RegisterUserRes {
	return &RegisterUserRes{
		User:           user,
		VerificationId: nil,
		IsEmailSent:    true,
		HasAuthError:   false,
	}
}

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
