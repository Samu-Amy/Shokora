package store

import (
	"database/sql"

	oauthstate "github.com/Samu-Amy/Shokora/internal/store/oauth-states"
	"github.com/Samu-Amy/Shokora/internal/store/product"
	rtoken "github.com/Samu-Amy/Shokora/internal/store/refresh-token"
	rstoken "github.com/Samu-Amy/Shokora/internal/store/reset-session-tokens"
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
	VToken            vtoken.IVTokenRepository
	ResetSessionToken rstoken.IResetSessionTokenRepository
	OAuthState        oauthstate.IOAuthStateRepository
	RefreshToken      rtoken.IRefreshTokenRepository
	UserSession       session.IUserSessionRepository
	User              user.IUserRepository
	UserSettings      usersettings.IUserSettingsRepository
	Product           product.IProductRepository
}

func NewPostgresStorage(db *sql.DB) *Storage {
	return &Storage{
		VToken:            vtoken.NewPostgresStore(db),
		ResetSessionToken: rstoken.NewPostgresStore(db),
		OAuthState:        oauthstate.NewPostgresStore(db),
		RefreshToken:      rtoken.NewPostgresStore(db),
		UserSession:       session.NewPostgresStore(db),
		User:              user.NewPostgresStore(db),
		UserSettings:      usersettings.NewPostgresStore(db),
		Product:           product.NewPostgresStore(db),
	}
}
