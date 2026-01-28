BEGIN;

CREATE TABLE IF NOT EXISTS reward_types (
    id BIGSERIAL PRIMARY KEY,
    slug VARCHAR(50) UNIQUE NOT NULL,
    name VARCHAR(100) NOT NULL,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE TABLE IF NOT EXISTS rewards (
    id BIGSERIAL PRIMARY KEY,
    reward_type_id BIGINT NOT NULL REFERENCES reward_types(id) ON DELETE CASCADE,
    slug VARCHAR(50) UNIQUE NOT NULL,
    name VARCHAR(100) UNIQUE NOT NULL,
    amount INT DEFAULT 0 NOT NULL,
    is_active BOOLEAN DEFAULT true NOT NULL,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE TABLE IF NOT EXISTS daily_rewards (
    id BIGSERIAL PRIMARY KEY,
    reward_id BIGINT NOT NULL REFERENCES rewards(id) ON DELETE CASCADE,
    day_number INT UNIQUE NOT NULL,
    is_active BOOLEAN DEFAULT true,
    description TEXT DEFAULT '',
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),

    CONSTRAINT check_day_number_positive CHECK (day_number > 0)
);

CREATE TABLE IF NOT EXISTS user_daily_reward_progress (
    id BIGSERIAL PRIMARY KEY,
    user_id BIGINT NOT NULL UNIQUE REFERENCES users(id) ON DELETE CASCADE,
    current_streak INT DEFAULT 0,
    longest_streak INT DEFAULT 0,
    last_claim_date DATE,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),

    CONSTRAINT check_streak_non_negative CHECK (current_streak >= 0)
);

CREATE INDEX IF NOT EXISTS idx_daily_rewards_day_active ON daily_rewards(day_number, is_active);

CREATE INDEX IF NOT EXISTS idx_daily_rewards_created_at ON daily_rewards(created_at);

CREATE INDEX IF NOT EXISTS idx_user_daily_reward_progress_user_id ON user_daily_reward_progress(user_id);

CREATE INDEX IF NOT EXISTS idx_user_daily_reward_progress_last_claim ON user_daily_reward_progress(last_claim_date);

COMMIT;
