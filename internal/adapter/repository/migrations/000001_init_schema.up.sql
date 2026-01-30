BEGIN;

-- Create the players table to store core pilot data
CREATE TABLE IF NOT EXISTS players (
    id UUID PRIMARY KEY,                   -- Internal Game ID
    external_id VARCHAR(255) UNIQUE NOT NULL, -- OAuth ID (e.g., "github_12345")
    name VARCHAR(50) NOT NULL,             -- Display Name
    credits INT DEFAULT 1000,              -- Starting wallet
    is_admin BOOLEAN DEFAULT FALSE,        -- Admin flag
    last_login TIMESTAMP DEFAULT NOW(),    -- Login tracking
    current_system_id UUID                 -- Location (to be linked to universe)
);

-- Create an index for external_id. 
-- Since we look this up every time someone logs in, an index makes it near-instant.
CREATE INDEX IF NOT EXISTS idx_players_external_id ON players(external_id);

COMMIT;
