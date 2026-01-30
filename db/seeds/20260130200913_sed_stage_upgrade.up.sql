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
    (
        'CAT_FRIEND', 
        'Cat Friend', 
        'Add +1 Helper', 
        15, 'coin', 'ADD_HELPER', 1, 'count', 'food', 
        (SELECT id FROM food_items WHERE slug = 'DIM_SUM_MENTAI'), 1
    ),
    (
        'PUT_OUT_PAMPHLET', 
        'Put out Pamphlet', 
        'Add +1 Customer slot', 
        40, 'coin', 'ADD_CUSTOMER_SLOT', 1, 'count', 'food', 
        (SELECT id FROM food_items WHERE slug = 'DIM_SUM_MENTAI'), 2
    ),
    (
        'TOO_MUCH_SAUCE', 
        'Too Much Sauce!', 
        'Reduce Dim Sum preparation time by half', 
        90, 'coin', 'REDUCE_COOKING_TIME', 0.5, 'multiplier', 'food', 
        (SELECT id FROM food_items WHERE slug = 'DIM_SUM_MENTAI'), 3
    ),
    (
        'GETTING_POPULAR', 
        'Getting Popular', 
        'Add +1 Customer slot', 
        230, 'coin', 'ADD_CUSTOMER_SLOT', 1, 'count', 'food', 
        (SELECT id FROM food_items WHERE slug = 'DIM_SUM_MENTAI'), 4
    ),
    (
        'JOB_FAIR_DAY', 
        'Job Fair Day', 
        'Add +1 Helper', 
        130, 'coin', 'ADD_HELPER', 1, 'count', 'food', 
        (SELECT id FROM food_items WHERE slug = 'DIM_SUM_MENTAI'), 5
    ),
    (
        'CAT_BROTHER', 
        'Cat Brother', 
        'Add +1 Helper', 
        38, 'coin', 'ADD_HELPER', 1, 'count', 'food', 
        (SELECT id FROM food_items WHERE slug = 'BOBA_TEA'), 7
    ),
    (
        'CAT_SISTER', 
        'Cat Sister', 
        'Add +1 Helper', 
        100, 'coin', 'ADD_HELPER', 1, 'count', 'food', 
        (SELECT id FROM food_items WHERE slug = 'BOBA_TEA'), 8
    ),
    (
        'YARD_SALE', 
        'Yard Sale', 
        'Add +1 Customer slot', 
        160, 'coin', 'ADD_CUSTOMER_SLOT', 1, 'count', 'food', 
        (SELECT id FROM food_items WHERE slug = 'BOBA_TEA'), 9
    ),
    (
        'BOBA_FOUNTAIN', 
        'Boba Fountain', 
        'Reduce Boba Tea time by half', 
        5169, 'coin', 'REDUCE_COOKING_TIME', 0.5, 'multiplier', 'food', 
        (SELECT id FROM food_items WHERE slug = 'BOBA_TEA'), 10
    ),
    (
        'CAT_COUSIN', 
        'Cat Cousin', 
        'Add +1 Helper', 
        334, 'coin', 'ADD_HELPER', 1, 'count', 'food', 
        (SELECT id FROM food_items WHERE slug = 'SALTED_BREAD'), 11
    ),
    (
        'ADVERTISE_YOUR_PLACE', 
        'Advertise Your Place', 
        'Add +1 Customer slot', 
        481, 'coin', 'ADD_CUSTOMER_SLOT', 1, 'count', 'food', 
        (SELECT id FROM food_items WHERE slug = 'SALTED_BREAD'), 12
    ),
    (
        'BRING_FAMILY', 
        'Bring Family', 
        'Add +1 Customer slot', 
        2070, 'coin', 'ADD_CUSTOMER_SLOT', 1, 'count', 'food', 
        (SELECT id FROM food_items WHERE slug = 'SALTED_BREAD'), 13
    ),
    (
        'MASTER_OF_SALT_BREAD', 
        'Master of Salt Bread', 
        'Reduce Salt Bread time by half', 
        26880, 'coin', 'REDUCE_COOKING_TIME', 0.5, 'multiplier', 'food', 
        (SELECT id FROM food_items WHERE slug = 'SALTED_BREAD'), 14
    ),
    (
        'NEW_RESTAURANT_2', 
        'New Restaurant', 
        'UNLOCK NEW RESTAURANT', 1035175, 
        'coin', 'UNLOCK_RESTAURANT', 1, 'boolean', 'restaurant ', 
        0, 15
    ),
    (
        'NEW_RESTAURANT', 
        'New Restaurant', 
        'UNLOCK NEW RESTAURANT', 
        2800, 'coin', 'UNLOCK_RESTAURANT', 1, 'boolean', 'restaurant', 
        0, 6 
    )
ON CONFLICT (slug) DO UPDATE 
SET 
    name = EXCLUDED.name,
    cost = EXCLUDED.cost,
    description = EXCLUDED.description,
    effect_value = EXCLUDED.effect_value,
    updated_at = NOW();


COMMIT;