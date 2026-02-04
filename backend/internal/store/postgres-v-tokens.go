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

func (store *PostgresVTokensStore) CreateTokens(ctx context.Context, userId int64, verificationTokens *auth.VerificationTokens) (int64, error) {
	// if pair (user_id, verification_type) exists -> update (set) columns with new values (tokens, exps) and reset otp attempts
	// else create new row
	query := `
		INSERT INTO verification_tokens (user_id, verification_type, magic_link_token, magic_link_token_exp, otp, otp_exp)
		VALUES ($1, $2, $3, $4, $5, $6)
		ON CONFLICT (user_id, verification_type)
		DO UPDATE SET magic_link_token = $3, magic_link_token_exp = $4, otp = $5, otp_exp = $6, otp_attempts = 0
		RETURNING id
	` // TODO: sistema "conflict"

	queryCtx, cancel := context.WithTimeout(ctx, medium_query_timeout)
	defer cancel()

	var verificationId int64

	err := store.db.QueryRowContext(
		queryCtx,
		query,
		userId,
		verificationTokens.VerificationType,
		verificationTokens.HashedMagicLinkToken,
		time.Now().Add(verificationTokens.MagicLinkTokenExp),
		verificationTokens.HashedOTP,
		time.Now().Add(verificationTokens.OTPExp),
	).Scan(
		&verificationId,
	)

	if err != nil {
		switch {
		case isPostgresErrorCode(err, UniqueViolationErr):
			return -1, ErrDuplicateToken
		default:
			return -1, err
		}
	}

	return verificationId, nil
}

// ----- UPDATE -----

func (store *PostgresVTokensStore) UpdateMagicLinkTokenFromId(ctx context.Context, verificationId int64, magicLinkTokenHash []byte, magicLinkTokenExp time.Duration) error {
	query := `
		UPDATE verification_tokens
		SET magic_link_token = $1, magic_link_token_exp = $2
		WHERE id = $3
	`

	queryCtx, cancel := context.WithTimeout(ctx, medium_query_timeout)
	defer cancel()

	_, err := store.db.ExecContext(
		queryCtx,
		query,
		magicLinkTokenHash,
		time.Now().Add(magicLinkTokenExp),
		verificationId,
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

func (store *PostgresVTokensStore) UpdateOTPFromId(ctx context.Context, verificationId int64, otpHash []byte, otpExp time.Duration) error {
	query := `
		UPDATE verification_tokens
		SET otp = $1, otp_exp = $2, otp_attempts = 0
		WHERE id = $3
	`

	queryCtx, cancel := context.WithTimeout(ctx, medium_query_timeout)
	defer cancel()

	_, err := store.db.ExecContext(
		queryCtx,
		query,
		otpHash,
		time.Now().Add(otpExp),
		verificationId,
	)

	if err != nil {
		return err
	}

	return nil
}
