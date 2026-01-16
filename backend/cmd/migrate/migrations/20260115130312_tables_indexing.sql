-- TODO: eliminare o spostare come ultimo questo file
-- TODO: mettere l'indexing insieme alle tabelle (?)

-- +goose Up
-- +goose StatementBegin
-- CREATE INDEX IF NOT EXISTS idx_users_id ON users(id);
-- CREATE INDEX IF NOT EXISTS idx_users_email ON users(email);

-- CREATE INDEX IF NOT EXISTS idx_products_id ON products(id);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
-- DROP INDEX IF EXISTS idx_products_id;

-- DROP INDEX IF EXISTS idx_users_email;
-- DROP INDEX IF EXISTS idx_users_id;
-- +goose StatementEnd
