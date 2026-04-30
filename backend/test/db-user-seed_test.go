package main

import (
	"database/sql"
	"strings"
	"testing"
	"time"

	authservice "github.com/Samu-Amy/Shokora/internal/service/auth"
	"github.com/Samu-Amy/Shokora/internal/store/user"
)

// TODO: aggiungi (in entrambi) utenti con GoogleId (?)

type User struct {
	Id            int64
	GoogleId      *int64
	FirstName     string
	LastName      string
	Email         string
	PlainPassword string
	Birthday      string
	Role          user.Role
	IsVerified    bool
	HasTwoAuth    bool
}

func seedUsers(t *testing.T, db *sql.DB) []User {
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

	users := make([]User, 0)

	// Create users in db
	for i := range min(seedUserNum, len(validFirstNames), len(validLastNames), len(validBirthdays), len(validEmails), len(validPasswords)) {

		// Start transaction
		tx, err := db.BeginTx(t.Context(), nil)
		if err != nil {
			t.Fatalf("failed to start transaction: %v", err)
		}

		// Create user and password
		user := User{
			FirstName:     strings.TrimSpace(validFirstNames[i]),
			LastName:      strings.TrimSpace(validLastNames[i]),
			Email:         strings.TrimSpace(validEmails[i]),
			PlainPassword: validPasswords[i],
			Birthday:      validBirthdays[i],
			IsVerified:    userVerified[i],
			Role:          userRoles[i],
			HasTwoAuth:    userTwoFactorAuth[i],
		}

		// Hash password
		hashedPssw, err := testService.Auth.HashPassword(strings.TrimSpace(user.PlainPassword))
		if err != nil {
			t.Errorf("Error hashing password: %v", err)
		}

		// Parse birthday
		var birthday time.Time
		if user.Birthday != "" {
			birthday, err = authservice.ConvertBirthdayToTime(strings.TrimSpace(user.Birthday))
			if err != nil {
				t.Errorf("Error converting birthday: %v", err)
			}
		}

		// Create user
		err = tx.QueryRowContext(
			t.Context(),
			userQuery,
			nil,
			user.FirstName,
			user.LastName,
			user.Email,
			hashedPssw,
			birthday,
			user.IsVerified,
			user.Role,
		).Scan(
			&user.Id,
		)
		if err != nil {
			tx.Rollback()
			t.Fatalf("failed to insert user: %v", err)
		}

		// Create User Settings
		_, err = tx.ExecContext(
			t.Context(),
			settingsQuery,
			user.Id,
			user.HasTwoAuth,
		)
		if err != nil {
			tx.Rollback()
			t.Fatalf("failed to insert user settings: %v", err)
		}

		err = tx.Commit()
		if err != nil {
			t.Fatalf("failed to commit transaction: %v", err)
		}

		users = append(users, user)
	}

	return users
}
