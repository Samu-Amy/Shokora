CREATE TABLE IF NOT EXISTS users(
  id bigint GENERATED ALWAYS AS IDENTITY PRIMARY KEY,
  first_name varchar(125) NOT NULL,
  last_name varchar(125) NOT NULL,
  email citext UNIQUE NOT NULL, -- TODO: aggiungi index (anche se in teoria con unique lo aggiunge già) (?)
  password bytea NOT NULL, -- TODO: usare TEXT (?)
  created_at timestamp(0) with time zone NOT NULL DEFAULT NOW(),
  updated_at timestamp(0) with time zone NOT NULL DEFAULT NOW(),
);

-- CREATE INDEX IF NOT EXISTS idx_users_email ON users(email); -- TODO: aggiungere per index (?)
