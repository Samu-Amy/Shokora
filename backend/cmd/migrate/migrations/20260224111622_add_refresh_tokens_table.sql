-- +goose Up
-- +goose StatementBegin
CREATE TABLE IF NOT EXISTS refresh_tokens(
  id bigint GENERATED ALWAYS AS IDENTITY PRIMARY KEY,

  session_id bigint NOT NULL REFERENCES user_sessions(id) ON DELETE CASCADE,
  
  token_hash bytea NOT NULL UNIQUE,
  expires_at timestamp(0) with time zone NOT NULL,
  
  replaces bigint REFERENCES refresh_tokens(id) ON DELETE RESTRICT,
  revoked_at timestamp(0) with time zone,

  created_at timestamp(0) with time zone NOT NULL DEFAULT NOW()
);

-- Avoid replaces same id multiple times
CREATE UNIQUE INDEX refresh_tokens_replaces_unique
ON refresh_tokens(replaces)
WHERE replaces IS NOT NULL;

CREATE INDEX idx_refresh_tokens_session_id
ON refresh_tokens(session_id);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS refresh_tokens;
-- +goose StatementEnd
