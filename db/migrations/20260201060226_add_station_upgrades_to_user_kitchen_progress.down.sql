BEGIN;

ALTER TABLE user_kitchen_progress
    DROP COLUMN IF EXISTS station_upgrades;

COMMIT;
