package usersettings

import (
	"time"
)

type UserSettings struct {
	Id               int64
	UserId           int64
	HasTwoFactorAuth bool
	Notifications    Notification
	CreatedAt        time.Time
	UpdatedAt        time.Time
}

// - Notification -

type Notification uint32

const (
	// Email
	Email Notification = 1 << iota // TODO: setta i tipi di notifiche

	// Push
	Push
)

func (settings *UserSettings) AddNotification(notification Notification) {
	settings.Notifications |= notification
}

func (settings *UserSettings) HasNotification(notification Notification) bool {
	return settings.Notifications&notification != 0
}

func (settings *UserSettings) RemoveNotification(notification Notification) {
	settings.Notifications &^= notification
}
