package user

import (
	"context"
	"database/sql"
)

// TODO: sistema (qua solo metodi "strettamente legati" a users (tabella) che fanno le query, poi quelli "composti" o con logica (es. retry) li si crea nel service usando questi)

type IUserRepository interface {
	Create(ctx context.Context, transaction *sql.Tx, user *User) error

	GetById(ctx context.Context, userId int64) (*User, error)
	GetByGoogleId(ctx context.Context, googleId string) (*User, error)
	GetByEmail(ctx context.Context, email string) (*User, error)
	GetByEmailForUpdate(ctx context.Context, transaction *sql.Tx, email string) (*User, error)
	GetUserVerificationDataByEmailForUpdate(ctx context.Context, transaction *sql.Tx, email string) (*UserVerificationData, error)
	GetPasswordForUpdate(ctx context.Context, transaction *sql.Tx, userId int64) ([]byte, error)

	UpdatePassword(ctx context.Context, transaction *sql.Tx, userId int64, hashedPassword []byte) error
	SetGoogleId(ctx context.Context, transaction *sql.Tx, userId int64, googleId string) error
	SetIsVerified(ctx context.Context, userId int64) error // Set is_verified to true
	// SetIsActive(ctx context.Context, userId int64, isActive bool) error // TODO: implementa (per bloccare/sbloccare users)

	Delete(ctx context.Context, userId int64) error
}
