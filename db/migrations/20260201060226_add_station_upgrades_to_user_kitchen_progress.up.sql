BEGIN;

ALTER TABLE user_kitchen_progress
    ADD COLUMN IF NOT EXISTS station_upgrades JSONB DEFAULT '[]'::jsonb NOT NULL;

COMMIT;
