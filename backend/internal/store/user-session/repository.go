package session

import (
	"context"
	"database/sql"
	"time"
)

type IUserSessionRepository interface {
	Create(ctx context.Context, transaction *sql.Tx, userId int64, sessionExp time.Duration) (*UserSession, error)
	Delete(ctx context.Context, transaction *sql.Tx, sessionId int64) error
}
