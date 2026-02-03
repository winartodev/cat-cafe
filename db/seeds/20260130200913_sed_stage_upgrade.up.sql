BEGIN;

INSERT INTO upgrades (
    slug, 
    name, 
    description, 
    cost, 
    cost_type, 
    effect_type, 
    effect_value, 
    effect_unit, 
    effect_target, 
    effect_target_id, 
    sequence
)
VALUES
    -- STG0001 Upgrades (Dim Sum)
    (
        'CAT_FRIEND', 'Cat Friend', 'Add +1 Helper',
        15, 'coin', 'add_helper', 1, 'count', 'food',
        (SELECT id FROM food_items WHERE slug = 'DIM_SUM_MENTAI'), 5
    ),
    (
        'PUT_OUT_PAMPHLET', 'Put out Pamphlet', 'Add +1 Customer slot',
        40, 'coin', 'add_customer', 1, 'count', 'food',
        (SELECT id FROM food_items WHERE slug = 'DIM_SUM_MENTAI'), 10
    ),
    (
        'TOO_MUCH_SAUCE', 'Too Much Sauce!', 'Reduce Dim Sum preparation time by half',
        90, 'coin', 'reduce_cooking_time', 50, 'percent', 'food',
        (SELECT id FROM food_items WHERE slug = 'DIM_SUM_MENTAI'), 15
    ),
    (
        'GETTING_POPULAR', 'Getting Popular', 'Add +1 Customer slot',
        230, 'coin', 'add_customer', 1, 'count', 'food',
        (SELECT id FROM food_items WHERE slug = 'DIM_SUM_MENTAI'), 20
    ),
    (
        'JOB_FAIR_DAY', 'Job Fair Day', 'Add +1 Helper',
        130, 'coin', 'add_helper', 1, 'count', 'food',
        (SELECT id FROM food_items WHERE slug = 'DIM_SUM_MENTAI'), 17
    ),

    -- STG0002 Upgrades (Boba Tea & Salt Bread)
    (
        'CAT_BROTHER', 'Cat Brother', 'Add +1 Helper',
        38, 'coin', 'add_helper', 1, 'count', 'food',
        (SELECT id FROM food_items WHERE slug = 'BOBA_TEA'), 12
    ),
    (
        'CAT_SISTER', 'Cat Sister', 'Add +1 Helper',
        100, 'coin', 'add_helper', 1, 'count', 'food',
        (SELECT id FROM food_items WHERE slug = 'BOBA_TEA'), 17
    ),
    (
        'YARD_SALE', 'Yard Sale', 'Add +1 Customer slot',
        160, 'coin', 'add_customer', 1, 'count', 'food',
        (SELECT id FROM food_items WHERE slug = 'BOBA_TEA'), 20
    ),
    (
        'BOBA_FOUNTAIN', 'Boba Fountain', 'Reduce Boba Tea time by half',
        5169, 'coin', 'reduce_cooking_time', 50, 'percent', 'food',
        (SELECT id FROM food_items WHERE slug = 'BOBA_TEA'), 35
    ),
    (
        'CAT_COUSIN', 'Cat Cousin', 'Add +1 Helper',
        334, 'coin', 'add_helper', 1, 'count', 'food',
        (SELECT id FROM food_items WHERE slug = 'SALTED_BREAD'), 15
    ),
    (
        'ADVERTISE_YOUR_PLACE', 'Advertise Your Place', 'Add +1 Customer slot',
        481, 'coin', 'add_customer', 1, 'count', 'food',
        (SELECT id FROM food_items WHERE slug = 'SALTED_BREAD'), 17
    ),
    (
        'BRING_FAMILY', 'Bring Family', 'Add +1 Customer slot',
        2070, 'coin', 'add_customer', 1, 'count', 'food',
        (SELECT id FROM food_items WHERE slug = 'SALTED_BREAD'), 25
    ),
    (
        'MASTER_OF_SALTED_BREAD', 'Master of Salt Bread', 'Reduce Salt Bread time by half',
        26880, 'coin', 'reduce_cooking_time', 50, 'percent', 'food',
        (SELECT id FROM food_items WHERE slug = 'SALTED_BREAD'), 35
    ),

    -- STG0003 Upgrades (Matcha, Waffle, Cheesecake)
    (
        'JIM_WANNA_JOIN', 'Jim wanna join', 'Add +1 Helper',
        150, 'coin', 'add_helper', 1, 'count', 'food',
        (SELECT id FROM food_items WHERE slug = 'MATCHA_LATTE'), 12
    ),
    (
        'IMPRESSIVE_CV', 'Impressive CV', 'Add +1 Helper',
        640, 'coin', 'add_helper', 1, 'count', 'food',
        (SELECT id FROM food_items WHERE slug = 'MATCHA_LATTE'), 20
    ),
    (
        'RENOVATION_BEGIN', 'Renovation begin', 'Add +1 Customer slot',
        370, 'coin', 'add_customer', 1, 'count', 'food',
        (SELECT id FROM food_items WHERE slug = 'MATCHA_LATTE'), 17
    ),
    (
        'MASTER_OF_ART_LATTE', 'Master of Art Latte', 'Reduce Matcha Latte time by half',
        3330000, 'coin', 'reduce_cooking_time', 50, 'percent', 'food',
        (SELECT id FROM food_items WHERE slug = 'MATCHA_LATTE'), 55
    ),
    (
        'BEST_SELLER_LATTE', 'Best Seller Latte', 'Matcha Latte profit +150%',
        8290000, 'coin', 'profit', 150, 'percent', 'food',
        (SELECT id FROM food_items WHERE slug = 'MATCHA_LATTE'), 60
    ),
    (
        'HIRING_DAY', 'Hiring day', 'Add +1 Helper',
        2040, 'coin', 'add_helper', 1, 'count', 'food',
        (SELECT id FROM food_items WHERE slug = 'STRAWBERRY_WAFFLE'), 17
    ),
    (
        'INTERN_BOY', 'Intern boy', 'Add +1 Helper',
        1420, 'coin', 'add_helper', 1, 'count', 'food',
        (SELECT id FROM food_items WHERE slug = 'STRAWBERRY_WAFFLE'), 15
    ),
    (
        'RENOVATION_ENDED', 'Renovation ended', 'Add +1 Customer slot',
        8800, 'coin', 'add_customer', 1, 'count', 'food',
        (SELECT id FROM food_items WHERE slug = 'STRAWBERRY_WAFFLE'), 25
    ),
    (
        'PROFESSIONAL_BAKER', 'Professional Baker', 'Reduce Strawberry Waffle time by half',
        18310000, 'coin', 'reduce_cooking_time', 50, 'percent', 'food',
        (SELECT id FROM food_items WHERE slug = 'STRAWBERRY_WAFFLE'), 55
    ),
    (
        'ENCHANTED_WAFFLE', 'Enchanted Waffle', 'Strawberry Waffle profit +150%',
        45560000, 'coin', 'profit', 150, 'percent', 'food',
        (SELECT id FROM food_items WHERE slug = 'STRAWBERRY_WAFFLE'), 60
    ),
    (
        'NEW_SPOT_IN_TOWN', 'New spot in town', 'Add +1 Customer slot',
        48000, 'coin', 'add_customer', 1, 'count', 'food',
        (SELECT id FROM food_items WHERE slug = 'BURNT_CHEESECAKE'), 20
    ),
    (
        'GETTING_VIRAL', 'Getting viral', 'Add +1 Customer slot',
        1560000, 'coin', 'add_customer', 1, 'count', 'food',
        (SELECT id FROM food_items WHERE slug = 'BURNT_CHEESECAKE'), 35
    ),
    (
        'SUPERSPEED_EQUIPMENT', 'Superspeed Equipment', 'Reduce Burnt Cheesecake time by half',
        249660000, 'coin', 'reduce_cooking_time', 50, 'percent', 'food',
        (SELECT id FROM food_items WHERE slug = 'BURNT_CHEESECAKE'), 55
    ),
    (
        'FRESH_FROM_THE_OVEN', 'Fresh from the Oven', 'Burnt Cheesecake profit +150%',
        621240000, 'coin', 'profit', 150, 'percent', 'food',
        (SELECT id FROM food_items WHERE slug = 'BURNT_CHEESECAKE'), 60
    ),

    -- Restaurant Unlock Upgrades
    (
        'NEW_RESTAURANT_1', 'New Restaurant', 'UNLOCK NEW RESTAURANT',
        2800, 'coin', 'unlock_restaurant', 1, 'boolean', 'restaurant',
        0, 25
    ),
    (
        'NEW_RESTAURANT_2', 'New Restaurant', 'UNLOCK NEW RESTAURANT',
        1035175, 'coin', 'unlock_restaurant', 1, 'boolean', 'restaurant',
        0, 50
    ),
    (
        'NEW_RESTAURANT_3', 'New Restaurant', 'UNLOCK NEW RESTAURANT',
        17530911000000, 'coin', 'unlock_restaurant', 1, 'boolean', 'restaurant',
        0, 100
    )
ON CONFLICT (slug) DO UPDATE 
SET 
    name = EXCLUDED.name,
    cost = EXCLUDED.cost,
    description = EXCLUDED.description,
    effect_value = EXCLUDED.effect_value,
    updated_at = NOW();

INSERT INTO game_stage_upgrades (game_stage_id, upgrade_id)
VALUES
    -- STG0001 Upgrades
    ((SELECT id FROM game_stages WHERE slug = 'STG0001'), (SELECT id FROM upgrades WHERE slug = 'CAT_FRIEND')),
    ((SELECT id FROM game_stages WHERE slug = 'STG0001'), (SELECT id FROM upgrades WHERE slug = 'PUT_OUT_PAMPHLET')),
    ((SELECT id FROM game_stages WHERE slug = 'STG0001'), (SELECT id FROM upgrades WHERE slug = 'TOO_MUCH_SAUCE')),
    ((SELECT id FROM game_stages WHERE slug = 'STG0001'), (SELECT id FROM upgrades WHERE slug = 'GETTING_POPULAR')),
    ((SELECT id FROM game_stages WHERE slug = 'STG0001'), (SELECT id FROM upgrades WHERE slug = 'JOB_FAIR_DAY')),
    ((SELECT id FROM game_stages WHERE slug = 'STG0001'), (SELECT id FROM upgrades WHERE slug = 'NEW_RESTAURANT_1')),

    -- STG0002 Upgrades
    ((SELECT id FROM game_stages WHERE slug = 'STG0002'), (SELECT id FROM upgrades WHERE slug = 'CAT_BROTHER')),
    ((SELECT id FROM game_stages WHERE slug = 'STG0002'), (SELECT id FROM upgrades WHERE slug = 'CAT_SISTER')),
    ((SELECT id FROM game_stages WHERE slug = 'STG0002'), (SELECT id FROM upgrades WHERE slug = 'YARD_SALE')),
    ((SELECT id FROM game_stages WHERE slug = 'STG0002'), (SELECT id FROM upgrades WHERE slug = 'CAT_COUSIN')),
    ((SELECT id FROM game_stages WHERE slug = 'STG0002'), (SELECT id FROM upgrades WHERE slug = 'ADVERTISE_YOUR_PLACE')),
    ((SELECT id FROM game_stages WHERE slug = 'STG0002'), (SELECT id FROM upgrades WHERE slug = 'BRING_FAMILY')),
    ((SELECT id FROM game_stages WHERE slug = 'STG0002'), (SELECT id FROM upgrades WHERE slug = 'BOBA_FOUNTAIN')),
    ((SELECT id FROM game_stages WHERE slug = 'STG0002'), (SELECT id FROM upgrades WHERE slug = 'MASTER_OF_SALTED_BREAD')),
    ((SELECT id FROM game_stages WHERE slug = 'STG0002'), (SELECT id FROM upgrades WHERE slug = 'NEW_RESTAURANT_2')),

    -- STG0003 Upgrades
    ((SELECT id FROM game_stages WHERE slug = 'STG0003'), (SELECT id FROM upgrades WHERE slug = 'JIM_WANNA_JOIN')),
    ((SELECT id FROM game_stages WHERE slug = 'STG0003'), (SELECT id FROM upgrades WHERE slug = 'IMPRESSIVE_CV')),
    ((SELECT id FROM game_stages WHERE slug = 'STG0003'), (SELECT id FROM upgrades WHERE slug = 'RENOVATION_BEGIN')),
    ((SELECT id FROM game_stages WHERE slug = 'STG0003'), (SELECT id FROM upgrades WHERE slug = 'MASTER_OF_ART_LATTE')),
    ((SELECT id FROM game_stages WHERE slug = 'STG0003'), (SELECT id FROM upgrades WHERE slug = 'BEST_SELLER_LATTE')),
    ((SELECT id FROM game_stages WHERE slug = 'STG0003'), (SELECT id FROM upgrades WHERE slug = 'HIRING_DAY')),
    ((SELECT id FROM game_stages WHERE slug = 'STG0003'), (SELECT id FROM upgrades WHERE slug = 'INTERN_BOY')),
    ((SELECT id FROM game_stages WHERE slug = 'STG0003'), (SELECT id FROM upgrades WHERE slug = 'RENOVATION_ENDED')),
    ((SELECT id FROM game_stages WHERE slug = 'STG0003'), (SELECT id FROM upgrades WHERE slug = 'ENCHANTED_WAFFLE')),
    ((SELECT id FROM game_stages WHERE slug = 'STG0003'), (SELECT id FROM upgrades WHERE slug = 'PROFESSIONAL_BAKER')),
    ((SELECT id FROM game_stages WHERE slug = 'STG0003'), (SELECT id FROM upgrades WHERE slug = 'NEW_SPOT_IN_TOWN')),
    ((SELECT id FROM game_stages WHERE slug = 'STG0003'), (SELECT id FROM upgrades WHERE slug = 'GETTING_VIRAL')),
    ((SELECT id FROM game_stages WHERE slug = 'STG0003'), (SELECT id FROM upgrades WHERE slug = 'SUPERSPEED_EQUIPMENT')),
    ((SELECT id FROM game_stages WHERE slug = 'STG0003'), (SELECT id FROM upgrades WHERE slug = 'FRESH_FROM_THE_OVEN')),
    ((SELECT id FROM game_stages WHERE slug = 'STG0003'), (SELECT id FROM upgrades WHERE slug = 'NEW_RESTAURANT_3'))
ON CONFLICT (game_stage_id, upgrade_id) DO NOTHING;


COMMIT;
