package auth

// type JWTAuthenticatorI interface {
// 	GenerateJWTToken(claims jwt.Claims) (string, error)
// 	ValidateJWTToken(token string) (*jwt.Token, error)
// }

// type TokenAuthenticatorI interface {
// 	// Create a struct with verification Tokens (plain and hashed), verification type and otp expiry
// 	CreateVerificationTokens(verificationType VerificationType) (*VerificationTokens, error)

// 	// Generate plain Tokens
// 	GenerateMagicLinkToken() (string, error)
// 	GenerateOTP() (string, error)

// 	// Hash Tokens
// 	HashMagicLinkToken(plainMagicLinkToken string) []byte
// 	HashOTP(plainOTP string, verificationType VerificationType) []byte

// 	// Verify Tokens
// 	VerifyMagicLinkToken(plainToken string, hashedToken []byte) bool
// 	VerifyOTP(plainOTP string, hashedToken []byte, verificationType VerificationType) bool

// 	// Utils
// 	getExpiration(verificationType VerificationType) time.Duration
// 	getVerificationTypeString(verificationType VerificationType) string
// }
