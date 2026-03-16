-- +goose Up
-- +goose StatementBegin
CREATE TABLE IF NOT EXISTS verification_tokens(
  id uuid PRIMARY KEY DEFAULT gen_random_uuid(),
  
  user_id bigint NOT NULL REFERENCES users(id) ON DELETE CASCADE,
  verification_type smallint NOT NULL,
  
  magic_link_token_hash bytea,
  magic_link_token_expires_at timestamp(0) with time zone,
  
  otp_hash bytea NOT NULL,
  otp_expires_at timestamp(0) with time zone NOT NULL,
  otp_attempts smallint NOT NULL DEFAULT 0,
  
  created_at timestamp(0) with time zone NOT NULL DEFAULT NOW(),
  updated_at timestamp(0) with time zone NOT NULL DEFAULT NOW(),
   
  
  -- CONSTRAINTS --

  CONSTRAINT v_tokens_user_id_and_verification_type_unique UNIQUE (user_id, verification_type),

  CONSTRAINT v_tokens_verification_type_range_check CHECK (verification_type BETWEEN 0 AND 2),
  
  CONSTRAINT v_tokens_otp_attempts_range_check CHECK (otp_attempts BETWEEN 0 AND 255),

  CONSTRAINT v_tokens_magic_link_token_hash_check CHECK (
    (magic_link_token_hash IS NULL AND magic_link_token_expires_at IS NULL)
    OR
    (magic_link_token_hash IS NOT NULL AND magic_link_token_expires_at IS NOT NULL)
  )
);


-- Avoid replaces same id multiple times with partial unique
CREATE UNIQUE INDEX v_tokens_magic_link_token_hash_unique
ON verification_tokens(magic_link_token_hash, verification_type)
WHERE magic_link_token_hash IS NOT NULL;

-- Updated At
CREATE TRIGGER update_verification_tokens_updated_at
BEFORE UPDATE ON verification_tokens
FOR EACH ROW
-- WHEN (OLD.* IS DISTINCT FROM NEW.*) -- TODO: usare questa riga (anche su altre tabelle)?
EXECUTE FUNCTION set_updated_at();
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS verification_tokens;
-- +goose StatementEnd
