CREATE SCHEMA money;

-- GLOBAL
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
  password_hash           TEXT NOT NULL,
  description             TEXT,
  created_at              TIMESTAMP DEFAULT NOW(),
  updated_at              TIMESTAMP DEFAULT NOW(),
  PRIMARY KEY(user_id)
);

CREATE TRIGGER trigger_users_updated_at BEFORE UPDATE
  ON money.users FOR EACH ROW EXECUTE PROCEDURE money.function_updated_at();


-- WALLETS
CREATE TABLE money.wallets (
  wallet_id               UUID NOT NULL,
  user_id                 TEXT NOT NULL,
  name                    TEXT NOT NULL,
  description             TEXT,
  password                TEXT,
  created_at              TIMESTAMP DEFAULT NOW(),
  updated_at              TIMESTAMP DEFAULT NOW(),
  PRIMARY KEY(wallet_id)
);

CREATE TRIGGER trigger_wallers_updated_at BEFORE UPDATE
  ON money.wallets FOR EACH ROW EXECUTE PROCEDURE money.function_updated_at();


-- IMAGES
CREATE TABLE money.images (
  image_id                UUID NOT NULL,
  user_id                 UUID NOT NULL,
  name                    TEXT NOT NULL,
  description             TEXT,
  url                     TEXT,
  file_name               TEXT,
  format                  TEXT,
  raw_image               BYTEA,
  created_at              TIMESTAMP DEFAULT NOW(),
  updated_at              TIMESTAMP DEFAULT NOW(),
  FOREIGN KEY(user_id) REFERENCES money.users(user_id),
  PRIMARY KEY(image_id)
);

CREATE TRIGGER trigger_images_updated_at BEFORE UPDATE
  ON money.images FOR EACH ROW EXECUTE PROCEDURE money.function_updated_at();


-- CATEGORIES
CREATE TABLE money.categories (
  category_id             UUID NOT NULL,
  user_id                 UUID NOT NULL,
  image_id                UUID NOT NULL,
  name                    TEXT NOT NULL,
  description             TEXT,
  created_at              TIMESTAMP DEFAULT NOW(),
  updated_at              TIMESTAMP DEFAULT NOW(),
  FOREIGN KEY(user_id) REFERENCES money.users(user_id),
  FOREIGN KEY(image_id) REFERENCES money.images(image_id),
  PRIMARY KEY(category_id)
);

CREATE TRIGGER trigger_categories_updated_at BEFORE UPDATE
  ON money.categories FOR EACH ROW EXECUTE PROCEDURE money.function_updated_at();


-- TRANSACTIONS
CREATE TABLE money.transactions (
  transaction_id          UUID NOT NULL,
  user_id                 UUID NOT NULL,
  wallet_id               UUID NOT NULL,
  category_id             UUID NOT NULL,
  price                   FLOAT NOT NULL,
  description             TEXT,
  date                    TIMESTAMP,
  created_at              TIMESTAMP DEFAULT NOW(),
  updated_at              TIMESTAMP DEFAULT NOW(),
  FOREIGN KEY(user_id) REFERENCES money.users(user_id),
  FOREIGN KEY(wallet_id) REFERENCES money.wallets(wallet_id),
  FOREIGN KEY(category_id) REFERENCES money.categories(category_id),
  PRIMARY KEY(transaction_id)
);

CREATE TRIGGER trigger_transactions_updated_at BEFORE UPDATE
  ON money.transactions FOR EACH ROW EXECUTE PROCEDURE money.function_updated_at();