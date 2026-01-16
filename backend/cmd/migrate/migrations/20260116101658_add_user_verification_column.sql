-- TODO: attiva il campo nella creazione della tabella ed elimina questo file

-- +goose Up
-- +goose StatementBegin
ALTER TABLE users
ADD COLUMN is_verified BOOLEAN NOT NULL DEFAULT FALSE;
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
ALTER TABLE users DROP COLUMN is_verified;
-- +goose StatementEnd
