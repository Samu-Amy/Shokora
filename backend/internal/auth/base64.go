package auth

import (
	"crypto/rand"
	"crypto/sha256"
	"encoding/base64"
)

// Generate Magic Link and Refresh Tokens
func GenerateBase64Token(size int) (*string, error) {
	buffer := make([]byte, size)

	if _, err := rand.Read(buffer); err != nil {
		return nil, err
	}

	token := base64.RawURLEncoding.EncodeToString(buffer)
	return &token, nil
}

// Hash Magic Link and Refresh Tokens
func HashBase64Token(plainMagicLinkToken *string) []byte {
	if plainMagicLinkToken == nil {
		return nil
	}
	hash := sha256.Sum256([]byte(*plainMagicLinkToken))
	return hash[:] // From [32]byte to []byte
}
