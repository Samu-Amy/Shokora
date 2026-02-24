-- +goose Up
-- +goose StatementBegin
CREATE TABLE IF NOT EXISTS user_sessions(
  id bigint GENERATED ALWAYS AS IDENTITY PRIMARY KEY,

  user_id bigint NOT NULL REFERENCES users(id) ON DELETE CASCADE,
  expires_at timestamp(0) with time zone NOT NULL,

  created_at timestamp(0) with time zone NOT NULL DEFAULT NOW()
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS user_sessions;
-- +goose StatementEnd
