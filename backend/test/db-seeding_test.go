package main

import (
	"context"
	"database/sql"
	"strings"
	"testing"
	"time"

	authservice "github.com/Samu-Amy/Shokora/internal/service/auth"
)

// TODO: aggiungi (in entrambi) utenti con GoogleId (?)

// Create data (seed db)
func seedUsers(t *testing.T, db *sql.DB) {
	t.Helper()

	userQuery := `
		INSERT INTO users (google_id, first_name, last_name, email, password, birthday, is_verified, user_role)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8)
		RETURNING id
	`

	settingsQuery := `
		INSERT INTO user_settings (user_id, two_factor_auth)
		VALUES ($1, $2)
	`

	// Create users in db
	for i := range min(seedUserNum, len(validFirstNames), len(validLastNames), len(validBirthdays), len(validEmails), len(validPasswords)) {

		// Start transaction
		tx, err := db.BeginTx(t.Context(), nil)
		if err != nil {
			t.Fatalf("failed to start transaction: %v", err)
		}

		// Create birthday and password
		strBirthday := validBirthdays[i]
		pssw := validPasswords[i]

		// Hash password
		hashedPssw, err := testService.Auth.HashPassword(strings.TrimSpace(pssw))
		if err != nil {
			t.Errorf("Error hashing password: %v", err)
		}

		// Parse birthday
		var birthday time.Time
		if strBirthday != "" {
			birthday, err = authservice.ConvertBirthdayToTime(strings.TrimSpace(strBirthday))
			if err != nil {
				t.Errorf("Error converting birthday: %v", err)
			}
		}

		var userId int64

		// Create user
		err = tx.QueryRowContext(
			t.Context(),
			userQuery,
			nil,
			strings.TrimSpace(validFirstNames[i]),
			strings.TrimSpace(validLastNames[i]),
			strings.TrimSpace(validEmails[i]),
			hashedPssw,
			birthday,
			userVerified[i],
			userRoles[i],
		).Scan(
			&userId,
		)
		if err != nil {
			tx.Rollback()
			t.Fatalf("failed to insert user: %v", err)
		}

		// Create User Settings
		_, err = tx.ExecContext(
			context.Background(),
			settingsQuery,
			userId,
			userTwoFactorAuth[i],
		)
		if err != nil {
			tx.Rollback()
			t.Fatalf("failed to insert user settings: %v", err)
		}

		err = tx.Commit()
		if err != nil {
			t.Fatalf("failed to commit transaction: %v", err)
		}
	}
}

func seedUsersFuzz(f *testing.F, db *sql.DB) {
	f.Helper()

	userQuery := `
		INSERT INTO users (google_id, first_name, last_name, email, password, birthday, is_verified, user_role)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8)
		RETURNING id
	`

	settingsQuery := `
		INSERT INTO user_settings (user_id, two_factor_auth)
		VALUES ($1, $2)
	`

	// Create users in db
	for i := range min(seedUserNum, len(validFirstNames), len(validLastNames), len(validBirthdays), len(validEmails), len(validPasswords)) {

		// Start transaction
		tx, err := db.BeginTx(context.Background(), nil)
		if err != nil {
			f.Fatalf("failed to start transaction: %v", err)
		}

		// Create birthday and password
		strBirthday := validBirthdays[i]
		pssw := validPasswords[i]

		// Hash password
		hashedPssw, err := testService.Auth.HashPassword(strings.TrimSpace(pssw))
		if err != nil {
			f.Errorf("Error hashing password: %v", err)
		}

		// Parse birthday
		var birthday time.Time
		if strBirthday != "" {
			birthday, err = authservice.ConvertBirthdayToTime(strings.TrimSpace(strBirthday))
			if err != nil {
				f.Errorf("Error converting birthday: %v", err)
			}
		}

		var userId int64

		// Create user
		err = tx.QueryRowContext(
			context.Background(),
			userQuery,
			nil,
			strings.TrimSpace(validFirstNames[i]),
			strings.TrimSpace(validLastNames[i]),
			strings.TrimSpace(validEmails[i]),
			hashedPssw,
			birthday,
			userVerified[i],
			userRoles[i],
		).Scan(
			&userId,
		)
		if err != nil {
			tx.Rollback()
			f.Fatalf("failed to insert user: %v", err)
		}

		// Create User Settings
		_, err = tx.ExecContext(
			context.Background(),
			settingsQuery,
			userId,
			userTwoFactorAuth[i],
		)
		if err != nil {
			tx.Rollback()
			f.Fatalf("failed to insert user settings: %v", err)
		}

		err = tx.Commit()
		if err != nil {
			f.Fatalf("failed to commit transaction: %v", err)
		}
	}
}
