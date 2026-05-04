package main

import (
	"database/sql"
	"testing"

	"github.com/Samu-Amy/Shokora/internal/auth"
	"github.com/google/uuid"
)

// Clear db
func clearTestDB(db *sql.DB) {
	tables := []string{
		"users",
		"user_settings",
		"user_sessions",
		"refresh_tokens",
		"reset_session_tokens",
		"verification_tokens",
		"products",
	}

	for _, table := range tables {
		_, err := db.Exec("TRUNCATE TABLE " + table + " RESTART IDENTITY CASCADE")
		if err != nil {
			panic("failed to clean table " + table + ": " + err.Error())
		}
	}
}

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
