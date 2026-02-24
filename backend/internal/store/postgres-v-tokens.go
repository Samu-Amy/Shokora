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
		INSERT INTO verification_tokens (user_id, verification_type, magic_link_token_hash, magic_link_token_exp, otp_hash, otp_exp)
		VALUES ($1, $2, $3, NOW() + $4, $5, NOW() + $6)
		ON CONFLICT (user_id, verification_type)
		DO UPDATE SET magic_link_token_hash = $3, magic_link_token_exp = NOW() + $4, otp_hash = $5, otp_exp = NOW() + $6, otp_attempts = 0
		RETURNING id
	`

	queryCtx, cancel := context.WithTimeout(ctx, MEDIUM_QUERY_TIMEOUT)
	defer cancel()

	// Fix magic link exp (nil if no magic link)
	var magicLinkExp any
	if verificationTokens.HashedMagicLinkToken != nil {
		magicLinkExp = verificationTokens.MagicLinkTokenExp // TODO: controlla che la scadenza sia giusta
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
		verificationTokens.OTPExp,
	).Scan(
		&verificationId,
	)

	if err != nil {
		switch {
		case isPostgresErrorCode(err, UNIQUE_VIOLATION_ERROR):
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
		SET magic_link_token_hash = $1, magic_link_token_exp = NOW() + $2
		WHERE id = $3
	`

	queryCtx, cancel := context.WithTimeout(ctx, MEDIUM_QUERY_TIMEOUT)
	defer cancel()

	_, err := store.db.ExecContext(
		queryCtx,
		query,
		magicLinkTokenHash,
		magicLinkTokenExp,
		verificationId,
	)

	if err != nil {
		switch {
		case isPostgresErrorCode(err, UNIQUE_VIOLATION_ERROR):
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
		SET otp_hash = $1, otp_exp = NOW() + $2, otp_attempts = 0
		WHERE id = $3
	`

	queryCtx, cancel := context.WithTimeout(ctx, MEDIUM_QUERY_TIMEOUT)
	defer cancel()

	_, err := store.db.ExecContext(
		queryCtx,
		query,
		otpHash,
		otpExp,
		verificationId,
	)

	if err != nil {
		return err
	}

	return nil
}

func (store *PostgresVTokensStore) UpdateOtpAttempts(ctx context.Context, verificationId int64, maxOTPAttempts uint8) error {
	query := `
		UPDATE verification_tokens
		SET otp_attempts = otp_attempts + 1
		WHERE id = $1 AND otp_attempts < $2
	`

	queryCtx, cancel := context.WithTimeout(ctx, MEDIUM_QUERY_TIMEOUT)
	defer cancel()

	res, err := store.db.ExecContext(
		queryCtx,
		query,
		verificationId,
		maxOTPAttempts,
	)

	if err != nil {
		switch {
		case errors.Is(err, sql.ErrNoRows):
			return errorcodes.ErrInvalid // VerificationId is not valid
		default:
			return err
		}
	}

	rows, err := res.RowsAffected()
	if err != nil {
		return err
	}

	if rows == 0 {
		return errorcodes.ErrMaxAttemptsExceeded // No errors (verificationId is correct), but no rows modified -> max attempts exceeded
	}

	return nil
}

// ----- GET -----

func (store *PostgresVTokensStore) GetOtpData(ctx context.Context, verificationId int64, verificationType auth.VerificationType) (*OTPPayload, error) {
	query := `
		SELECT user_id, otp_hash, otp_attempts, otp_exp
		FROM verification_tokens
		WHERE id = $1 AND verification_type = $2
	`

	queryCtx, cancel := context.WithTimeout(ctx, MEDIUM_QUERY_TIMEOUT)
	defer cancel()

	var otpPayload OTPPayload

	err := store.db.QueryRowContext(
		queryCtx,
		query,
		verificationId,
		verificationType, // Only for best practice
	).Scan(
		&otpPayload.UserId,
		&otpPayload.HashedOtp,
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

// ----- VERIFY -----

func (store *PostgresVTokensStore) VerifyMagicLink(ctx context.Context, hashedToken []byte, verificationType auth.VerificationType) (*MagicLinkTokenPayload, error) {
	query := `
		SELECT id, user_id
		FROM verification_tokens
		WHERE magic_link_token_hash = $1 AND verification_type = $2 AND magic_link_token_exp > NOW()
	`

	queryCtx, cancel := context.WithTimeout(ctx, MEDIUM_QUERY_TIMEOUT)
	defer cancel()

	var magicLinkTokenPayload MagicLinkTokenPayload

	err := store.db.QueryRowContext(
		queryCtx,
		query,
		hashedToken,
		verificationType, // Only for best practice
	).Scan(
		&magicLinkTokenPayload.VerificationId,
		&magicLinkTokenPayload.UserId,
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

// ----- DELETE -----
func (store *PostgresVTokensStore) Delete(ctx context.Context, verificationId int64) error {
	query := `
	DELETE from verification_tokens WHERE id = $1
	`

	queryCtx, cancel := context.WithTimeout(ctx, MEDIUM_QUERY_TIMEOUT)
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
