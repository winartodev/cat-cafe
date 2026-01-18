BEGIN;

CREATE TABLE IF NOT EXISTS reward_types (
    id BIGSERIAL PRIMARY KEY,
    slug VARCHAR(50) UNIQUE NOT NULL,
    name VARCHAR(100) NOT NULL,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE TABLE IF NOT EXISTS daily_rewards (
    id BIGSERIAL PRIMARY KEY,
    reward_type_id BIGINT NOT NULL REFERENCES reward_types(id) ON DELETE CASCADE,
    day_number INT NOT NULL,
    reward_amount INT NOT NULL,
    is_active BOOLEAN DEFAULT true,
    description TEXT DEFAULT '',
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE TABLE IF NOT EXISTS user_daily_rewards (
    id BIGSERIAL PRIMARY KEY,
    user_id BIGINT NOT NULL UNIQUE REFERENCES users(id) ON DELETE CASCADE,
    current_streak INT DEFAULT 0,
    last_claim_date DATE,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE INDEX IF NOT EXISTS idx_daily_rewards_day_active ON daily_rewards(day_number, is_active);

CREATE INDEX IF NOT EXISTS idx_daily_rewards_created_at ON daily_rewards(created_at);

CREATE INDEX IF NOT EXISTS idx_user_daily_rewards_user_id ON user_daily_rewards(user_id);

CREATE INDEX IF NOT EXISTS idx_user_daily_rewards_last_claim ON user_daily_rewards(last_claim_date);

COMMIT;
