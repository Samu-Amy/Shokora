package payloads

import (
	"time"
)

// ----- REGISTER -----

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

// ----- LOGIN -----

// The response sent to the frontend (with soft failure report)
type LoginUserRes struct {
	User UserRes `json:"user"`
	// VerificationId *int64  `json:"verification_id,omitempty"` // TODO: forse può servire per 2fa
}

// Create a new RegisterUserRes with the user data and intializing the other fields
func NewLoginUserRes(user UserRes) *LoginUserRes {
	return &LoginUserRes{
		User: user,
		// VerificationId: nil,
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

/*
Contains AuthTokensDto plus data for the middleware
*/
type AuthTokensCheckDto struct {
	IsAccessTokenValid bool
	UserId             int64
	SessionId          int64
	TokensDto          AuthTokensDto
}
