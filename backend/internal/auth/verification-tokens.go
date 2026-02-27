package auth

import (
	"crypto/hmac"
	"crypto/rand"
	"crypto/sha256"
	"fmt"
	"math/big"
	"time"

	"github.com/Samu-Amy/Shokora/internal/config"
)

// TODO: forse molti metodi si possono rendere privati (es. quelli per generate e hash tokens)

// - Authenticator -
type TokenAuthenticator struct {
	MagicLink  config.MagicLinkConfig
	OTP        config.OTPConfig
	MaxRetries uint8 // Counting the first attempt
	secret     string
}

// - Tokens -
type VerificationTokens struct {
	VerificationType     VerificationType
	PlainMagicLinkToken  *string
	HashedMagicLinkToken []byte
	PlainOTP             string
	HashedOTP            []byte
	MagicLinkTokenExp    time.Duration
	OTPExp               time.Duration
}

// Verification Token and OTP
type VerificationType uint8

const (
	EmailVerification VerificationType = 0
	PasswordReset     VerificationType = 1
	TwoFactorAuth     VerificationType = 2
)

// - Constructor -

func NewTokenAuthenticator(MagicLink config.MagicLinkConfig, OTP config.OTPConfig, MaxRetries uint8, secret string) *TokenAuthenticator {
	return &TokenAuthenticator{MagicLink, OTP, MaxRetries, secret}
}

// - Methods -

func (tokenAuthenticator *TokenAuthenticator) CreateVerificationTokens(verificationType VerificationType) (*VerificationTokens, error) {

	// Generate and hash verification Token (only for email verification and password reset)
	var plainMagicLinkToken *string = nil
	var hashedMagicLinkToken []byte = nil
	var err error

	if verificationType != TwoFactorAuth {
		plainMagicLinkToken, err = GenerateBase64Token(tokenAuthenticator.MagicLink.ByteSize)
		if err != nil {
			return nil, err
		}

		hashedMagicLinkToken = HashBase64Token(plainMagicLinkToken)
	}

	// Generate and hash OTP
	plainOTP, err := tokenAuthenticator.generateOTP()
	if err != nil {
		return nil, err
	}

	hashedOTP := tokenAuthenticator.HashOTP(plainOTP, verificationType)

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

func (tokenAuthenticator *TokenAuthenticator) HashOTP(plainOTP string, verificationType VerificationType) []byte {
	mac := hmac.New(sha256.New, []byte(tokenAuthenticator.secret))
	mac.Write([]byte(plainOTP + tokenAuthenticator.getVerificationTypeString(verificationType)))
	return mac.Sum(nil)
}

// Regenerate
func (tokenAuthenticator *TokenAuthenticator) RegenerateMagicLinkToken(verificationTokens *VerificationTokens) error {
	newMagicLinkToken, err := GenerateBase64Token(tokenAuthenticator.MagicLink.ByteSize)
	if err != nil {
		return err
	}

	verificationTokens.PlainMagicLinkToken = newMagicLinkToken
	verificationTokens.HashedMagicLinkToken = HashBase64Token(newMagicLinkToken)

	return nil
}

func (tokenAuthenticator *TokenAuthenticator) RegenerateOTP(verificationTokens *VerificationTokens) error {
	newOTP, err := tokenAuthenticator.generateOTP()
	if err != nil {
		return err
	}

	verificationTokens.PlainOTP = newOTP
	verificationTokens.HashedOTP = tokenAuthenticator.HashOTP(newOTP, verificationTokens.VerificationType)

	return nil
}

// Verification
// func (tokenAuthenticator *TokenAuthenticator) VerifyMagicLinkToken(plainToken *string, hashedToken []byte) bool {
// 	if plainToken == nil || hashedToken == nil { // theoretically hashedToken should be also checked from len(hashedToken), since would return 0 if nil (so is redundant)
// 		return false
// 	}
// 	hash := tokenAuthenticator.HashMagicLinkToken(plainToken)

// 	//? si può aggiungere padding al/ai token e fare comunque il controllo per evitare timing leak (ma rischiando di validare token sbagliati)
// 	if len(hash) != 32 || len(hashedToken) != 32 {
// 		return false
// 	}

// 	return subtle.ConstantTimeCompare(hash, hashedToken) == 1
// }

// Verification
func (tokenAuthenticator *TokenAuthenticator) VerifyOTP(hashedOtp1 []byte, hashedOtp2 []byte) bool {
	return hmac.Equal(hashedOtp1, hashedOtp2)
}

// ----- ----- ----- PRIVATES ----- ----- -----

func (tokenAuthenticator *TokenAuthenticator) generateOTP() (string, error) { // TODO RICORDA: per la verifica dell'otp si usa anche lo user_id nella richiesta (l'otp potrebbe non essere unico nel db)
	length := tokenAuthenticator.OTP.Length

	max := new(big.Int).Exp(big.NewInt(10), big.NewInt(int64(length)), nil) // Create a new *big.Int as 10^length (big.NewInt(10) ^ big.NewInt(int64(length)))

	otp, err := rand.Int(rand.Reader, max) // Max for length = 6 -> 1000000 (values in range[000000, 999999])
	if err != nil {
		return "", err
	}

	return fmt.Sprintf("%0*d", length, otp), nil // Format with [length] numbers/zeros
}

// - Utils -

func (tokenAuthenticator *TokenAuthenticator) getExpiration(verificationType VerificationType) time.Duration {
	var exp time.Duration

	switch verificationType {

	case EmailVerification, PasswordReset:
		exp = tokenAuthenticator.OTP.LongExp

	case TwoFactorAuth:
		exp = tokenAuthenticator.OTP.BaseExp

	default:
		exp = tokenAuthenticator.OTP.BaseExp
	}

	return exp
}

func (tokenAuthenticator *TokenAuthenticator) getVerificationTypeString(verificationType VerificationType) string {
	verificationTypeString := "verification"

	switch verificationType {

	case EmailVerification:
		verificationTypeString = "email_verification"

	case PasswordReset:
		verificationTypeString = "password_reset"

	case TwoFactorAuth:
		verificationTypeString = "two_factor_auth"
	}

	return verificationTypeString
}
