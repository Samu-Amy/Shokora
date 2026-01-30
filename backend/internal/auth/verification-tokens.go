package auth

import (
	"crypto/hmac"
	"crypto/rand"
	"crypto/sha256"
	"crypto/subtle"
	"encoding/base64"
	"fmt"
	"math/big"
	"time"
)

// TODO: forse molti metodi si possono rendere privati (es. quelli per generate e hash tokens)

// - Authenticator -
type TokenAuthenticator struct {
	MagicLink  MagicLinkConfig
	OTP        OTPConfig
	MaxRetries uint8
	secret     string
}

type MagicLinkConfig struct {
	ByteSize int
	Exp      time.Duration
}

type OTPConfig struct {
	Length      uint8
	MaxAttempts uint8
	LongExp     time.Duration // For Email Verification
	BaseExp     time.Duration // For Password Reset and 2FA
	// CriticalExp time.Duration // For critical operations (es. 30s)
}

// - Tokens -
type VerificationTokens struct {
	VerificationType     VerificationType
	PlainMagicLinkToken  string
	HashedMagicLinkToken []byte
	PlainOTP             string
	HashedOTP            []byte
	MagicLinkTokenExp    time.Duration
	OTPExp               time.Duration
}

// Verification Token and OTP
type VerificationType uint8

const (
	TokenEmailVerification VerificationType = 0
	TokenPasswordReset     VerificationType = 1
	TokenTwoFactorAuth     VerificationType = 2
)

// - Constructor -

func NewTokenAuthenticator(MagicLink MagicLinkConfig, OTP OTPConfig, MaxRetries uint8, secret string) *TokenAuthenticator {
	return &TokenAuthenticator{MagicLink, OTP, MaxRetries, secret}
}

// - Methods -

func (tokenAuthenticator *TokenAuthenticator) CreateVerificationTokens(verificationType VerificationType) (*VerificationTokens, error) {
	// Generate verification Token and OTP
	plainMagicLinkToken, err := tokenAuthenticator.generateMagicLinkToken()
	if err != nil {
		return nil, err
	}

	plainOTP, err := tokenAuthenticator.generateOTP()
	if err != nil {
		return nil, err
	}

	// Hash verification Token and OTP
	hashedMagicLinkToken := tokenAuthenticator.hashMagicLinkToken(plainMagicLinkToken)
	hashedOTP := tokenAuthenticator.hashOTP(plainOTP, verificationType)

	return &VerificationTokens{
		VerificationType:     verificationType,
		PlainMagicLinkToken:  plainMagicLinkToken,
		HashedMagicLinkToken: hashedMagicLinkToken,
		PlainOTP:             plainOTP,
		HashedOTP:            hashedOTP,
		MagicLinkTokenExp:    tokenAuthenticator.MagicLink.Exp,
		OTPExp:               tokenAuthenticator.getExpiration(verificationType),
	}, nil
}

// Regenerate
func (tokenAuthenticator *TokenAuthenticator) RegenerateMagicLinkToken(verificationTokens *VerificationTokens) error {
	newMagicLinkToken, err := tokenAuthenticator.generateMagicLinkToken()
	if err != nil {
		return err
	}

	verificationTokens.PlainMagicLinkToken = newMagicLinkToken
	verificationTokens.HashedMagicLinkToken = tokenAuthenticator.hashMagicLinkToken(newMagicLinkToken)

	return nil
}

func (tokenAuthenticator *TokenAuthenticator) RegenerateOTP(verificationTokens *VerificationTokens) error {
	newOTP, err := tokenAuthenticator.generateOTP()
	if err != nil {
		return err
	}

	verificationTokens.PlainOTP = newOTP
	verificationTokens.HashedOTP = tokenAuthenticator.hashOTP(newOTP, verificationTokens.VerificationType)

	return nil
}

// Verification
func (tokenAuthenticator *TokenAuthenticator) VerifyMagicLinkToken(plainToken string, hashedToken []byte) bool {
	hash := tokenAuthenticator.hashMagicLinkToken(plainToken)

	//? si può aggiungere padding al/ai token e fare comunque il controllo per evitare timing leak (ma rischiando di validare token sbagliati)
	if len(hash) != 32 || len(hashedToken) != 32 {
		return false
	}

	return subtle.ConstantTimeCompare(hash, hashedToken) == 1
}

func (tokenAuthenticator *TokenAuthenticator) VerifyOTP(plainOTP string, hashedToken []byte, verificationType VerificationType) bool {
	hash := tokenAuthenticator.hashOTP(plainOTP, verificationType)
	return hmac.Equal(hash, hashedToken)
}

// ----- ----- ----- PRIVATES ----- ----- -----

// Generate
func (tokenAuthenticator *TokenAuthenticator) generateMagicLinkToken() (string, error) {
	buffer := make([]byte, tokenAuthenticator.MagicLink.ByteSize)

	if _, err := rand.Read(buffer); err != nil {
		return "", err
	}

	return base64.RawURLEncoding.EncodeToString(buffer), nil
}

func (tokenAuthenticator *TokenAuthenticator) generateOTP() (string, error) { // TODO RICORDA: per la verifica dell'otp si usa anche lo user_id nella richiesta (l'otp potrebbe non essere unico nel db)
	length := tokenAuthenticator.OTP.Length

	max := new(big.Int).Exp(big.NewInt(10), big.NewInt(int64(length)), nil) // Create a new *big.Int as 10^length (big.NewInt(10) ^ big.NewInt(int64(length)))

	otp, err := rand.Int(rand.Reader, max) // Max for length = 6 -> 1000000 (values in range[000000, 999999])
	if err != nil {
		return "", err
	}

	return fmt.Sprintf("%0*d", length, otp), nil // Format with [length] numbers/zeros
}

// Hash
func (tokenAuthenticator *TokenAuthenticator) hashMagicLinkToken(plainMagicLinkToken string) []byte {
	hash := sha256.Sum256([]byte(plainMagicLinkToken))
	return hash[:] // From [32]byte to []byte
}

func (tokenAuthenticator *TokenAuthenticator) hashOTP(plainOTP string, verificationType VerificationType) []byte {
	mac := hmac.New(sha256.New, []byte(tokenAuthenticator.secret))
	mac.Write([]byte(plainOTP + tokenAuthenticator.getVerificationTypeString(verificationType)))
	return mac.Sum(nil)
}

// - Utils -

func (tokenAuthenticator *TokenAuthenticator) getExpiration(verificationType VerificationType) time.Duration {
	var exp time.Duration

	switch verificationType {

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

func (tokenAuthenticator *TokenAuthenticator) getVerificationTypeString(verificationType VerificationType) string {
	verificationTypeString := "verification"

	switch verificationType {

	case TokenEmailVerification:
		verificationTypeString = "email_verification"

	case TokenPasswordReset:
		verificationTypeString = "password_reset"

	case TokenTwoFactorAuth:
		verificationTypeString = "two_factor_auth"
	}

	return verificationTypeString
}
