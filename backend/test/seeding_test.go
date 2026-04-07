package main

import (
	"context"
	"database/sql"
	"strings"
	"testing"
	"time"

	authservice "github.com/Samu-Amy/Shokora/internal/service/auth"
)

// Create data (seed db)
func seedUsers(t *testing.T, db *sql.DB) {
	t.Helper()

	query := `
			INSERT INTO users (google_id, first_name, last_name, email, password, birthday, is_verified)
			VALUES ($1, $2, $3, $4, $5, $6, $7)
		`

	// Create users in db
	for i := range min(seedUserNum, len(validFirstNames), len(validLastNames), len(validBirthdays), len(validEmails), len(validPasswords)) {

		firstName := validFirstNames[i]
		lastName := validLastNames[i]
		strBirthday := validBirthdays[i]
		email := validEmails[i]
		pssw := validPasswords[i]

		hashedPssw, err := testService.Auth.HashPassword(strings.TrimSpace(pssw))
		if err != nil {
			t.Errorf("Error hashing password: %v", err)
		}

		var birthday time.Time
		if strBirthday != "" {
			birthday, err = authservice.ConvertBirthdayToTime(strings.TrimSpace(strBirthday))
			if err != nil {
				t.Errorf("Error converting birthday: %v", err)
			}
		}

		_, err = db.ExecContext(
			context.Background(),
			query,
			nil,
			firstName,
			lastName,
			email,
			hashedPssw,
			birthday,
			false,
		)
	}
}

func seedUsersFuzz(f *testing.F, db *sql.DB) {
	f.Helper()

	query := `
			INSERT INTO users (google_id, first_name, last_name, email, password, birthday, is_verified)
			VALUES ($1, $2, $3, $4, $5, $6, $7)
		`

	// Create users in db
	for i := range min(seedUserNum, len(validFirstNames), len(validLastNames), len(validBirthdays), len(validEmails), len(validPasswords)) {

		firstName := validFirstNames[i]
		lastName := validLastNames[i]
		strBirthday := validBirthdays[i]
		email := validEmails[i]
		pssw := validPasswords[i]

		hashedPssw, err := testService.Auth.HashPassword(strings.TrimSpace(pssw))
		if err != nil {
			f.Errorf("Error hashing password: %v", err)
		}

		var birthday time.Time
		if strBirthday != "" {
			birthday, err = authservice.ConvertBirthdayToTime(strings.TrimSpace(strBirthday))
			if err != nil {
				f.Errorf("Error converting birthday: %v", err)
			}
		}

		_, err = db.ExecContext(
			context.Background(),
			query,
			nil,
			firstName,
			lastName,
			email,
			hashedPssw,
			birthday,
			false,
		)
	}
}
