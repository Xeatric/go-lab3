CREATE TABLE IF NOT EXISTS tiles (
    id SERIAL PRIMARY KEY,
    name VARCHAR(200) NOT NULL,
    shape VARCHAR(50) NOT NULL,
    color VARCHAR(50) NOT NULL,
    size VARCHAR(50) NOT NULL,
    material VARCHAR(100),
    price_per_m2 NUMERIC(10,2) NOT NULL,
    stock INTEGER DEFAULT 0,
    description TEXT,
    image_url VARCHAR(500),
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    deleted_at TIMESTAMP WITH TIME ZONE,
    
    CONSTRAINT valid_shape CHECK (shape IN ('square', 'rectangle', 'hexagon', 'circle')),
    CONSTRAINT positive_price CHECK (price_per_m2 >= 0),
    CONSTRAINT non_negative_stock CHECK (stock >= 0)
);

CREATE INDEX idx_tiles_deleted_at ON tiles(deleted_at);
CREATE INDEX idx_tiles_shape ON tiles(shape);
CREATE INDEX idx_tiles_color ON tiles(color);
CREATE INDEX idx_tiles_price ON tiles(price_per_m2);

-- Триггер для автоматического обновления updated_at
CREATE OR REPLACE FUNCTION update_updated_at_column()
RETURNS TRIGGER AS $$
BEGIN
    NEW.updated_at = CURRENT_TIMESTAMP;
    RETURN NEW;
END;
$$ language 'plpgsql';

CREATE TRIGGER update_tiles_updated_at 
    BEFORE UPDATE ON tiles 
    FOR EACH ROW 
    EXECUTE FUNCTION update_updated_at_column();