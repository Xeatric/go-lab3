DROP TRIGGER IF EXISTS update_tiles_updated_at ON tiles;
DROP FUNCTION IF EXISTS update_updated_at_column();
DROP INDEX IF EXISTS idx_tiles_deleted_at;
DROP INDEX IF EXISTS idx_tiles_shape;
DROP INDEX IF EXISTS idx_tiles_color;
DROP INDEX IF EXISTS idx_tiles_price;
DROP TABLE IF EXISTS tiles;