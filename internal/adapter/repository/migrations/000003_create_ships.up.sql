BEGIN;

CREATE TABLE IF NOT EXISTS ships (
    id UUID PRIMARY KEY,
    player_id UUID NOT NULL REFERENCES players(id) ON DELETE CASCADE,
    
    -- Flavor & Classification
    name VARCHAR(255) NOT NULL,
    type VARCHAR(50) NOT NULL, -- e.g. "SCOUT"
    
    -- Fuel System
    fuel INT NOT NULL,
    max_fuel INT NOT NULL,
    
    -- Capabilities
    fuel_efficiency DOUBLE PRECISION NOT NULL, -- Maps to float64
    max_range DOUBLE PRECISION NOT NULL,       -- Maps to float64
    
    -- Metadata
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT NOW()
);

-- Index for rapid retrieval of a player's active ship during login/session creation
CREATE INDEX IF NOT EXISTS idx_ships_player_id ON ships(player_id);

COMMIT;