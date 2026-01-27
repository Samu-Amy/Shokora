-- +goose Up
-- +goose StatementBegin
CREATE TABLE IF NOT EXISTS verification_tokens(
  user_id bigint NOT NULL REFERENCES users(id) ON DELETE CASCADE,
  verification_type smallint NOT NULL CHECK (verification_type BETWEEN 0 AND 2),
  
  magic_link_token bytea UNIQUE NOT NULL,
  magic_link_token_exp timestamp(0) with time zone NOT NULL,
  
  otp bytea NOT NULL,
  otp_exp timestamp(0) with time zone NOT NULL,
  otp_attempts smallint NOT NULL DEFAULT 0,
  
  created_at timestamp(0) with time zone NOT NULL DEFAULT NOW(),
  updated_at timestamp(0) with time zone NOT NULL DEFAULT NOW(),
  
  PRIMARY KEY (user_id, verification_type)
  -- UNIQUE (user_id, otp_hash)
);

CREATE TRIGGER update_verification_tokens_updated_at
BEFORE UPDATE ON verification_tokens
FOR EACH ROW
EXECUTE FUNCTION set_updated_at();
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS verification_tokens;
-- +goose StatementEnd
