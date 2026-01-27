BEGIN;

INSERT INTO food_items (slug, name, initial_cost,initial_profit, cooking_time)
VALUES
    ('DIM_SUM_MENTAI', 'Dim Sum Mentai', 0, 4, 3),
    ('BOBA_TEA', 'Boba Tea', 5, 2, 3),
    ('SALTED_BREAD', 'Salted Bread', 26, 25, 7),
    ('MATCHA_LATTE', 'Matcha Latte', 20, 15, 3),
    ('STRAWBERRY_WAFFLE', 'Strawberry Waffle', 110, 200, 5),
    ('BURNT_CHEESECAKE', 'Burnt Cheesecake', 1500, 1200, 7),
    ('KLEPON', 'Klepon', 50, 120, 5),
    ('ONDE_ONDE', 'Onde - Onde', 260, 5000, 7),
    ('SERABI', 'Serabi', 3340, 120000, 9),
    ('GETUK_LINDRI', 'Getuk Lindri', 42890, 65000000, 11)
ON CONFLICT (slug) DO UPDATE
    SET
        name = EXCLUDED.name,
        initial_cost = EXCLUDED.initial_cost,
        initial_profit = EXCLUDED.initial_profit,
        cooking_time = EXCLUDED.cooking_time,
        updated_at = NOW();

WITH items AS (
    SELECT id, slug, cooking_time FROM food_items WHERE slug IN (
        'DIM_SUM_MENTAI', 'BOBA_TEA', 'SALTED_BREAD', 'MATCHA_LATTE', 
        'STRAWBERRY_WAFFLE', 'BURNT_CHEESECAKE', 'KLEPON', 'ONDE_ONDE', 
        'SERABI', 'GETUK_LINDRI'
    )
)
INSERT INTO food_level_overrides (food_item_id, level, cost, profit, preparation_time)
VALUES
    -- Dim Sum Mentai
    ((SELECT id FROM items WHERE slug = 'DIM_SUM_MENTAI'), 1, 0, 4, (SELECT cooking_time FROM items WHERE slug = 'DIM_SUM_MENTAI')),
    ((SELECT id FROM items WHERE slug = 'DIM_SUM_MENTAI'), 2, 3, 5, (SELECT cooking_time FROM items WHERE slug = 'DIM_SUM_MENTAI')),

    -- Boba Tea
    ((SELECT id FROM items WHERE slug = 'BOBA_TEA'), 1, 5, 5, (SELECT cooking_time FROM items WHERE slug = 'BOBA_TEA')),
    ((SELECT id FROM items WHERE slug = 'BOBA_TEA'), 2, 3, 6, (SELECT cooking_time FROM items WHERE slug = 'BOBA_TEA')),

    -- Salted Bread
    ((SELECT id FROM items WHERE slug = 'SALTED_BREAD'), 1, 26, 20, (SELECT cooking_time FROM items WHERE slug = 'SALTED_BREAD')),
    ((SELECT id FROM items WHERE slug = 'SALTED_BREAD'), 2, 3, 22, (SELECT cooking_time FROM items WHERE slug = 'SALTED_BREAD')),

    -- Matcha Latte
    ((SELECT id FROM items WHERE slug = 'MATCHA_LATTE'), 1, 20, 20, (SELECT cooking_time FROM items WHERE slug = 'MATCHA_LATTE')),
    ((SELECT id FROM items WHERE slug = 'MATCHA_LATTE'), 2, 3, 23, (SELECT cooking_time FROM items WHERE slug = 'MATCHA_LATTE')),

    -- Strawberry Waffle
    ((SELECT id FROM items WHERE slug = 'STRAWBERRY_WAFFLE'), 1, 110, 83, (SELECT cooking_time FROM items WHERE slug = 'STRAWBERRY_WAFFLE')),
    ((SELECT id FROM items WHERE slug = 'STRAWBERRY_WAFFLE'), 2, 3, 95, (SELECT cooking_time FROM items WHERE slug = 'STRAWBERRY_WAFFLE')),

    -- Burnt Cheesecake
    ((SELECT id FROM items WHERE slug = 'BURNT_CHEESECAKE'), 1, 1500, 1125, (SELECT cooking_time FROM items WHERE slug = 'BURNT_CHEESECAKE')),
    ((SELECT id FROM items WHERE slug = 'BURNT_CHEESECAKE'), 2, 3, 1294, (SELECT cooking_time FROM items WHERE slug = 'BURNT_CHEESECAKE')),

    -- Klepon
    ((SELECT id FROM items WHERE slug = 'KLEPON'), 1, 50, 50, (SELECT cooking_time FROM items WHERE slug = 'KLEPON')),
    ((SELECT id FROM items WHERE slug = 'KLEPON'), 2, 60, 58, (SELECT cooking_time FROM items WHERE slug = 'KLEPON')),

    -- Onde-Onde
    ((SELECT id FROM items WHERE slug = 'ONDE_ONDE'), 1, 260, 195, (SELECT cooking_time FROM items WHERE slug = 'ONDE_ONDE')),
    ((SELECT id FROM items WHERE slug = 'ONDE_ONDE'), 2, 312, 224, (SELECT cooking_time FROM items WHERE slug = 'ONDE_ONDE')),

    -- Serabi
    ((SELECT id FROM items WHERE slug = 'SERABI'), 1, 3340, 2505, (SELECT cooking_time FROM items WHERE slug = 'SERABI')),
    ((SELECT id FROM items WHERE slug = 'SERABI'), 2, 4008, 2881, (SELECT cooking_time FROM items WHERE slug = 'SERABI')),

    -- Getuk Lindri
    ((SELECT id FROM items WHERE slug = 'GETUK_LINDRI'), 1, 42890, 32168, (SELECT cooking_time FROM items WHERE slug = 'GETUK_LINDRI')),
    ((SELECT id FROM items WHERE slug = 'GETUK_LINDRI'), 2, 51468, 36993, (SELECT cooking_time FROM items WHERE slug = 'GETUK_LINDRI'))

ON CONFLICT (food_item_id, level) DO UPDATE 
SET 
    cost = EXCLUDED.cost, 
    profit = EXCLUDED.profit,
    preparation_time = EXCLUDED.preparation_time,
    updated_at = NOW();

COMMIT;
