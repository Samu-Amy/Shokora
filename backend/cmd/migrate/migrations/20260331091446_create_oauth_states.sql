-- +goose Up
-- +goose StatementBegin
CREATE TABLE IF NOT EXISTS oauth_states (
  state      VARCHAR(64) PRIMARY KEY,
  created_at timestamp(0) with time zone NOT NULL DEFAULT NOW()
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS oauth_states;
-- +goose StatementEnd
