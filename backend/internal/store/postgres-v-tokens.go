package store

import (
	"context"
	"database/sql"
	"time"

	"github.com/Samu-Amy/Shokora/internal/auth"
)

type PostgresVTokensStore struct {
	db *sql.DB
}

func NewPostgresVTokenStore(db *sql.DB) *PostgresVTokensStore {
	return &PostgresVTokensStore{db: db}
}

// ----- CREATE -----

func (store *PostgresVTokensStore) CreateTokens(ctx context.Context, userId int64, verificationTokens *auth.VerificationTokens) error {
	// if user_id and verification_type -> update (set) columns with new values (tokens, exps) and reset otp attempts
	// else create new row
	query := `
		INSERT INTO verification_tokens (user_id, verification_type, magic_link_token, magic_link_token_exp, otp, otp_exp)
		VALUES ($1, $2, $3, $4, $5, $6)
		ON CONFLICT (user_id, verification_type)
		DO UPDATE SET magic_link_token = $3, magic_link_token_exp = $4, otp = $5, otp_exp = $6, otp_attempts = 0
	`

	queryCtx, cancel := context.WithTimeout(ctx, medium_query_timeout)
	defer cancel()

	_, err := store.db.ExecContext(
		queryCtx,
		query,
		userId,
		verificationTokens.VerificationType,
		verificationTokens.HashedMagicLinkToken,
		time.Now().Add(verificationTokens.MagicLinkTokenExp),
		verificationTokens.HashedOTP,
		time.Now().Add(verificationTokens.OTPExp),
	)

	if err != nil {
		switch {
		case isPostgresErrorCode(err, UniqueViolationErr):
			return ErrDuplicateToken
		default:
			return err
		}
	}

	return nil
}

// ----- UPDATE -----

func (store *PostgresVTokensStore) UpdateMagicLinkToken(ctx context.Context, userId int64, verificationType auth.VerificationType, magicLinkTokenHash []byte, magicLinkTokenExp time.Duration) error {
	query := `
		UPDATE verification_tokens
		SET magic_link_token = $1, magic_link_token_exp = $2
		WHERE user_id = $3 AND verification_type = $4
	`

	queryCtx, cancel := context.WithTimeout(ctx, medium_query_timeout)
	defer cancel()

	_, err := store.db.ExecContext(
		queryCtx,
		query,
		magicLinkTokenHash,
		time.Now().Add(magicLinkTokenExp),
		userId,
		verificationType,
	)

	if err != nil {
		switch {
		case isPostgresErrorCode(err, UniqueViolationErr):
			return ErrDuplicateToken
		default:
			return err
		}
	}

	return nil
}

func (store *PostgresVTokensStore) UpdateOTP(ctx context.Context, userId int64, verificationType auth.VerificationType, otpHash []byte, otpExp time.Duration) error {
	query := `
		UPDATE verification_tokens
		SET otp = $1, otp_exp = $2, otp_attempts = 0
		WHERE user_id = $3 AND verification_type = $4
	`

	queryCtx, cancel := context.WithTimeout(ctx, medium_query_timeout)
	defer cancel()

	_, err := store.db.ExecContext(
		queryCtx,
		query,
		otpHash,
		time.Now().Add(otpExp),
		userId,
		verificationType,
	)

	if err != nil {
		return err
	}

	return nil
}
