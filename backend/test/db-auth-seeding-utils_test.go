package main

import (
	"database/sql"
	"testing"
)

type AuthState struct {
	Users    []User
	Sessions map[int64]int64 // UserId -> SessionId
}

func seedAuthState(t *testing.T, db *sql.DB) *AuthState {
	t.Helper()

	users := seedUsers(t, db)

	sessions := seedSessions(t, db, users)

	return &AuthState{
		Users:    users,
		Sessions: sessions,
	}
}
