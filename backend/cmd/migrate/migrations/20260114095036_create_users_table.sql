-- +goose Up
-- +goose StatementBegin
CREATE TABLE IF NOT EXISTS users(
  id bigint GENERATED ALWAYS AS IDENTITY PRIMARY KEY,
  first_name varchar(125) NOT NULL,
  last_name varchar(125) NOT NULL,
  email citext UNIQUE NOT NULL,
  password bytea NOT NULL,
  image_url text,
  birth_date date,
  is_verified boolean NOT NULL DEFAULT FALSE,
  is_active boolean NOT NULL DEFAULT TRUE,
  user_role smallint NOT NULL DEFAULT 0 CHECK (user_role BETWEEN 0 AND 3),
  created_at timestamp(0) with time zone NOT NULL DEFAULT NOW(),
  updated_at timestamp(0) with time zone NOT NULL DEFAULT NOW()
);

CREATE TRIGGER update_users_updated_at
BEFORE UPDATE ON users
FOR EACH ROW
EXECUTE FUNCTION set_updated_at();
-- +goose StatementEnd


-- +goose Down
-- +goose StatementBegin
DROP TRIGGER IF EXISTS update_users_updated_at ON users;
DROP TABLE IF EXISTS users;
-- +goose StatementEnd
