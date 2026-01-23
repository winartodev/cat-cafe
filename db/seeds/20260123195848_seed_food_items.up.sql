BEGIN;

INSERT INTO food_items (slug, name, starting_price, starting_preparation)
VALUES
    ('DIM_SUM_MENTAI', 'Dim Sum Mentai', 4, 3),

    ('BOBA_TEA', 'Boba Tea', 2, 3),
    ('SALTED_BREAD', 'Salted Bread', 25, 7),

    ('MATCHA_LATTE', 'Matcha Latte', 15, 3),
    ('STRAWBERRY_WAFFLE', 'Strawberry Waffle', 200, 5),
    ('BURNT_CHEESECAKE', 'Burnt Cheesecake', 1200, 7),

    ('KLEPON', 'Klepon', 120, 5),
    ('ONDE_ONDE', 'Onde - Onde', 5000, 7),
    ('SERABI', 'Serabi', 120000, 9),
    ('GETUK_LINDRI', 'Getuk Lindri', 65000000, 11)
ON CONFLICT (slug) DO UPDATE
    SET
        name = EXCLUDED.name,
        starting_price = EXCLUDED.starting_price,
        starting_preparation = EXCLUDED.starting_preparation,
        updated_at = NOW();

COMMIT;
