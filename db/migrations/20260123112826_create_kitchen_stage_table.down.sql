BEGIN;


DROP TABLE IF EXISTS user_kitchen_phase_reward_claims CASCADE;

DROP TABLE IF EXISTS user_kitchen_phase_progress CASCADE;

DROP TABLE IF EXISTS user_kitchen_progress CASCADE;

DROP TABLE IF EXISTS kitchen_phase_completion_rewards CASCADE;

DROP TABLE IF EXISTS stage_kitchen_configs CASCADE;

DROP TABLE IF EXISTS kitchen_stations CASCADE;

DROP TABLE IF EXISTS food_level_overrides CASCADE;

DROP TABLE IF EXISTS food_items CASCADE;

COMMIT;
