BEGIN;

INSERT INTO reward_types (slug, name, created_at, updated_at)
VALUES
    ('COIN', 'Coin', NOW(), NOW()),
    ('GEM', 'Gem', NOW(), NOW()),
    ('GOPAY_COIN', 'Gopay Coin', NOW(), NOW())
ON CONFLICT (slug) DO NOTHING;

INSERT INTO rewards (reward_type_id, slug, name, amount, is_active, created_at, updated_at)
SELECT
    rt.id,
    rt.slug || '_' || s.val || '_PACK',
    s.val || ' ' || rt.name || ' Reward',
    (s.val * 10),
    true,
    NOW(),
    NOW()
FROM reward_types rt
         CROSS JOIN generate_series(1, 10) AS s(val)
WHERE rt.slug IN ('COIN', 'GEM', 'GOPAY_COIN')
ON CONFLICT DO NOTHING;

INSERT INTO daily_rewards (reward_id, day_number, is_active, description, created_at, updated_at)
SELECT
    -- Picks a random reward_id from the rewards table for each day
    (SELECT id FROM rewards ORDER BY random() LIMIT 1),
    s.day,
    true,
    'Reward for Day ' || s.day,
    NOW(),
    NOW()
FROM generate_series(1, 7) AS s(day)
ON CONFLICT (day_number) DO NOTHING;

COMMIT;
