DROP TABLE IF EXISTS money.users;
DROP TRIGGER IF EXISTS trigger_users_updated_at on money.users;

DROP TABLE IF EXISTS money.wallets;
DROP TRIGGER IF EXISTS trigger_wallets_updated_at on money.wallets;

DROP TABLE IF EXISTS money.images;
DROP TRIGGER IF EXISTS trigger_images_updated_at on money.images;

DROP TABLE IF EXISTS money.categories;
DROP TRIGGER IF EXISTS trigger_categories_updated_at on money.categories;

DROP TABLE IF EXISTS money.transactions;
DROP TRIGGER IF EXISTS trigger_transactions_updated_at on money.transactions;

DROP FUNCTION IF EXISTS money.function_updated_at();
DROP SCHEMA money;