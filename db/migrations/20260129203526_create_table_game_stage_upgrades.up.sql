BEGIN;

CREATE TABLE IF NOT EXISTS upgrades (
    id BIGSERIAL PRIMARY KEY,
    slug VARCHAR(50) UNIQUE NOT NULL,
    name VARCHAR(100) NOT NULL,
    description TEXT DEFAULT '' NOT NULL,
    cost BIGINT NOT NULL,
    cost_type VARCHAR(50) DEFAULT 'coin' NOT NULL,
    effect_type VARCHAR(50) NOT NULL,
    effect_value NUMERIC NOT NULL,
    effect_unit VARCHAR(50) NOT NULL,
    effect_target VARCHAR(50) NOT NULL,
    effect_target_id BIGINT NOT NULL,
    is_active BOOLEAN DEFAULT true NOT NULL,
    sequence INT DEFAULT 0 NOT NULL,
    created_at TIMESTAMPTZ DEFAULT NOW() NOT NULL,
    updated_at TIMESTAMPTZ DEFAULT NOW() NOT NULL
);

CREATE TABLE IF NOT EXISTS game_stage_upgrades (
    id BIGSERIAL PRIMARY KEY,
    game_stage_id BIGINT NOT NULL REFERENCES game_stages(id) ON DELETE CASCADE,
    upgrade_id BIGINT NOT NULL REFERENCES upgrades(id) ON DELETE CASCADE,
    created_at TIMESTAMPTZ DEFAULT NOW() NOT NULL,
    updated_at TIMESTAMPTZ DEFAULT NOW() NOT NULL,
    UNIQUE(game_stage_id, upgrade_id)
);

CREATE TABLE IF NOT EXISTS user_stage_upgrades (
    id BIGSERIAL PRIMARY KEY,
    user_id BIGINT NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    game_stage_id BIGINT NOT NULL REFERENCES game_stages(id) ON DELETE CASCADE,
    game_stage_upgrade_id BIGINT NOT NULL REFERENCES game_stage_upgrades(id) ON DELETE CASCADE,
    purchased_at TIMESTAMPTZ DEFAULT NOW() NOT NULL,
    created_at TIMESTAMPTZ DEFAULT NOW() NOT NULL,
    updated_at TIMESTAMPTZ DEFAULT NOW() NOT NULL,
    UNIQUE(user_id, game_stage_upgrade_id)
);

CREATE INDEX idx_upgrades_slug ON upgrades(slug);
CREATE INDEX idx_upgrades_effect ON upgrades(effect_target, effect_target_id);
CREATE INDEX idx_upgrades_active_sequence ON upgrades(is_active, sequence) WHERE is_active = true;

CREATE INDEX idx_game_stage_upgrades_game_stage ON game_stage_upgrades(game_stage_id);
CREATE INDEX idx_game_stage_upgrades_upgrade ON game_stage_upgrades(upgrade_id);

CREATE INDEX idx_user_stage_upgrades_user ON user_stage_upgrades(user_id);
CREATE INDEX idx_user_stage_upgrades_stage ON user_stage_upgrades(game_stage_id);
CREATE INDEX idx_user_stage_upgrades_upgrade ON user_stage_upgrades(game_stage_upgrade_id);
CREATE INDEX idx_user_stage_upgrades_purchased ON user_stage_upgrades(purchased_at);

COMMIT;