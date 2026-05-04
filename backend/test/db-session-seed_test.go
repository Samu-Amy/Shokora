package main

import (
	"context"
	"database/sql"
	"time"

	"github.com/Samu-Amy/Shokora/internal/auth"
)

type Sessions = map[int64]RefreshToken

type RefreshToken struct {
	Id         int64
	PlainToken string
}

func seedRefreshTokens(ctx context.Context, db *sql.DB, users []User) (Sessions, error) {

	sessionQuery := `
		INSERT INTO user_sessions (user_id, expires_at)
		VALUES ($1, $2)
		RETURNING id
	`

	refreshTokenQuery := `
		INSERT INTO refresh_tokens (session_id, token_hash, expires_at, replaces)
		VALUES ($1, $2, $3, $4)
	`

	sessions := make(map[int64]RefreshToken)

	for _, user := range users {

		// Start transaction
		tx, err := db.BeginTx(ctx, nil)
		if err != nil {
			return nil, err
		}

		// Create Session
		var sessionId int64

		err = tx.QueryRowContext(
			ctx,
			sessionQuery,
			user.Id,
			time.Now().Add(24*time.Hour),
		).Scan(
			&sessionId,
		)
		if err != nil {
			tx.Rollback()
			return nil, err
		}

		// Create Refresh Token
		plainToken, err := auth.GenerateBase64Token(32)
		if err != nil {
			tx.Rollback()
			return nil, err
		}

		tokenHash := auth.HashBase64Token(plainToken)

		_, err = tx.ExecContext(
			ctx,
			refreshTokenQuery,
			sessionId,
			tokenHash,
			time.Now().Add(30*time.Minute).UTC(),
			nil,
		)
		if err != nil {
			tx.Rollback()
			return nil, err
		}

		err = tx.Commit()
		if err != nil {
			tx.Rollback()
			return nil, err
		}

		sessions[user.Id] = RefreshToken{
			Id:         sessionId,
			PlainToken: plainToken,
		}
	}

	return sessions, nil
}
