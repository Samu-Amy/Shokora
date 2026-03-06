package auth

import "github.com/golang-jwt/jwt/v5"

type UserClaims struct {
	// TokenId int64 `json:"token_id"` // Può essere usato per revoca token
	UserId    int64 `json:"user_id,omitempty"`
	SessionId int64 `json:"session_id,omitempty"`
	jwt.RegisteredClaims
}
