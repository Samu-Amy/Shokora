-- +goose Up
-- +goose StatementBegin
CREATE TABLE IF NOT EXISTS products(
  id bigint GENERATED ALWAYS AS IDENTITY PRIMARY KEY,
  name varchar(150) NOT NULL,
  description text NOT NULL,
  image_url text NOT NULL,
  price numeric(10, 2) NOT NULL,
  discount numeric(4, 3) NOT NULL DEFAULT 0 CHECK (discount >= 0 AND discount <= 1),
  created_at timestamp(0) with time zone NOT NULL DEFAULT NOW(),
  updated_at timestamp(0) with time zone NOT NULL DEFAULT NOW()
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS products;
-- +goose StatementEnd
