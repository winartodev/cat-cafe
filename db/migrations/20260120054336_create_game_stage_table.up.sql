BEGIN;

CREATE TABLE IF NOT EXISTS game_stages (
    id BIGSERIAL PRIMARY KEY,
    slug VARCHAR(50) UNIQUE NOT NULL,
    name VARCHAR(100) NOT NULL,
    description TEXT DEFAULT '' NOT NULL,
    starting_coin BIGINT DEFAULT 0 NOT NULL,
    stage_prize BIGINT DEFAULT 0 NOT NULL,
    is_active BOOL DEFAULT true NOT NULL,
    sequence INT DEFAULT 0 NOT NULL,
    created_at TIMESTAMPTZ DEFAULT NOW() NOT NULL,
    updated_at TIMESTAMPTZ DEFAULT NOW() NOT NULL,

    CONSTRAINT check_sequence_non_negative CHECK (sequence >= 0),
    CONSTRAINT check_starting_coin_non_negative CHECK (starting_coin >= 0),
    CONSTRAINT check_stage_price_non_negative CHECK (stage_prize >= 0)
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

    UNIQUE(stage_id),

    CONSTRAINT check_customer_spawn_time_positive CHECK (customer_spawn_time > 0),
    CONSTRAINT check_max_order_counts_positive CHECK (
        max_customer_order_count > 0 AND
        max_customer_order_variant > 0 AND
        starting_order_table_count > 0
    )
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

CREATE UNIQUE INDEX unique_active_sequence_per_stage
    ON game_stages (sequence)
    WHERE is_active = true;

CREATE INDEX idx_game_stages_slug ON game_stages(slug);
CREATE INDEX idx_game_stages_active_sequence ON game_stages(is_active, sequence) WHERE is_active = true;

CREATE INDEX idx_stage_customer_config_stage ON stage_customer_configs(stage_id);
CREATE INDEX idx_stage_staff_config_stage ON stage_staff_configs(stage_id);

CREATE INDEX idx_stage_camera_config_stage ON stage_camera_configs(stage_id);

CREATE INDEX idx_user_stage_progress_user ON user_stage_progress(user_id);
CREATE INDEX idx_user_stage_progress_stage ON user_stage_progress(stage_id);

COMMIT;
