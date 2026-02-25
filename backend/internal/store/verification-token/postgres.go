package v_token

import (
	"context"
	"database/sql"
	"time"

	"github.com/Samu-Amy/Shokora/internal/auth"
	"github.com/Samu-Amy/Shokora/internal/db"
)

type PostgresVTokenStore struct {
	db *sql.DB
}

func NewPostgresStore(db *sql.DB) *PostgresVTokenStore {
	return &PostgresVTokenStore{db: db}
}

// ----- CREATE -----

func (store *PostgresVTokenStore) CreateTokens(ctx context.Context, userId int64, verificationTokens *auth.VerificationTokens) (*int64, error) {
	// if pair (user_id, verification_type) exists -> update (set) columns with new values (tokens, exps) and reset otp attempts
	// else create new row
	query := `
		INSERT INTO verification_tokens (user_id, verification_type, magic_link_token_hash, magic_link_token_exp, otp_hash, otp_exp)
		VALUES ($1, $2, $3, NOW() + $4, $5, NOW() + $6)
		ON CONFLICT (user_id, verification_type)
		DO UPDATE SET magic_link_token_hash = $3, magic_link_token_exp = NOW() + $4, otp_hash = $5, otp_exp = NOW() + $6, otp_attempts = 0
		RETURNING id
	`

	queryCtx, cancel := context.WithTimeout(ctx, db.MEDIUM_QUERY_TIMEOUT)
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

	return &verificationId, db.ParseDbError(err)
}

// ----- UPDATE -----

func (store *PostgresVTokenStore) UpdateMagicLinkTokenFromId(ctx context.Context, verificationId int64, magicLinkTokenHash []byte, magicLinkTokenExp time.Duration) error {
	query := `
		UPDATE verification_tokens
		SET magic_link_token_hash = $1, magic_link_token_exp = NOW() + $2
		WHERE id = $3
	`

	queryCtx, cancel := context.WithTimeout(ctx, db.MEDIUM_QUERY_TIMEOUT)
	defer cancel()

	return db.HandleExecContextResult(store.db.ExecContext(
		queryCtx,
		query,
		magicLinkTokenHash,
		magicLinkTokenExp,
		verificationId,
	))
}

func (store *PostgresVTokenStore) UpdateOTPFromId(ctx context.Context, verificationId int64, otpHash []byte, otpExp time.Duration) error {
	query := `
		UPDATE verification_tokens
		SET otp_hash = $1, otp_exp = NOW() + $2, otp_attempts = 0
		WHERE id = $3
	`

	queryCtx, cancel := context.WithTimeout(ctx, db.MEDIUM_QUERY_TIMEOUT)
	defer cancel()

	return db.HandleExecContextResult(store.db.ExecContext(
		queryCtx,
		query,
		otpHash,
		otpExp,
		verificationId,
	))
}

func (store *PostgresVTokenStore) UpdateOtpAttempts(ctx context.Context, verificationId int64, maxOTPAttempts uint8) error {
	query := `
		UPDATE verification_tokens
		SET otp_attempts = otp_attempts + 1
		WHERE id = $1 AND otp_attempts < $2
	`

	queryCtx, cancel := context.WithTimeout(ctx, db.MEDIUM_QUERY_TIMEOUT)
	defer cancel()

	return db.HandleExecContextResult(store.db.ExecContext(
		queryCtx,
		query,
		verificationId,
		maxOTPAttempts,
	))
}

// ----- GET -----

func (store *PostgresVTokenStore) GetOtpData(ctx context.Context, verificationId int64, verificationType auth.VerificationType) (*OTPPayload, error) {
	query := `
		SELECT user_id, otp_hash, otp_attempts, otp_exp
		FROM verification_tokens
		WHERE id = $1 AND verification_type = $2
	`

	queryCtx, cancel := context.WithTimeout(ctx, db.MEDIUM_QUERY_TIMEOUT)
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

	return &otpPayload, db.ParseDbError(err)
}

// ----- VERIFY -----

func (store *PostgresVTokenStore) VerifyMagicLink(ctx context.Context, hashedToken []byte, verificationType auth.VerificationType) (*MagicLinkTokenPayload, error) {
	query := `
		SELECT id, user_id
		FROM verification_tokens
		WHERE magic_link_token_hash = $1 AND verification_type = $2 AND magic_link_token_exp > NOW()
	`

	queryCtx, cancel := context.WithTimeout(ctx, db.MEDIUM_QUERY_TIMEOUT)
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

	return &magicLinkTokenPayload, db.ParseDbError(err)
}

// ----- DELETE -----
func (store *PostgresVTokenStore) Delete(ctx context.Context, verificationId int64) error {
	query := `
	DELETE from verification_tokens WHERE id = $1
	`

	queryCtx, cancel := context.WithTimeout(ctx, db.MEDIUM_QUERY_TIMEOUT)
	defer cancel()

	return db.HandleExecContextResult(store.db.ExecContext(
		queryCtx,
		query,
		verificationId,
	))
}
