BEGIN;

-- Drop the index first, then the table
DROP INDEX IF EXISTS idx_players_external_id;
DROP TABLE IF EXISTS players;

COMMIT;
