BEGIN;

CREATE TABLE IF NOT EXISTS game_stages (
    id BIGSERIAL PRIMARY KEY,
    slug VARCHAR(50) UNIQUE NOT NULL,
    name VARCHAR(100) NOT NULL,
    starting_coin BIGINT DEFAULT 0 NOT NULL,
    stage_prize BIGINT DEFAULT 0 NOT NULL,
    is_active BOOL DEFAULT true NOT NULL,
    sequence INT DEFAULT 0 NOT NULL,
    created_at TIMESTAMPTZ DEFAULT NOW() NOT NULL,
    updated_at TIMESTAMPTZ DEFAULT NOW() NOT NULL
);

CREATE TABLE IF NOT EXISTS stage_customer_configs (
    id BIGSERIAL PRIMARY KEY,
    stage_id BIGINT NOT NULL REFERENCES game_stages(id) ON DELETE CASCADE,
    customer_spawn_time NUMERIC DEFAULT 0.0 NOT NULL,
    max_customer_order_count INT DEFAULT 0 NOT NULL,
    max_customer_order_variant INT DEFAULT  0 NOT NULL,
    starting_order_table_count INT DEFAULT 0 NOT NULL,
    created_at TIMESTAMPTZ DEFAULT NOW() NOT NULL,
    updated_at TIMESTAMPTZ DEFAULT NOW() NOT NULL,

    UNIQUE(stage_id)
);

CREATE TABLE IF NOT EXISTS stage_staff_configs (
    id BIGSERIAL PRIMARY KEY ,
    stage_id BIGINT NOT NULL REFERENCES game_stages(id) ON DELETE CASCADE,
    starting_staff_manager VARCHAR(50) DEFAULT '' NOT NULL,
    starting_staff_helper VARCHAR(50) DEFAULT '' NOT NULL,
    created_at TIMESTAMPTZ DEFAULT NOW() NOT NULL,
    updated_at TIMESTAMPTZ DEFAULT NOW() NOT NULL,

    UNIQUE(stage_id)
);

CREATE TABLE IF NOT EXISTS stage_kitchen_configs (
    id BIGSERIAL PRIMARY KEY,
    stage_id BIGINT NOT NULL REFERENCES game_stages(id) ON DELETE CASCADE,
    max_level INT DEFAULT 0 NOT NULL,
    upgrade_profit_multiply NUMERIC DEFAULT 0 NOT NULL,
    upgrade_cost_multiply NUMERIC DEFAULT  0 NOT NULL,
    transition_phase_levels INT[] DEFAULT '{}' NOT NULL,
    phase_profit_multipliers NUMERIC[] DEFAULT '{}' NOT NULL,
    phase_upgrade_cost_multipliers NUMERIC[] DEFAULT '{}' NOT NULL,
    table_count_per_phases INT[] DEFAULT '{}' NOT NULL,

    created_at TIMESTAMPTZ DEFAULT NOW() NOT NULL,
    updated_at TIMESTAMPTZ DEFAULT NOW() NOT NULL,

    UNIQUE (stage_id)
);

CREATE TABLE IF NOT EXISTS stage_camera_configs (
   id BIGSERIAL PRIMARY KEY,
   stage_id BIGINT NOT NULL REFERENCES game_stages(id) ON DELETE CASCADE,
   zoom_size NUMERIC DEFAULT 0.0 NOT NULL,
   min_bound_x NUMERIC DEFAULT 0.0 NOT NULL,
   min_bound_y NUMERIC DEFAULT 0.0 NOT NULL,
   max_bound_x NUMERIC DEFAULT 0.0 NOT NULL,
   max_bound_y NUMERIC DEFAULT 0.0 NOT NULL,
   created_at TIMESTAMPTZ DEFAULT NOW() NOT NULL,
   updated_at TIMESTAMPTZ DEFAULT NOW() NOT NULL,

   UNIQUE(stage_id)
);

CREATE TABLE IF NOT EXISTS kitchen_phase_completion_rewards (
    id BIGSERIAL PRIMARY KEY,
    kitchen_config_id BIGINT NOT NULL REFERENCES stage_kitchen_configs ON DELETE CASCADE,
    phase_number INT DEFAULT 0 NOT NULL,
    reward_id BIGINT NOT NULL REFERENCES rewards(id) ON DELETE CASCADE,

    created_at TIMESTAMPTZ DEFAULT NOW() NOT NULL,
    updated_at TIMESTAMPTZ DEFAULT NOW() NOT NULL,

    UNIQUE (kitchen_config_id, phase_number, reward_id)
);

CREATE TABLE IF NOT EXISTS user_stage_progress (
   id BIGSERIAL PRIMARY KEY,
   user_id BIGINT NOT NULL REFERENCES  users(id) ON DELETE CASCADE,
   stage_id BIGINT NOT NULL REFERENCES  game_stages(id) ON DELETE CASCADE,
   is_complete BOOL DEFAULT false NOT NULL,
   completed_at TIMESTAMPTZ,
   created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
   updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),

   UNIQUE (user_id, stage_id)
);

CREATE TABLE IF NOT EXISTS user_kitchen_phase_progress (
   id BIGSERIAL PRIMARY KEY,
   user_id BIGINT NOT NULL REFERENCES users(id) ON DELETE CASCADE,
   kitchen_config_id BIGINT NOT NULL REFERENCES stage_kitchen_configs(id) ON DELETE CASCADE,
   current_phase INT DEFAULT 1 NOT NULL,
   completed_phases JSONB DEFAULT '[]'::jsonb NOT NULL,
   created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
   updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),

   UNIQUE(user_id, kitchen_config_id)
);

CREATE TABLE IF NOT EXISTS user_kitchen_phase_reward_claims (
    id BIGSERIAL PRIMARY KEY,
    user_id BIGINT NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    kitchen_config_id BIGINT NOT NULL REFERENCES stage_kitchen_configs(id) ON DELETE CASCADE,
    phase_number INT NOT NULL,
    reward_id BIGINT NOT NULL REFERENCES rewards(id) ON DELETE CASCADE,
    claimed_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),

    UNIQUE(user_id, kitchen_config_id, phase_number, reward_id)
);

CREATE UNIQUE INDEX unique_active_sequence_per_stage
    ON game_stages (sequence)
    WHERE is_active = true;

CREATE INDEX idx_stage_customer_config_stage ON stage_customer_configs(stage_id);
CREATE INDEX idx_stage_staff_config_stage ON stage_staff_configs(stage_id);
CREATE INDEX idx_stage_kitchen_config_stage ON stage_kitchen_configs(stage_id);
CREATE INDEX idx_stage_camera_config_stage ON stage_camera_configs(stage_id);
CREATE INDEX idx_kitchen_phase_rewards_config ON kitchen_phase_completion_rewards(kitchen_config_id);
CREATE INDEX idx_kitchen_phase_rewards_reward ON kitchen_phase_completion_rewards(reward_id);

CREATE INDEX idx_user_stage_progress_user ON user_stage_progress(user_id);
CREATE INDEX idx_user_stage_progress_stage ON user_stage_progress(stage_id);
CREATE INDEX idx_user_kitchen_progress_user ON user_kitchen_phase_progress(user_id);
CREATE INDEX idx_user_kitchen_progress_config ON user_kitchen_phase_progress(kitchen_config_id);
CREATE INDEX idx_user_kitchen_claims_user ON user_kitchen_phase_reward_claims(user_id);

CREATE INDEX idx_kitchen_transition_levels ON stage_kitchen_configs USING GIN (transition_phase_levels);
CREATE INDEX idx_kitchen_profit_multipliers ON stage_kitchen_configs USING GIN (phase_profit_multipliers);
CREATE INDEX idx_user_kitchen_completed_phases ON user_kitchen_phase_progress USING GIN (completed_phases);

COMMIT;
