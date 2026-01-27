package auth

import (
	"time"

	"github.com/golang-jwt/jwt/v5"
)

type JWTService interface {
	GenerateJWTToken(claims jwt.Claims) (string, error)
	ValidateJWTToken(token string) (*jwt.Token, error)
}

type TokenService interface {
	CreateVerificationTokens(verificationType VerificationType) (*VerificationTokens, error)

	GenerateMagicLinkToken() (string, error)
	GenerateOTP() (string, error)

	HashMagicLinkToken(plainMagicLinkToken string) []byte
	HashOTP(plainOTP string, verificationType VerificationType) []byte

	VerifyMagicLinkToken(plainToken string, hashedToken []byte) bool
	VerifyOTP(plainOTP string, hashedToken []byte, verificationType VerificationType) bool

	getExpiration(verificationType VerificationType) time.Duration
	getVerificationTypeString(verificationType VerificationType) string
}
