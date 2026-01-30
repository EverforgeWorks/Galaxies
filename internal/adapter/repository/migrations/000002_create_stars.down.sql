BEGIN;

-- Drop the index first (optional, as dropping table usually drops its indexes)
DROP INDEX IF EXISTS idx_stars_coordinates;

-- Drop the table
DROP TABLE IF EXISTS stars;

COMMIT;
