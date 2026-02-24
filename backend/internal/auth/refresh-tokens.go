package auth

import (
	"time"
)

// Used to create a new token or a one that replaces another (in -> data into db, out -> data from the db)
type RefreshToken struct {
	Id          *int64
	UserId      int64
	HashedToken []byte
	Exp         time.Duration // the expiration duration for the token (from config), used in db to set expires_at
	SessionExp  time.Duration // the expiration date of the session
	Replaces    *int64        // in
	RevokedAt   *time.Time    // out
	ExpiresAt   *time.Time    // out (expiration date, correspond to the "expires_at" set in db (using "Exp"))
	CreatedAt   *time.Time    // out
}

type CreateRefreshTokenPayload struct {
	PlainToken string
	ExpiresAt  time.Time
}
