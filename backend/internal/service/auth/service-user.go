package authservice

import (
	"context"
	"database/sql"

	"github.com/Samu-Amy/Shokora/internal/api/payloads"
	domerrors "github.com/Samu-Amy/Shokora/internal/errors/dom"
)

// - Update -

func (service *AuthService) UpdatePassword(ctx context.Context, userId int64, payload *payloads.UpdatePasswordReq) error {

	err := service.txManager.WithTx(ctx, func(tx *sql.Tx) error {

	})

	return domerrors.ParseIntError(err)
}
