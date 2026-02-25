package authservice

import (
	"context"
	"database/sql"
	"time"

	"github.com/Samu-Amy/Shokora/internal/db"
	"github.com/Samu-Amy/Shokora/internal/store"
	r_token "github.com/Samu-Amy/Shokora/internal/store/refresh-token.go"
)

// ----- CREATE AND ROTATE REFRESH TOKENS -----

// Create
func (service *AuthService) CreateRefreshToken(ctx context.Context, session *store.UserSession, refreshToken *r_token.RefreshToken, sessionExp, tokenExp time.Duration) error {
	// TODO: fai transaction per creare sia sessione che refresh token

	return db.WithTransaction(service.db, ctx, func(tx *sql.Tx) error {
		// Create session
		err := service.userSessionRepo.Create(ctx, tx, session, sessionExp)
		if err != nil {
			return err
		}

		// Create token
		return service.refreshTokenRepo.CreateToken(ctx, tx, refreshToken, tokenExp)
	})

}

// Rotate // TODO: aggiorna a session + refresh token
// func (service *AuthService) RotateRefreshToken(ctx context.Context, oldHashedToken []byte, newRefreshToken *r_token.RefreshToken) error { // TODO: ritorna dati (es. expires_at per cookies)
// 	err := db.WithTransaction(service.db, ctx, func(tx *sql.Tx) error {

// 		// Get token
// 		oldRefreshToken, err := service.refreshTokenRepo.GetToken(ctx, tx, oldHashedToken)
// 		if err != nil {
// 			return err
// 		}

// 		// Validate token - Expired or revoked
// 		if oldRefreshToken.ExpiresAt.Before(time.Now()) || oldRefreshToken.RevokedAt != nil {
// 			return errorcodes.ErrInvalid
// 		}

// 		// Validate token - wrong RefreshToken, User or Session id
// 		if oldRefreshToken.UserId != newRefreshToken.UserId || oldRefreshToken.SessionId != newRefreshToken.SessionId || oldRefreshToken.Id != newRefreshToken.Replaces {
// 			return errorcodes.ErrInvalid
// 		}

// 		// TODO: implementa estensione expiry - usando costanti "SESSION_EXTENSION_<...>" (e nel token nuovo usa il conto di estensioni da quello vecchio (più eventualmente quella appena fatta))

// 		// Create new token
// 		// (create token using same session_id of the old one and using its id as replaces)
// 		err = service.refreshTokenRepo.CreateToken(ctx, tx, newRefreshToken)
// 		if err != nil {
// 			return err
// 		}

// 		if oldRefreshToken.Id == nil || newRefreshToken.CreatedAt == nil {
// 			return errorcodes.ErrInvalid
// 		}

// 		// Revoke (update) old token
// 		err = service.refreshTokenRepo.RevokeTokenById(ctx, tx, *oldRefreshToken.Id, *newRefreshToken.CreatedAt)
// 		if err != nil {
// 			// TODO: gestire token not found or already revoked (?)
// 			switch {
// 			case errors.Is(err, errorcodes.InternalErrNoRowsAffected):
// 				return errorcodes.InternalErrTokenNotFoundOrAlreadyRevoked
// 			default:
// 				return err
// 			}
// 		}

// 		return nil
// 	})

// 	// Revoke session (in case of Reuse Detection)
// 	if errors.Is(err, errorcodes.InternalErrReusedToken) {
// 		err = service.refreshTokenRepo.DeleteSessionById(ctx, newRefreshToken.UserId, newRefreshToken.SessionId)
// 	}

// 	return err
// }
