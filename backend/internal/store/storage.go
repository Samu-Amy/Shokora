package store

import (
	"database/sql"

	"github.com/Samu-Amy/Shokora/internal/store/product"
	refreshtoken "github.com/Samu-Amy/Shokora/internal/store/refresh-token.go"
	"github.com/Samu-Amy/Shokora/internal/store/user"
	session "github.com/Samu-Amy/Shokora/internal/store/user-session"
	usersettings "github.com/Samu-Amy/Shokora/internal/store/user-settings"
	vtoken "github.com/Samu-Amy/Shokora/internal/store/verification-token"
)

/*
The Repository layer, it manages the database interactions using queries.
It is divided in Repositories (one for every db table) and it is used by the Service layer.
*/
type Storage struct {
	VToken       vtoken.IVTokenRepository
	UserSession  session.IUserSessionRepository
	RefreshToken refreshtoken.IRefreshTokenRepository
	User         user.IUserRepository
	UserSettings usersettings.IUserSettingsRepository
	Product      product.IProductRepository
}

func NewPostgresStorage(db *sql.DB) *Storage {
	return &Storage{
		VToken:       vtoken.NewPostgresStore(db),
		UserSession:  session.NewPostgresStore(db),
		RefreshToken: refreshtoken.NewPostgresStore(db),
		User:         user.NewPostgresStore(db),
		UserSettings: usersettings.NewPostgresStore(db),
		Product:      product.NewPostgresStore(db),
	}
}
