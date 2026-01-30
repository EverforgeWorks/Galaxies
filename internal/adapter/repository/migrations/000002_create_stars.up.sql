BEGIN;

CREATE TABLE IF NOT EXISTS stars (
    id UUID PRIMARY KEY,
    name VARCHAR(255) NOT NULL,
    x INT NOT NULL,
    y INT NOT NULL,
    created_at TIMESTAMP DEFAULT NOW(),
    updated_at TIMESTAMP DEFAULT NOW()
);

-- Ensure no two stars can occupy the same coordinate
CREATE UNIQUE INDEX IF NOT EXISTS idx_stars_coordinates ON stars(x, y);

COMMIT;
