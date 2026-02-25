package session

import (
	"context"
	"database/sql"
	"time"
)

type IUserSessionRepository interface {
	Create(ctx context.Context, transaction *sql.Tx, session *UserSession, sessionExp time.Duration) error
	Delete(ctx context.Context, transaction *sql.Tx, sessionId int64) error
}
