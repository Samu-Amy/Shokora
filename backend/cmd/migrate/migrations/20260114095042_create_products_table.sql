-- TODO: rifai questa migartion dopo tutte quelle legate a user e auth

-- +goose Up
-- +goose StatementBegin
CREATE TABLE IF NOT EXISTS products(
  id int GENERATED ALWAYS AS IDENTITY PRIMARY KEY,

  name varchar(150) NOT NULL, -- TODO: ottimizzare name e description per text search (text GIN (to_tsvector('italian', name) (?))
  description text NOT NULL,
  image_url text NOT NULL,
  price numeric(10, 2) NOT NULL,
  discount numeric(4, 3) NOT NULL DEFAULT 0,

  version INT NOT NULL DEFAULT 0,

  created_at timestamp(0) with time zone NOT NULL DEFAULT NOW(),
  updated_at timestamp(0) with time zone NOT NULL DEFAULT NOW(),
 
  
  -- CONSTRAINTS --

  CONSTRAINT products_price_range_check CHECK (price > 0),

  CONSTRAINT products_discount_range_check CHECK (discount BETWEEN 0 AND 1)
);

CREATE TRIGGER update_products_updated_at
BEFORE UPDATE ON products
FOR EACH ROW
EXECUTE FUNCTION set_updated_at();

-- TODO: aggiungi index (es. su nome e descrizione)?
-- +goose StatementEnd


-- +goose Down
-- +goose StatementBegin
DROP TRIGGER IF EXISTS update_products_updated_at ON products;
DROP TABLE IF EXISTS products;
-- +goose StatementEnd
