package store

import (
	"context"
	"database/sql"
	"errors"
	"time"

	"github.com/Samu-Amy/Shokora/internal/auth"
	"github.com/Samu-Amy/Shokora/internal/errorcodes"
)

type PostgresVTokensStore struct {
	db *sql.DB
}

func NewPostgresVTokenStore(db *sql.DB) *PostgresVTokensStore {
	return &PostgresVTokensStore{db: db}
}

// ----- CREATE -----

func (store *PostgresVTokensStore) CreateTokens(ctx context.Context, userId int64, verificationTokens *auth.VerificationTokens) (*int64, error) {
	// if pair (user_id, verification_type) exists -> update (set) columns with new values (tokens, exps) and reset otp attempts
	// else create new row
	query := `
		INSERT INTO verification_tokens (user_id, verification_type, magic_link_token, magic_link_token_exp, otp, otp_exp)
		VALUES ($1, $2, $3, $4, $5, $6)
		ON CONFLICT (user_id, verification_type)
		DO UPDATE SET magic_link_token = $3, magic_link_token_exp = $4, otp = $5, otp_exp = $6, otp_attempts = 0
		RETURNING id
	`

	queryCtx, cancel := context.WithTimeout(ctx, medium_query_timeout)
	defer cancel()

	// Fix magic link exp
	var magicLinkExp any
	if verificationTokens.HashedMagicLinkToken != nil {
		magicLinkExp = time.Now().Add(verificationTokens.MagicLinkTokenExp)
	}

	// Create tokens
	var verificationId int64

	err := store.db.QueryRowContext(
		queryCtx,
		query,
		userId,
		verificationTokens.VerificationType,
		verificationTokens.HashedMagicLinkToken,
		magicLinkExp,
		verificationTokens.HashedOTP,
		time.Now().Add(verificationTokens.OTPExp),
	).Scan(
		&verificationId,
	)

	if err != nil {
		switch {
		case isPostgresErrorCode(err, UniqueViolationErr):
			return nil, errorcodes.InternalErrDuplicateToken
		default:
			return nil, err
		}
	}

	return &verificationId, nil
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
			return errorcodes.InternalErrDuplicateToken
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

// ----- VERIFY -----

func (store *PostgresVTokensStore) VerifyMagicLink(ctx context.Context, hashedToken []byte) (*MagicLinkTokenPayload, error) {
	query := `
		SELECT id, user_id, verification_type, magic_link_token_exp
		FROM verification_tokens
		WHERE magic_link_token = $1
	`

	queryCtx, cancel := context.WithTimeout(ctx, medium_query_timeout)
	defer cancel()

	var magicLinkTokenPayload MagicLinkTokenPayload

	err := store.db.QueryRowContext(
		queryCtx,
		query,
		hashedToken,
	).Scan(
		&magicLinkTokenPayload.VerificationId,
		&magicLinkTokenPayload.UserId,
		&magicLinkTokenPayload.VerificationType,
		&magicLinkTokenPayload.Exp,
	)

	if err != nil {
		switch {
		case errors.Is(err, sql.ErrNoRows):
			return nil, errorcodes.ErrNotFound
		default:
			return nil, err
		}
	}

	return &magicLinkTokenPayload, nil
}

func (store *PostgresVTokensStore) VerifyOTP(ctx context.Context, verificationId int64, hashedOTP []byte) (*OTPPayload, error) {
	query := `
		SELECT user_id, verification_type, otp_attempts, otp_exp
		FROM verification_tokens
		WHERE id = $1 AND otp = $2
	`

	queryCtx, cancel := context.WithTimeout(ctx, medium_query_timeout)
	defer cancel()

	var otpPayload OTPPayload

	err := store.db.QueryRowContext(
		queryCtx,
		query,
		verificationId,
		hashedOTP,
	).Scan(
		&otpPayload.UserId,
		&otpPayload.VerificationType,
		&otpPayload.Attempts,
		&otpPayload.Exp,
	)

	if err != nil {
		switch {
		case errors.Is(err, sql.ErrNoRows):
			return nil, errorcodes.ErrNotFound
		default:
			return nil, err
		}
	}

	return &otpPayload, nil
}

// ----- DELETE -----
func (store *PostgresVTokensStore) Delete(ctx context.Context, verificationId int64) error {
	query := `
	DELETE from verification_tokens WHERE id = $1
	`

	queryCtx, cancel := context.WithTimeout(ctx, medium_query_timeout)
	defer cancel()

	_, err := store.db.ExecContext(
		queryCtx,
		query,
		verificationId,
	)

	if err != nil {
		return err
	}

	return nil
}
