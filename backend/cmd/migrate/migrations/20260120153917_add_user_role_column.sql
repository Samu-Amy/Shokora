-- +goose Up
-- +goose StatementBegin
ALTER TABLE users
ADD COLUMN user_role smallint NOT NULL DEFAULT 0 CHECK (user_role BETWEEN 0 AND 3);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
ALTER TABLE users DROP COLUMN user_role;
-- +goose StatementEnd
