package vtoken

import (
	"context"
	"database/sql"
	"time"

	"github.com/Samu-Amy/Shokora/internal/auth"
	"github.com/Samu-Amy/Shokora/internal/database"
	"github.com/google/uuid"
)

type PostgresVTokenStore struct {
	db *sql.DB
}

func NewPostgresStore(db *sql.DB) *PostgresVTokenStore {
	return &PostgresVTokenStore{db: db}
}

// ----- CREATE -----

func (store *PostgresVTokenStore) Create(ctx context.Context, userId int64, verificationTokens *auth.VerificationTokens) error {
	// if pair (user_id, verification_type) exists -> update (set) columns with new values (tokens, exps) and reset otp attempts
	// else create new row
	query := `
		INSERT INTO verification_tokens (user_id, verification_type, magic_link_token_hash, magic_link_token_expires_at, otp_hash, otp_expires_at)
		VALUES ($1, $2, $3, $4, $5, $6)
		ON CONFLICT (user_id, verification_type)
		DO UPDATE SET magic_link_token_hash = $3, magic_link_token_expires_at = $4, otp_hash = $5, otp_expires_at = $6, otp_attempts = 0
		RETURNING id
	`

	queryCtx, cancel := context.WithTimeout(ctx, database.MediumQueryTimeout)
	defer cancel()

	// Fix magic link exp (nil if no magic link)
	var magicLinkExp any = nil
	if verificationTokens.HashedMagicLinkToken != nil {
		magicLinkExp = time.Now().Add(verificationTokens.MagicLinkTokenExp).UTC()
	}

	// Create tokens

	err := store.db.QueryRowContext(
		queryCtx,
		query,
		userId,
		verificationTokens.VerificationType,
		verificationTokens.HashedMagicLinkToken,
		magicLinkExp,
		verificationTokens.HashedOTP,
		time.Now().Add(verificationTokens.OTPExp).UTC(),
	).Scan(
		&verificationTokens.VerificationId,
	)

	return database.ParseDbError(err)
}

// ----- UPDATE -----

func (store *PostgresVTokenStore) UpdateMagicLinkTokenFromId(ctx context.Context, verificationId uuid.UUID, magicLinkTokenHash []byte, magicLinkTokenExp time.Duration) error {
	query := `
		UPDATE verification_tokens
		SET magic_link_token_hash = $1, magic_link_token_expires_at = $2
		WHERE id = $3
	`

	queryCtx, cancel := context.WithTimeout(ctx, database.MediumQueryTimeout)
	defer cancel()

	return database.HandleExecContextResult(store.db.ExecContext(
		queryCtx,
		query,
		magicLinkTokenHash,
		time.Now().Add(magicLinkTokenExp).UTC(),
		verificationId,
	))
}

func (store *PostgresVTokenStore) UpdateOTPFromId(ctx context.Context, verificationId uuid.UUID, otpHash []byte, otpExp time.Duration) error {
	query := `
		UPDATE verification_tokens
		SET otp_hash = $1, otp_expires_at = $2, otp_attempts = 0
		WHERE id = $3
	`

	queryCtx, cancel := context.WithTimeout(ctx, database.MediumQueryTimeout)
	defer cancel()

	return database.HandleExecContextResult(store.db.ExecContext(
		queryCtx,
		query,
		otpHash,
		time.Now().Add(otpExp).UTC(),
		verificationId,
	))
}

func (store *PostgresVTokenStore) IncrementOtpAttempts(ctx context.Context, transaction *sql.Tx, verificationId uuid.UUID, maxOTPAttempts uint8) error {
	query := `
		UPDATE verification_tokens
		SET otp_attempts = otp_attempts + 1
		WHERE id = $1 AND otp_attempts < $2
	`

	queryCtx, cancel := context.WithTimeout(ctx, database.MediumQueryTimeout)
	defer cancel()

	return database.HandleExecContextResult(transaction.ExecContext(
		queryCtx,
		query,
		verificationId,
		maxOTPAttempts,
	))
}

// ----- GET -----

func (store *PostgresVTokenStore) GetOtpData(ctx context.Context, transaction *sql.Tx, verificationId uuid.UUID, verificationType auth.VerificationType) (*OTPVerificationData, error) {
	query := `
		SELECT user_id, otp_hash, otp_attempts, otp_expires_at
		FROM verification_tokens
		WHERE id = $1 AND verification_type = $2
		FOR UPDATE
	`

	queryCtx, cancel := context.WithTimeout(ctx, database.MediumQueryTimeout)
	defer cancel()

	var otpPayload OTPVerificationData

	err := transaction.QueryRowContext(
		queryCtx,
		query,
		verificationId,
		verificationType, // Only for best practice
	).Scan(
		&otpPayload.UserId,
		&otpPayload.HashedOtp,
		&otpPayload.Attempts,
		&otpPayload.ExpiresAt,
	)

	if err != nil {
		return nil, database.ParseDbError(err)
	}

	return &otpPayload, nil
}

// ----- VERIFY -----

func (store *PostgresVTokenStore) GetValidMagicLinkData(ctx context.Context, transaction *sql.Tx, hashedToken []byte, verificationType auth.VerificationType) (*MagicLinkVerificationData, error) {
	query := `
		SELECT id, user_id
		FROM verification_tokens
		WHERE magic_link_token_hash = $1 AND verification_type = $2 AND magic_link_token_expires_at > NOW()
		FOR UPDATE
	`

	queryCtx, cancel := context.WithTimeout(ctx, database.MediumQueryTimeout)
	defer cancel()

	var magicLinkTokenPayload MagicLinkVerificationData

	err := transaction.QueryRowContext(
		queryCtx,
		query,
		hashedToken,
		verificationType, // Only for best practice
	).Scan(
		&magicLinkTokenPayload.VerificationId,
		&magicLinkTokenPayload.UserId,
	)

	if err != nil {
		return nil, database.ParseDbError(err)
	}

	return &magicLinkTokenPayload, nil
}

// ----- DELETE -----
func (store *PostgresVTokenStore) Delete(ctx context.Context, transaction *sql.Tx, verificationId uuid.UUID) error {
	query := `
		DELETE from verification_tokens WHERE id = $1
	`

	queryCtx, cancel := context.WithTimeout(ctx, database.MediumQueryTimeout)
	defer cancel()

	return database.HandleExecContextResult(transaction.ExecContext(
		queryCtx,
		query,
		verificationId,
	))
}
