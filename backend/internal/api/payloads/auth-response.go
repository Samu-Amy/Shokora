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
	AuthError         bool               `json:"auth_error,omitempty"`
}

// The data required to set cookies for auth
type AuthTokensDto struct {
	AccessToken           string
	PlainRefreshToken     string
	AccessTokenExpiresAt  time.Time
	RefreshTokenExpiresAt time.Time
}

// The data returned from RegisterUser (auth service)
type RegisterUserDto struct {
	RegisterUserRes
	AuthTokensDto
}
