package session

import (
	"context"
	"database/sql"
	"time"
)

type IUserSessionRepository interface {
	Create(ctx context.Context, transaction *sql.Tx, userId int64, sessionExp time.Duration) (int64, error)

	Delete(ctx context.Context, sessionId int64) error

	DeleteExpired(ctx context.Context, userId int64) error
	DeleteOtherUserSessions(ctx context.Context, transaction *sql.Tx, userId, sessionId int64) error
}
