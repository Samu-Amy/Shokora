package store

// import (
// 	"context"
// 	"database/sql"
// 	"errors"
// 	"time"

// 	"github.com/Samu-Amy/Shokora/internal/errorcodes"
// )

// TODO: fai lookup anche per verification_type (token + verification_type | OPT + email + verification_type + attempts) -- SET otp_attempts = otp_attempts + 1 (aggiorna atomicamente attempts)

// // ----- PRIVATES -----

// TODO: sposta in vtokens
// func (store *PostgresUserStore) getUserFromEmailVerificationToken(ctx context.Context, transaction *sql.Tx, hashedToken string) (*User, error) {
// 	query := `
// 	SELECT u.id, u.first_name, u.last_name, u.email, u.is_verified, u.created_at, u.updated_at
// 	FROM users u
// 	JOIN email_verification_tokens e ON u.id = e.user_id
// 	WHERE e.token = $1 AND e.expiry > $2
// 	`
// TODO: nel caso sia scaduto (bisogna fare un controllo separato ed eliminare sia token che user - mandare un errore ErrExpired)?

// 	queryCtx, cancel := context.WithTimeout(ctx, MEDIUM_QUERY_TIMEOUT)
// 	defer cancel()

// 	user := &User{}

// 	err := transaction.QueryRowContext(
// 		queryCtx,
// 		query,
// 		hashedToken,
// 		time.Now(),
// 	).Scan(
// 		&user.Id,
// 		&user.FirstName,
// 		&user.LastName,
// 		&user.Email,
// 		&user.IsVerified,
// 		&user.CreatedAt,
// 		&user.UpdatedAt,
// 	)

// 	if err != nil {
// 		switch {
// 		case errors.Is(err, sql.ErrNoRows):
// 			return nil, errorcodes.ErrNotFound
// 		default:
// 			return nil, err
// 		}
// 	}

// 	return user, nil
// }

// TODO: sposta in vtokens
// func (store *PostgresUserStore) deleteEmailVerificationToken(ctx context.Context, transaction *sql.Tx, userId int64) error {
// 	query := `
// 		DELETE FROM email_verification_tokens
// 		WHERE user_id = $1
// 	`
// 	queryCtx, cancel := context.WithTimeout(ctx, MEDIUM_QUERY_TIMEOUT)
// 	defer cancel()

// 	_, err := transaction.ExecContext(queryCtx, query, userId)
// 	if err != nil {
// TODO: migliorare error handling (?)
// 		return err
// 	}

// 	return nil
// }
