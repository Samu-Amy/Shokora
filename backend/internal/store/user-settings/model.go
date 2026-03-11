package usersettings

import (
	"time"
)

type UserSettings struct {
	Id               int64
	UserId           int64
	HasTwoFactorAuth bool
	CreatedAt        time.Time
	UpdatedAt        time.Time
}
