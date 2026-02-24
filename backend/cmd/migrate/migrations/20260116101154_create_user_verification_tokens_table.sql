-- +goose Up
-- +goose StatementBegin
CREATE TABLE IF NOT EXISTS verification_tokens(
  id bigint GENERATED ALWAYS AS IDENTITY PRIMARY KEY,
  
  user_id bigint NOT NULL REFERENCES users(id) ON DELETE CASCADE,
  verification_type smallint NOT NULL CHECK (verification_type BETWEEN 0 AND 2),
  
  magic_link_token_hash bytea UNIQUE,
  magic_link_token_exp timestamp(0) with time zone,
  
  otp_hash bytea NOT NULL,
  otp_exp timestamp(0) with time zone NOT NULL,
  otp_attempts smallint NOT NULL DEFAULT 0 CHECK (otp_attempts BETWEEN 0 AND 255),
  
  created_at timestamp(0) with time zone NOT NULL DEFAULT NOW(),
  updated_at timestamp(0) with time zone NOT NULL DEFAULT NOW(),
  
  UNIQUE (user_id, verification_type),

  CONSTRAINT magic_link_consistency
  CHECK (
    (magic_link_token_hash IS NULL AND magic_link_token_exp IS NULL)
    OR
    (magic_link_token_hash IS NOT NULL AND magic_link_token_exp IS NOT NULL)
  )
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
