-- +goose Up
-- +goose StatementBegin
CREATE TABLE IF NOT EXISTS refresh_tokens(
  id bigint GENERATED ALWAYS AS IDENTITY PRIMARY KEY,

  user_id bigint NOT NULL REFERENCES users(id) ON DELETE CASCADE,
  session_id uuid NOT NULL,
  
  token_hash bytea NOT NULL UNIQUE,
  expires_at timestamp(0) with time zone NOT NULL,
  
  replaces bigint REFERENCES refresh_tokens(id),
  revoked_at timestamp(0) with time zone,

  created_at timestamp(0) with time zone NOT NULL DEFAULT NOW()
);

CREATE INDEX idx_refresh_token_user_session
ON refresh_tokens(user_id, session_id);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS refresh_tokens;
-- +goose StatementEnd
