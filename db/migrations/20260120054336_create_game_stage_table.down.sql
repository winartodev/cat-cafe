BEGIN;

DROP TABLE IF EXISTS game_stages CASCADE;

DROP TABLE IF EXISTS stage_customer_configs CASCADE;

DROP TABLE IF EXISTS stage_staff_configs CASCADE;

DROP TABLE IF EXISTS stage_kitchen_configs CASCADE;

DROP TABLE IF EXISTS stage_camera_configs CASCADE;

DROP TABLE IF EXISTS kitchen_phase_completion_rewards CASCADE;

DROP TABLE IF EXISTS user_stage_progress CASCADE;

DROP TABLE IF EXISTS user_kitchen_phase_progress CASCADE;

DROP TABLE IF EXISTS user_kitchen_phase_reward_claims CASCADE;

COMMIT;
