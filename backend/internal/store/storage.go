package store

import (
	"database/sql"

	"github.com/Samu-Amy/Shokora/internal/store/product"
	refreshtoken "github.com/Samu-Amy/Shokora/internal/store/refresh-token.go"
	"github.com/Samu-Amy/Shokora/internal/store/user"
	session "github.com/Samu-Amy/Shokora/internal/store/user-session"
	vtoken "github.com/Samu-Amy/Shokora/internal/store/verification-token"
)

type Storage struct {
	User         user.UserRepositoryI
	Product      product.ProductRepositoryI
	VToken       vtoken.VTokenRepositoryI
	UserSession  session.UserSessionI
	RefreshToken refreshtoken.RefreshTokenRepositoryI
}

func NewPostgresStorage(db *sql.DB) *Storage {
	return &Storage{
		User:         user.NewPostgresStore(db),
		Product:      product.NewPostgresStore(db),
		VToken:       vtoken.NewPostgresStore(db),
		UserSession:  session.NewPostgresStore(db),
		RefreshToken: refreshtoken.NewPostgresStore(db),
	}
}
