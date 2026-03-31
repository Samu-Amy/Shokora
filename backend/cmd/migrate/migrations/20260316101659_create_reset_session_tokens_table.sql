-- +goose Up
-- +goose StatementBegin
CREATE TABLE IF NOT EXISTS reset_session_tokens(
  token_hash bytea NOT NULL PRIMARY KEY,
  user_id bigint NOT NULL REFERENCES users(id) ON DELETE CASCADE,
  expires_at timestamp(0) with time zone NOT NULL,
  
  created_at timestamp(0) with time zone NOT NULL DEFAULT NOW()
);

CREATE UNIQUE INDEX idx_reset_session_tokens_user_id
ON reset_session_tokens(user_id);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS reset_session_tokens;
-- +goose StatementEnd
