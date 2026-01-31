BEGIN;

ALTER TABLE user_stage_progress
DROP COLUMN IF EXISTS last_started_at;  

COMMIT;