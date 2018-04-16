CREATE SCHEMA money;


CREATE OR REPLACE FUNCTION money.function_updated_at()
  RETURNS TRIGGER AS $$
  BEGIN
   NEW.updated_at = now();
   RETURN NEW;
  END;
  $$ LANGUAGE 'plpgsql';


-- USERS
CREATE TABLE money.users (
  user_id                 UUID NOT NULL,
  name                    TEXT NOT NULL,
  email                   TEXT NOT NULL UNIQUE,
  password                TEXT NOT NULL,
  description             TEXT NOT NULL,
  created_at              TIMESTAMP DEFAULT NOW(),
  updated_at              TIMESTAMP DEFAULT NOW(),
  PRIMARY KEY(user_id)
);

CREATE TRIGGER money.trigger_users_updated_at BEFORE UPDATE
  ON money.users FOR EACH ROW EXECUTE PROCEDURE money.function_updated_at();


-- TRANSACTIONS
CREATE TABLE money.transactions (
  transaction_id          UUID NOT NULL,
  user_id                 UUID NOT NULL,
  category_id             UUID NOT NULL,
  price                   FLOAT NOT NULL,
  email                   TEXT NOT NULL UNIQUE,
  description             TEXT NOT NULL,
  date                    TIMESTAMP,
  created_at              TIMESTAMP DEFAULT NOW(),
  updated_at              TIMESTAMP DEFAULT NOW(),
  PRIMARY KEY(user_id)
);

CREATE TRIGGER money.trigger_transactions_updated_at BEFORE UPDATE
  ON money.transactions FOR EACH ROW EXECUTE PROCEDURE money.function_updated_at();