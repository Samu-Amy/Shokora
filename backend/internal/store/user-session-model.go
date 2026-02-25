package store

import (
	"context"
	"database/sql"
	"time"
)

type UserSession struct {
	Id        int64     `json:"id"` // Generated
	UserId    int64     `json:"user_id"`
	ExpiresAt time.Time `json:"expires_at"`
	CreatedAt time.Time `json:"created_at"` // Default now()
}

type UserSessionI interface {
	Create(ctx context.Context, transaction *sql.Tx, session *UserSession, sessionExp time.Duration) error
	Delete(ctx context.Context, transaction *sql.Tx, sessionId int64) error
}
