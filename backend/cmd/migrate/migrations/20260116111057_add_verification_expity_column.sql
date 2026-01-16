-- TODO: attiva il campo nella creazione della tabella ed elimina questo file

-- +goose Up
-- +goose StatementBegin
ALTER TABLE user_verification_tokens
ADD COLUMN expiry TIMESTAMP(0) with time zone NOT NULL;
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
ALTER TABLE user_verification_tokens DROP COLUMN expiry;
-- +goose StatementEnd
