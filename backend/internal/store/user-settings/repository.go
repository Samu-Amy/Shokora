package usersettings

import (
	"context"
	"database/sql"
)

type IUserSettingsRepository interface {
	Create(ctx context.Context, transaction *sql.Tx, userId int64) (int64, error)
	Update(ctx context.Context, settings *UserSettings) error
	Delete(ctx context.Context, sessionId int64) error
}
