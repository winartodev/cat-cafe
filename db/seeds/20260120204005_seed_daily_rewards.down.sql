BEGIN;

DELETE FROM daily_rewards
WHERE day_number BETWEEN 1 AND 7;

DELETE FROM rewards
WHERE slug LIKE 'COIN_%'
   OR slug LIKE 'GEM_%'
   OR slug LIKE 'GOPAY_COIN_%';

DELETE FROM reward_types
WHERE slug IN ('COIN', 'GEM', 'GOPAY_COIN');

COMMIT;
