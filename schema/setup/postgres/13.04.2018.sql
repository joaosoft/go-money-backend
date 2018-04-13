CREATE SCHEMA money;


CREATE OR REPLACE FUNCTION money.function_updated_at()
  RETURNS TRIGGER AS $$
  BEGIN
   NEW.updated_at = now();
   RETURN NEW;
  END;
  $$ LANGUAGE 'plpgsql';


CREATE TABLE money.users (
  user_id                 UUID NOT NULL,
  name                    TEXT NOT NULL,
  email                    TEXT NOT NULL,
  password                TEXT NOT NULL,
  description             TEXT NOT NULL,
  created_at              TIMESTAMP DEFAULT NOW(),
  updated_at              TIMESTAMP DEFAULT NOW(),
  PRIMARY KEY(user_id)
);


CREATE TRIGGER money.trigger_users_updated_at BEFORE UPDATE
  ON money.users FOR EACH ROW EXECUTE PROCEDURE money.function_updated_at();