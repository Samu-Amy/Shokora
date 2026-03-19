package user

import (
	"context"
	"database/sql"
)

// TODO: sistema (qua solo metodi "strettamente legati" a users (tabella) che fanno le query, poi quelli "composti" o con logica (es. retry) li si crea nel service usando questi)

type IUserRepository interface {
	Create(ctx context.Context, transaction *sql.Tx, user *User) error

	GetById(ctx context.Context, userId int64) (*User, error)
	GetByEmail(ctx context.Context, email string) (*User, error)

	GetUserVerificationDataByEmail(ctx context.Context, transaction *sql.Tx, email string) (*UserVerificationData, error)

	GetPassword(ctx context.Context, transaction *sql.Tx, userId int64) ([]byte, error)

	UpdatePassword(ctx context.Context, transaction *sql.Tx, userId int64, hashedPassword []byte) error

	SetIsVerified(ctx context.Context, userId int64) error // Set is_verified to true
	// SetIsActive(ctx context.Context, userId int64, isActive bool) error // TODO: implementa (per bloccare/sbloccare users)

	Delete(ctx context.Context, userId int64) error
}
