package auth

import (
	"github.com/golang-jwt/jwt/v5"
)

type JWTService interface {
	GenerateJWTToken(claims jwt.Claims) (string, error)
	ValidateJWTToken(token string) (*jwt.Token, error)
}

type TokenService interface {
	CreateVerificationTokens(tokenType TokenType) (*VerificationTokens, error)
	GenerateVerificationToken() (string, error)
	GenerateOTP() (string, error)
	HashToken(plainToken string) []byte
	VerifyToken(plainToken string, hashedToken []byte) bool
}
