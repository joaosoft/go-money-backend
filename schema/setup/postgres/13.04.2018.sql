CREATE OR REPLACE FUNCTION function_updated_at()
RETURNS TRIGGER AS $$
BEGIN
   NEW.updated_at = now();
   RETURN NEW;
END;
$$ language 'plpgsql';

CREATE TABLE users (
  user_id                 UUID NOT NULL,
  name                    TEXT NOT NULL,
  password                TEXT NOT NULL,
  description             TEXT NOT NULL,
  created_at              TIMESTAMP DEFAULT NOW(),
  updated_at              TIMESTAMP DEFAULT NOW(),
  PRIMARY KEY(user_id)
);

CREATE TRIGGER trigger_users_updated_at BEFORE UPDATE
  ON users FOR EACH ROW EXECUTE PROCEDURE function_updated_at();
