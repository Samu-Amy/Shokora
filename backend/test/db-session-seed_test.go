package main

import (
	"context"
	"database/sql"
	"time"
)

func seedSessions(ctx context.Context, db *sql.DB, users []User) (map[int64]int64, error) {

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
			ctx,
			query,
			user.Id,
			time.Now().Add(24*time.Hour),
		).Scan(
			&sessionId,
		)
		if err != nil {
			return nil, err
		}

		sessions[user.Id] = sessionId
	}

	return sessions, nil
}
