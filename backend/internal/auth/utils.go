package auth

import (
	"crypto/rand"
	"encoding/base64"
)

// Generate Magic Link and Refresh Tokens
func GenerateToken(size int) (*string, error) {
	buffer := make([]byte, size)

	if _, err := rand.Read(buffer); err != nil {
		return nil, err
	}

	token := base64.RawURLEncoding.EncodeToString(buffer)
	return &token, nil
}
