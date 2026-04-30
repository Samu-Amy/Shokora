package main

import (
	"database/sql"
	"testing"
	"time"
)

func seedSessions(t *testing.T, db *sql.DB, users []User) map[int64]int64 {
	t.Helper()

	query := `
		INSERT INTO user_sessions (user_id, expires_at)
		VALUES ($1, $2)
		RETURNING id
	`

	sessions := make(map[int64]int64)

	// Create users in db
	for _, user := range users {

		var sessionId int64

		// Create user
		err := db.QueryRowContext(
			t.Context(),
			query,
			user.Id,
			time.Now().Add(24*time.Hour),
		).Scan(
			&sessionId,
		)
		if err != nil {
			t.Fatalf("failed to insert user: %v", err)
		}

		sessions[user.Id] = sessionId
	}

	return sessions
}
