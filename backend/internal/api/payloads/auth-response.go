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
DO NOT send this to frontend
*/
type AuthTokensDto struct {
	AccessToken           string
	PlainRefreshToken     string
	AccessTokenExpiresAt  time.Time
	RefreshTokenExpiresAt time.Time
}

/*
The data returned from RegisterUser (auth service).
SEND ONLY UserRes to frontend
*/
// type RegisterUserDto struct {
// 	UserRes    RegisterUserRes
// 	AuthTokens AuthTokensDto
// }
