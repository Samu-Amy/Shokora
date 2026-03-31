-- +goose Up
-- +goose StatementBegin
CREATE TABLE IF NOT EXISTS user_settings(
  id bigint GENERATED ALWAYS AS IDENTITY PRIMARY KEY,

  user_id bigint NOT NULL REFERENCES users(id) ON DELETE CASCADE,

  two_factor_auth boolean NOT NULL DEFAULT FALSE,
  notifications smallint NOT NULL DEFAULT 0, -- TODO: aggiorna (metti notifiche di base - magari nel codice go durante la creazione)

  created_at timestamp(0) with time zone NOT NULL DEFAULT NOW(),
  updated_at timestamp(0) with time zone NOT NULL DEFAULT NOW(),

  -- CONSTRAINTS --

  CONSTRAINT user_settings_user_id_unique UNIQUE (user_id)
);

-- Updated At
CREATE TRIGGER update_user_settings_updated_at
BEFORE UPDATE ON user_settings
FOR EACH ROW
EXECUTE FUNCTION set_updated_at();
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS user_settings;
-- +goose StatementEnd
