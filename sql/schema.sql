-- 1. ENABLE UUID EXTENSION
CREATE EXTENSION IF NOT EXISTS "uuid-ossp";

-- 2. SYSTEMS
CREATE TABLE systems (
    id UUID PRIMARY KEY,
    name TEXT NOT NULL,
    x INT NOT NULL,
    y INT NOT NULL,
    political TEXT,
    economic TEXT,
    social TEXT,
    stats JSONB NOT NULL DEFAULT '{}'::jsonb
);
CREATE INDEX idx_systems_coords ON systems (x, y);

-- 3. PLAYERS
CREATE TABLE players (
    id UUID PRIMARY KEY,
    name TEXT NOT NULL UNIQUE,
    external_id TEXT UNIQUE,  -- <--- ADDED THIS
    credits INT NOT NULL DEFAULT 0,
    current_system_id UUID REFERENCES systems(id),
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    last_login TIMESTAMP WITH TIME ZONE DEFAULT NOW()
);

-- 4. SHIPS
CREATE TABLE ships (
    id UUID PRIMARY KEY,
    player_id UUID REFERENCES players(id) UNIQUE NOT NULL,
    name TEXT,
    model_name TEXT,
    current_hull FLOAT NOT NULL,
    current_shield FLOAT NOT NULL,
    current_fuel FLOAT NOT NULL,
    stats JSONB NOT NULL DEFAULT '{}'::jsonb,
    cargo JSONB NOT NULL DEFAULT '[]'::jsonb,
    passengers JSONB NOT NULL DEFAULT '[]'::jsonb,
    crew JSONB NOT NULL DEFAULT '[]'::jsonb
);

-- 5. SHIP INVENTORY
CREATE TABLE IF NOT EXISTS ship_inventory (
    id UUID PRIMARY KEY,
    ship_id UUID REFERENCES ships(id) ON DELETE CASCADE,
    name TEXT NOT NULL,
    category TEXT NOT NULL,
    base_value INT NOT NULL,
    rarity INT NOT NULL,
    is_illegal BOOLEAN NOT NULL,
    
    quantity INT NOT NULL,
    avg_cost FLOAT NOT NULL DEFAULT 0.0 -- Tracks weighted average cost of acquisition
);

CREATE INDEX idx_inventory_ship ON ship_inventory(ship_id);

CREATE TABLE IF NOT EXISTS market_listings (
    id UUID PRIMARY KEY,
    system_id UUID REFERENCES systems(id) ON DELETE CASCADE,
    name TEXT NOT NULL,
    category TEXT NOT NULL,
    base_value INT NOT NULL,
    rarity INT NOT NULL,
    is_illegal BOOLEAN NOT NULL,
    quantity INT NOT NULL
);

CREATE INDEX idx_market_system ON market_listings(system_id);

-- Add this to your schema or run it manually if you don't want to re-seed yet
ALTER TABLE players ADD COLUMN IF NOT EXISTS is_admin BOOLEAN NOT NULL DEFAULT FALSE;