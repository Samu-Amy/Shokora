-- +goose Up
-- +goose StatementBegin
CREATE TABLE IF NOT EXISTS users(
  id bigint GENERATED ALWAYS AS IDENTITY PRIMARY KEY,
  first_name varchar(125) NOT NULL,
  last_name varchar(125) NOT NULL,
  email citext UNIQUE NOT NULL,
  password bytea NOT NULL,
  -- is_verified BOOLEAN NOT NULL DEFAULT FALSE; -- TODO: attiva questo campo (per ora migration alter)
  created_at timestamp(0) with time zone NOT NULL DEFAULT NOW(),
  updated_at timestamp(0) with time zone NOT NULL DEFAULT NOW()
);

CREATE TRIGGER update_users_updated_at
BEFORE UPDATE ON users
FOR EACH ROW
EXECUTE FUNCTION set_updated_at();
-- +goose StatementEnd

-- TODO: (email) - aggiungi index (con gin extension?) (anche se in teoria con unique lo aggiunge già) (?)
-- TODO: (password) - usare TEXT (?)
-- CREATE INDEX IF NOT EXISTS idx_users_id ON users(id);
-- CREATE INDEX IF NOT EXISTS idx_users_email ON users(email); -- TODO: aggiungere per index (?)
-- TODO: per foreign keys usare "FOREIGN KEY (<field>) REFERENCES <table> (<field>) ON DELETE <delete_option>" -- attenzione all'opzione di delete


-- +goose Down
-- +goose StatementBegin
DROP TRIGGER IF EXISTS update_users_updated_at ON users;
DROP TABLE IF EXISTS users;
-- +goose StatementEnd
