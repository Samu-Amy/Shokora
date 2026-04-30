package main

import (
	"context"
	"database/sql"
)

type AuthState struct {
	Users    []User
	Sessions map[int64]int64 // UserId -> SessionId
}

func seedAuthState(ctx context.Context, db *sql.DB) (*AuthState, error) {

	users, err := seedUsers(ctx, db)
	if err != nil {
		return nil, err
	}

	sessions, err := seedSessions(ctx, db, users)
	if err != nil {
		return nil, err
	}

	return &AuthState{
		Users:    users,
		Sessions: sessions,
	}, nil
}
