package auth

import (
	"bytes"
	"crypto/rand"
	"crypto/sha256"
	"encoding/base64"
	"fmt"
	"math/big"
	"time"
)

// Authenticator
type TokenAuthenticator struct {
	MagicLink MagicLinkConfig
	OTP       OTPConfig
}

type MagicLinkConfig struct {
	ByteSize int
	Exp      time.Duration
}

type OTPConfig struct {
	Length      int8
	MaxAttempts int8
	LongExp     time.Duration // For Email Verification
	BaseExp     time.Duration // For Password Reset and 2FA
	// CriticalExp time.Duration // For critical operations (es. 30s)
}

// Tokens
type VerificationTokens struct { // TODO: fai methodi per seplificarne la creazione (?)
	TokenType   TokenType
	PlainToken  string
	HashedToken []byte
	PlainOTP    string
	HashedOTP   []byte
	OTPExp      time.Duration
}

// Verification Token and OTP
type TokenType uint8

const (
	TokenEmailVerification TokenType = 0
	TokenPasswordReset     TokenType = 1
	TokenTwoFactorAuth     TokenType = 2
)

// - Constructor -

func NewTokenAuthenticator(MagicLink MagicLinkConfig, OTP OTPConfig) *TokenAuthenticator {
	return &TokenAuthenticator{MagicLink, OTP}
}

// - Methods -

func (tokenAuthenticator *TokenAuthenticator) CreateVerificationTokens(tokenType TokenType) (*VerificationTokens, error) {
	// Generate verification Token and OTP
	plainToken, err := tokenAuthenticator.GenerateVerificationToken() // TODO: nell'handler gestire il retry nel caso non dovesse essere unico
	if err != nil {
		return nil, err
	}

	plainOTP, err := tokenAuthenticator.GenerateOTP()
	if err != nil {
		return nil, err
	}

	// Hash verification Token and OTP
	hashedToken := tokenAuthenticator.HashToken(plainToken)
	hashedOTP := tokenAuthenticator.HashToken(plainOTP)

	return &VerificationTokens{
		TokenType:   tokenType,
		PlainToken:  plainToken,
		HashedToken: hashedToken,
		PlainOTP:    plainOTP,
		HashedOTP:   hashedOTP,
		OTPExp:      tokenAuthenticator.getExp(tokenType),
	}, nil
}

// - Utils -

func (tokenAuthenticator *TokenAuthenticator) getExp(tokenType TokenType) time.Duration {
	var exp time.Duration

	switch tokenType {

	case TokenEmailVerification:
		exp = tokenAuthenticator.OTP.LongExp

	case TokenPasswordReset:
	case TokenTwoFactorAuth:
		exp = tokenAuthenticator.OTP.BaseExp

	default:
		exp = tokenAuthenticator.OTP.BaseExp
	}

	return exp
}

// Generate
func (tokenAuthenticator *TokenAuthenticator) GenerateVerificationToken() (string, error) {
	buffer := make([]byte, tokenAuthenticator.MagicLink.ByteSize)

	if _, err := rand.Read(buffer); err != nil {
		return "", err
	}

	return base64.RawURLEncoding.EncodeToString(buffer), nil
}

func (tokenAuthenticator *TokenAuthenticator) GenerateOTP() (string, error) { // TODO RICORDA: per la verifica dell'otp si usa anche lo user_id nella richiesta (l'otp potrebbe non essere unico nel db)
	length := tokenAuthenticator.OTP.Length

	max := new(big.Int).Exp(big.NewInt(10), big.NewInt(int64(length)), nil) // Create a new *big.Int as 10^length (big.NewInt(10) ^ big.NewInt(int64(length)))

	otp, err := rand.Int(rand.Reader, max) // Max for length = 6 -> 1000000 (values in range[000000, 999999])
	if err != nil {
		return "", err
	}

	return fmt.Sprintf("%0*d", length, otp), nil // Format with [length] numbers/zeros
}

// Hash and Verification
func (tokenAuthenticator *TokenAuthenticator) HashToken(plainToken string) []byte {
	hash := sha256.Sum256([]byte(plainToken)) // TODO: aggiungere pepper (secret)?
	return hash[:]                            // From [32]byte to []byte
}

func (tokenAuthenticator *TokenAuthenticator) VerifyToken(plainToken string, hashedToken []byte) bool {
	hash := tokenAuthenticator.HashToken(plainToken)
	return bytes.Equal(hash, hashedToken) // TODO: usa subtle.ConstantTimeCompare(hash1, hash2) == 1
}
