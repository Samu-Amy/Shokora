package payloads

import (
	"time"

	softerrors "github.com/Samu-Amy/Shokora/internal/errors/soft"
)

// The response sent to the frontend (with soft failure report)
type RegisterUserRes struct {
	User              UserRes            `json:"user"`
	VerificationId    *int64             `json:"verification_id,omitempty"`
	VerificationError softerrors.SoftErr `json:"verification_error,omitempty"`
	HasAuthError      bool               `json:"auth_error,omitempty"`
}

// Create a new RegisterUserRes with the user data and intializing the other fields
func NewRegisterUserRes(user UserRes) *RegisterUserRes {
	return &RegisterUserRes{
		User:              user,
		VerificationId:    nil,
		VerificationError: "",
		HasAuthError:      false,
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

// Create a new AuthTokensDto with the refresh token data
func NewAuthTokensDto(plainRefreshToken string, refreshTokenExpiresAt time.Time) *AuthTokensDto {
	return &AuthTokensDto{
		PlainRefreshToken:     plainRefreshToken,
		RefreshTokenExpiresAt: refreshTokenExpiresAt,
	}
}
