-- +goose Up
-- +goose StatementBegin
CREATE TABLE IF NOT EXISTS email_verification_tokens(
  token bytea PRIMARY KEY,
  user_id bigint NOT NULL,
  expiry TIMESTAMP(0) with time zone NOT NULL
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS email_verification_tokens;
-- +goose StatementEnd
