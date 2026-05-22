DROP TRIGGER IF EXISTS update_tokens_updated_at ON tokens;
DROP TRIGGER IF EXISTS update_users_updated_at ON users;
DROP FUNCTION IF EXISTS update_updated_at_column();

DROP INDEX IF EXISTS idx_tiles_user_id;
ALTER TABLE tiles DROP COLUMN IF EXISTS user_id;

DROP TABLE IF EXISTS tokens;
DROP TABLE IF EXISTS users;