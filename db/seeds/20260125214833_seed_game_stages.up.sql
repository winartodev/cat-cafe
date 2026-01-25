BEGIN;

-- Insert Game Stages
INSERT INTO game_stages (slug, name, starting_coin, stage_prize, is_active, sequence)
VALUES
    ('STG0001', 'DIMSUM CART', 0, 1500, true, 1),
    ('STG0002', 'BOBA STALL', 10, 250000, true, 2),
    ('STG0003', 'DESERT CAFE', 20, 1500000, true, 3);

-- Insert Customer Configs
INSERT INTO stage_customer_configs (stage_id, customer_spawn_time, max_customer_order_count, max_customer_order_variant, starting_order_table_count)
SELECT id, 2, 1, 1, 1 FROM game_stages WHERE slug = 'STG0001';
INSERT INTO stage_customer_configs (stage_id, customer_spawn_time, max_customer_order_count, max_customer_order_variant, starting_order_table_count)
SELECT id, 1.5, 2, 1, 1 FROM game_stages WHERE slug = 'STG0002';
INSERT INTO stage_customer_configs (stage_id, customer_spawn_time, max_customer_order_count, max_customer_order_variant, starting_order_table_count)
SELECT id, 1.5, 3, 1, 1 FROM game_stages WHERE slug = 'STG0003';

-- Insert Staff Configs
INSERT INTO stage_staff_configs (stage_id, starting_staff_manager, starting_staff_helper)
SELECT id, 'General', '' FROM game_stages WHERE slug IN ('STG0001', 'STG0002', 'STG0003');

-- Insert Camera Configs
INSERT INTO stage_camera_configs (stage_id, zoom_size, min_bound_x, min_bound_y, max_bound_x, max_bound_y)
SELECT id, 9.0, 0.0, 0.0, 0.0, 0.0 FROM game_stages WHERE slug = 'STG0001';
INSERT INTO stage_camera_configs (stage_id, zoom_size, min_bound_x, min_bound_y, max_bound_x, max_bound_y)
SELECT id, 10.5, 0.0, 0.0, 0.0, 0.0 FROM game_stages WHERE slug = 'STG0002';
INSERT INTO stage_camera_configs (stage_id, zoom_size, min_bound_x, min_bound_y, max_bound_x, max_bound_y)
SELECT id, 9.5, 0.0, -100.0, 15.0, 0.0 FROM game_stages WHERE slug = 'STG0003';

-- Insert Kitchen Configs
INSERT INTO stage_kitchen_configs
(stage_id, max_level, upgrade_profit_multiply, upgrade_cost_multiply, transition_phase_levels, phase_profit_multipliers, phase_upgrade_cost_multipliers, table_count_per_phases)
VALUES
    ((SELECT id FROM game_stages WHERE slug = 'STG0001'), 25, 115, 120, '{1, 13}', '{1, 2}', '{1, 2.1}', '{1, 2}'),
    ((SELECT id FROM game_stages WHERE slug = 'STG0002'), 50, 115, 120, '{1, 17, 33}', '{1, 2, 4}', '{1, 2.1, 4.2}', '{1, 2, 2}'),
    ((SELECT id FROM game_stages WHERE slug = 'STG0003'), 75, 115, 120, '{1, 20, 39, 58}', '{1, 2, 4, 6}', '{1, 2.1, 4.2, 6.2}', '{1, 2, 2, 2}');

-- Insert Kitchen Stations
INSERT INTO kitchen_stations (stage_id, food_item_id, auto_unlock)
VALUES
    ((SELECT id FROM game_stages WHERE slug = 'STG0001'), (SELECT id FROM food_items WHERE slug = 'DIM_SUM_MENTAI'), true),
    ((SELECT id FROM game_stages WHERE slug = 'STG0002'), (SELECT id FROM food_items WHERE slug = 'BOBA_TEA'), false),
    ((SELECT id FROM game_stages WHERE slug = 'STG0002'), (SELECT id FROM food_items WHERE slug = 'SALTED_BREAD'), false),
    ((SELECT id FROM game_stages WHERE slug = 'STG0003'), (SELECT id FROM food_items WHERE slug = 'MATCHA_LATTE'), false),
    ((SELECT id FROM game_stages WHERE slug = 'STG0003'), (SELECT id FROM food_items WHERE slug = 'STRAWBERRY_WAFFLE'), false),
    ((SELECT id FROM game_stages WHERE slug = 'STG0003'), (SELECT id FROM food_items WHERE slug = 'BURNT_CHEESECAKE'), false);

-- Insert Phase Rewards
INSERT INTO kitchen_phase_completion_rewards (kitchen_config_id, phase_number, reward_id)
VALUES
    ((SELECT id FROM stage_kitchen_configs WHERE stage_id = (SELECT id FROM game_stages WHERE slug = 'STG0001')), 1, (SELECT id FROM rewards WHERE slug = 'COIN_1_PACK')),
    ((SELECT id FROM stage_kitchen_configs WHERE stage_id = (SELECT id FROM game_stages WHERE slug = 'STG0001')), 2, (SELECT id FROM rewards WHERE slug = 'GEM_1_PACK')),
    ((SELECT id FROM stage_kitchen_configs WHERE stage_id = (SELECT id FROM game_stages WHERE slug = 'STG0002')), 1, (SELECT id FROM rewards WHERE slug = 'COIN_2_PACK')),
    ((SELECT id FROM stage_kitchen_configs WHERE stage_id = (SELECT id FROM game_stages WHERE slug = 'STG0002')), 2, (SELECT id FROM rewards WHERE slug = 'GEM_2_PACK')),
    ((SELECT id FROM stage_kitchen_configs WHERE stage_id = (SELECT id FROM game_stages WHERE slug = 'STG0002')), 3, (SELECT id FROM rewards WHERE slug = 'GOPAY_COIN_2_PACK')),
    ((SELECT id FROM stage_kitchen_configs WHERE stage_id = (SELECT id FROM game_stages WHERE slug = 'STG0003')), 1, (SELECT id FROM rewards WHERE slug = 'COIN_3_PACK')),
    ((SELECT id FROM stage_kitchen_configs WHERE stage_id = (SELECT id FROM game_stages WHERE slug = 'STG0003')), 2, (SELECT id FROM rewards WHERE slug = 'GEM_3_PACK')),
    ((SELECT id FROM stage_kitchen_configs WHERE stage_id = (SELECT id FROM game_stages WHERE slug = 'STG0003')), 3, (SELECT id FROM rewards WHERE slug = 'GOPAY_COIN_3_PACK')),
    ((SELECT id FROM stage_kitchen_configs WHERE stage_id = (SELECT id FROM game_stages WHERE slug = 'STG0003')), 4, (SELECT id FROM rewards WHERE slug = 'COIN_4_PACK'));

COMMIT;
