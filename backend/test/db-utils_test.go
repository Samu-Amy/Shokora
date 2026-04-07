package main

import (
	"testing"

	"github.com/Samu-Amy/Shokora/internal/auth"
	"github.com/google/uuid"
)

func getVerificationType(t *testing.T, verificationId *uuid.UUID) auth.VerificationType {
	t.Helper()

	query := `
		SELECT verification_type
		FROM verification_tokens
		WHERE id = $1
	`

	var verificationType auth.VerificationType

	err := db.QueryRowContext(
		t.Context(),
		query,
		verificationId,
	).Scan(
		&verificationType,
	)

	if err != nil {
		t.Errorf("error getting verification type: %v", err)
	}

	return verificationType
}
