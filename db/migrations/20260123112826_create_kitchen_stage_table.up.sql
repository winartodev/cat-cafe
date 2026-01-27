BEGIN;

CREATE TABLE IF NOT EXISTS food_items (
    id BIGSERIAL PRIMARY KEY,
    slug VARCHAR(50) UNIQUE NOT NULL,
    name VARCHAR(100) DEFAULT '' NOT NULL,
    initial_cost BIGINT DEFAULT 0 NOT NULL,
    initial_profit BIGINT DEFAULT 0 NOT NULL,
    cooking_time NUMERIC DEFAULT 0 NOT NULL,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    CONSTRAINT check_initial_cost_non_negative CHECK (initial_cost >= 0),
    CONSTRAINT check_initial_profit_non_negative CHECK (initial_profit >= 0),
    CONSTRAINT check_cooking_time_positive CHECK (cooking_time > 0)
);

CREATE TABLE IF NOT EXISTS food_level_overrides (
    id BIGSERIAL PRIMARY KEY,
    food_item_id BIGINT NOT NULL REFERENCES food_items (id) ON DELETE CASCADE,
    level BIGINT DEFAULT 0 NOT NULL,
    cost BIGINT DEFAULT 0 NOT NULL,
    profit BIGINT DEFAULT 0 NOT NULL,
    preparation_time NUMERIC DEFAULT 0 NOT NULL,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    UNIQUE (food_item_id, level),
    CONSTRAINT check_level_positive CHECK (level > 0)
);

CREATE TABLE IF NOT EXISTS kitchen_stations (
    id BIGSERIAL PRIMARY KEY,
    stage_id BIGINT NOT NULL REFERENCES game_stages (id) ON DELETE CASCADE,
    food_item_id BIGINT NOT NULL REFERENCES food_items (id) ON DELETE CASCADE,
    auto_unlock BOOL DEFAULT false NOT NULL,
    unlock_phase INT DEFAULT 1 NOT NULL,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    UNIQUE (stage_id, food_item_id)
);

CREATE TABLE IF NOT EXISTS stage_kitchen_configs (
    id BIGSERIAL PRIMARY KEY,
    stage_id BIGINT NOT NULL REFERENCES game_stages (id) ON DELETE CASCADE,
    max_level INT DEFAULT 0 NOT NULL,
    upgrade_profit_multiply NUMERIC DEFAULT 0 NOT NULL,
    upgrade_cost_multiply NUMERIC DEFAULT 0 NOT NULL,
    transition_phase_levels INT[] DEFAULT '{}' NOT NULL,
    phase_profit_multipliers NUMERIC[] DEFAULT '{}' NOT NULL,
    phase_upgrade_cost_multipliers NUMERIC[] DEFAULT '{}' NOT NULL,
    table_count_per_phases INT[] DEFAULT '{}' NOT NULL,
    created_at TIMESTAMPTZ DEFAULT NOW() NOT NULL,
    updated_at TIMESTAMPTZ DEFAULT NOW() NOT NULL,
    UNIQUE (stage_id)
);

CREATE TABLE IF NOT EXISTS kitchen_phase_completion_rewards (
    id BIGSERIAL PRIMARY KEY,
    kitchen_config_id BIGINT NOT NULL REFERENCES stage_kitchen_configs ON DELETE CASCADE,
    phase_number INT DEFAULT 0 NOT NULL,
    reward_id BIGINT NOT NULL REFERENCES rewards (id) ON DELETE CASCADE,
    created_at TIMESTAMPTZ DEFAULT NOW() NOT NULL,
    updated_at TIMESTAMPTZ DEFAULT NOW() NOT NULL,
    UNIQUE (
        kitchen_config_id,
        phase_number,
        reward_id
    )
);

CREATE TABLE IF NOT EXISTS user_kitchen_progress (
    id BIGSERIAL PRIMARY KEY,
    user_id BIGINT NOT NULL REFERENCES users (id) ON DELETE CASCADE,
    stage_id BIGINT NOT NULL REFERENCES game_stages (id) ON DELETE CASCADE,
    station_levels JSONB DEFAULT '[]'::jsonb NOT NULL,
    unlocked_stations JSONB DEFAULT '[]'::jsonb NOT NULL,
    created_at TIMESTAMPTZ DEFAULT NOW() NOT NULL,
    updated_at TIMESTAMPTZ DEFAULT NOW() NOT NULL,
    UNIQUE (user_id, stage_id)
);

CREATE TABLE IF NOT EXISTS user_kitchen_phase_progress (
    id BIGSERIAL PRIMARY KEY,
    user_id BIGINT NOT NULL REFERENCES users (id) ON DELETE CASCADE,
    kitchen_config_id BIGINT NOT NULL REFERENCES stage_kitchen_configs (id) ON DELETE CASCADE,
    current_phase INT DEFAULT 1 NOT NULL,
    completed_phases JSONB DEFAULT '[]'::jsonb NOT NULL,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    UNIQUE (user_id, kitchen_config_id)
);

CREATE TABLE IF NOT EXISTS user_kitchen_phase_reward_claims (
    id BIGSERIAL PRIMARY KEY,
    user_id BIGINT NOT NULL REFERENCES users (id) ON DELETE CASCADE,
    kitchen_config_id BIGINT NOT NULL REFERENCES stage_kitchen_configs (id) ON DELETE CASCADE,
    phase_number INT NOT NULL,
    reward_id BIGINT NOT NULL REFERENCES rewards (id) ON DELETE CASCADE,
    claimed_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    UNIQUE (
        user_id,
        kitchen_config_id,
        phase_number,
        reward_id
    )
);

CREATE INDEX idx_food_items_slug ON food_items (slug);

CREATE INDEX idx_stage_kitchen_config_stage ON stage_kitchen_configs (stage_id);

CREATE INDEX idx_kitchen_phase_rewards_config ON kitchen_phase_completion_rewards (kitchen_config_id);

CREATE INDEX idx_kitchen_phase_rewards_reward ON kitchen_phase_completion_rewards (reward_id);

CREATE INDEX idx_user_kitchen_progress_user ON user_kitchen_phase_progress (user_id);

CREATE INDEX idx_user_kitchen_progress_config ON user_kitchen_phase_progress (kitchen_config_id);

CREATE INDEX idx_user_kitchen_claims_user ON user_kitchen_phase_reward_claims (user_id);

CREATE INDEX idx_kitchen_transition_levels ON stage_kitchen_configs USING GIN (transition_phase_levels);

CREATE INDEX idx_kitchen_profit_multipliers ON stage_kitchen_configs USING GIN (phase_profit_multipliers);

CREATE INDEX idx_user_kitchen_completed_phases ON user_kitchen_phase_progress USING GIN (completed_phases);

COMMIT;