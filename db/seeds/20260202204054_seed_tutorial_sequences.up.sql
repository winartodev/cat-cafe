-- Tutorial Sequences Database Seeding
-- Tutorial: Restaurant Introduction and Gameplay

-- ============================================
-- TUTORIAL 1: Introducing the Restaurant
-- Location: World Map
-- ============================================

-- Sequence 0: Cat Manager greeting
INSERT INTO tutorial_sequences (tutorial_key, location, sequence)
VALUES ('restaurant_introduction', 'world_map', 0);

INSERT INTO tutorial_sequence_translations (tutorial_sequence_id, language_code, title, description)
VALUES
    (currval('tutorial_sequences_id_seq'), 'id', 'Cat Manager', 'Hei hei, kamu! Lihat dong, kami baru aja buka kafe pertama kami, keren nggak tuh?'),
    (currval('tutorial_sequences_id_seq'), 'en', 'Cat Manager', 'Hey hey, you! Look, we just opened our first cafe, isn''t it cool?');

-- Sequence 1: FOMO explanation
INSERT INTO tutorial_sequences (tutorial_key, location, sequence)
VALUES ('restaurant_introduction', 'world_map', 1);

INSERT INTO tutorial_sequence_translations (tutorial_sequence_id, language_code, title, description)
VALUES
    (currval('tutorial_sequences_id_seq'), 'id', 'Cat Manager', 'Aku denger-denger, katanya sekarang buka kafe lagi ngetren, kami FOMO deh. Meow.'),
    (currval('tutorial_sequences_id_seq'), 'en', 'Cat Manager', 'I heard that opening cafes is trending now, we got FOMO. Meow.');

-- Sequence 2: Ask for help
INSERT INTO tutorial_sequences (tutorial_key, location, sequence)
VALUES ('restaurant_introduction', 'world_map', 2);

INSERT INTO tutorial_sequence_translations (tutorial_sequence_id, language_code, title, description)
VALUES
    (currval('tutorial_sequences_id_seq'), 'id', 'Cat Manager', 'Mau nggak bantuin kami ngurus kafenya?'),
    (currval('tutorial_sequences_id_seq'), 'en', 'Cat Manager', 'Would you like to help us manage the cafe?');

-- Sequence 3: Thank you
INSERT INTO tutorial_sequences (tutorial_key, location, sequence)
VALUES ('restaurant_introduction', 'world_map', 3);

INSERT INTO tutorial_sequence_translations (tutorial_sequence_id, language_code, title, description)
VALUES
    (currval('tutorial_sequences_id_seq'), 'id', 'Cat Manager', 'Kamu mau? Waah, makasih banget!'),
    (currval('tutorial_sequences_id_seq'), 'en', 'Cat Manager', 'You will? Wow, thank you so much!');

-- Sequence 4: Restaurant icon unlocked
INSERT INTO tutorial_sequences (tutorial_key, location, sequence)
VALUES ('restaurant_introduction', 'world_map', 4);

INSERT INTO tutorial_sequence_translations (tutorial_sequence_id, language_code, title, description)
VALUES
    (currval('tutorial_sequences_id_seq'), 'id', 'System', 'unlock_restaurant'),
    (currval('tutorial_sequences_id_seq'), 'en', 'System', 'unlock_restaurant');

-- Sequence 5: Let's go
INSERT INTO tutorial_sequences (tutorial_key, location, sequence)
VALUES ('restaurant_introduction', 'world_map', 5);

INSERT INTO tutorial_sequence_translations (tutorial_sequence_id, language_code, title, description)
VALUES
    (currval('tutorial_sequences_id_seq'), 'id', 'Cat Manager', 'Ayo deh!'),
    (currval('tutorial_sequences_id_seq'), 'en', 'Cat Manager', 'Let''s go!');

-- Sequence 6: Point to restaurant icon
INSERT INTO tutorial_sequences (tutorial_key, location, sequence)
VALUES ('restaurant_introduction', 'world_map', 6);

INSERT INTO tutorial_sequence_translations (tutorial_sequence_id, language_code, title, description)
VALUES
    (currval('tutorial_sequences_id_seq'), 'id', 'System', 'point_restaurant_icon'),
    (currval('tutorial_sequences_id_seq'), 'en', 'System', 'point_restaurant_icon');

-- Sequence 7: Player clicks restaurant (input action)
INSERT INTO tutorial_sequences (tutorial_key, location, sequence)
VALUES ('restaurant_introduction', 'world_map', 7);

INSERT INTO tutorial_sequence_translations (tutorial_sequence_id, language_code, title, description)
VALUES
    (currval('tutorial_sequences_id_seq'), 'id', 'Player', 'click_restaurant_icon'),
    (currval('tutorial_sequences_id_seq'), 'en', 'Player', 'click_restaurant_icon');


-- ============================================
-- TUTORIAL 2: Introducing Gameplay
-- Location: 1st Restaurant
-- ============================================

-- Sequence 0: Welcome to cafe
INSERT INTO tutorial_sequences (tutorial_key, location, sequence)
VALUES ('gameplay_introduction', 'STG0001', 0);

INSERT INTO tutorial_sequence_translations (tutorial_sequence_id, language_code, title, description)
VALUES
    (currval('tutorial_sequences_id_seq'), 'id', 'Cat Manager', 'Jeng jeng jeng, ini dia kafe kecil kesayangan kami!'),
    (currval('tutorial_sequences_id_seq'), 'en', 'Cat Manager', 'Ta-da, this is our beloved little cafe!');

-- Sequence 1: Customer enters
INSERT INTO tutorial_sequences (tutorial_key, location, sequence)
VALUES ('gameplay_introduction', 'STG0001', 1);

INSERT INTO tutorial_sequence_translations (tutorial_sequence_id, language_code, title, description)
VALUES
    (currval('tutorial_sequences_id_seq'), 'id', 'System', 'customer_enter'),
    (currval('tutorial_sequences_id_seq'), 'en', 'System', 'customer_enter');

-- Sequence 2: First customer
INSERT INTO tutorial_sequences (tutorial_key, location, sequence)
VALUES ('gameplay_introduction', 'STG0001', 2);

INSERT INTO tutorial_sequence_translations (tutorial_sequence_id, language_code, title, description)
VALUES
    (currval('tutorial_sequences_id_seq'), 'id', 'Cat Manager', 'Ooh, lihat lihat, pelanggan pertama nih, bentar ya aku tanya dia mau beli apa.'),
    (currval('tutorial_sequences_id_seq'), 'en', 'Cat Manager', 'Ooh, look look, our first customer, let me ask what they want to order.');

-- Sequence 3: Taking order
INSERT INTO tutorial_sequences (tutorial_key, location, sequence)
VALUES ('gameplay_introduction', 'STG0001', 3);

INSERT INTO tutorial_sequence_translations (tutorial_sequence_id, language_code, title, description)
VALUES
    (currval('tutorial_sequences_id_seq'), 'id', 'System', 'cat_take_order'),
    (currval('tutorial_sequences_id_seq'), 'en', 'System', 'cat_take_order');

-- Sequence 4: Order dimsum
INSERT INTO tutorial_sequences (tutorial_key, location, sequence)
VALUES ('gameplay_introduction', 'STG0001', 4);

INSERT INTO tutorial_sequence_translations (tutorial_sequence_id, language_code, title, description)
VALUES
    (currval('tutorial_sequences_id_seq'), 'id', 'Cat Manager', 'Dia mau beli Dimsum, katanya. Sekarang waktunya bikin dapur buat masak pesanannya!'),
    (currval('tutorial_sequences_id_seq'), 'en', 'Cat Manager', 'They want to buy Dimsum, they said. Now it''s time to build a kitchen to cook their order!');

-- Sequence 5: Highlight dimsum stand
INSERT INTO tutorial_sequences (tutorial_key, location, sequence)
VALUES ('gameplay_introduction', 'STG0001', 5);

INSERT INTO tutorial_sequence_translations (tutorial_sequence_id, language_code, title, description)
VALUES
    (currval('tutorial_sequences_id_seq'), 'id', 'System', 'highlight_dimsum_stand'),
    (currval('tutorial_sequences_id_seq'), 'en', 'System', 'highlight_dimsum_stand');

-- Sequence 6: Build dimsum stand (input action)
INSERT INTO tutorial_sequences (tutorial_key, location, sequence)
VALUES ('gameplay_introduction', 'STG0001', 6);

INSERT INTO tutorial_sequence_translations (tutorial_sequence_id, language_code, title, description)
VALUES
    (currval('tutorial_sequences_id_seq'), 'id', 'Player', 'build_dimsum_stand'),
    (currval('tutorial_sequences_id_seq'), 'en', 'Player', 'build_dimsum_stand');

-- Sequence 7: Cat delivers order
INSERT INTO tutorial_sequences (tutorial_key, location, sequence)
VALUES ('gameplay_introduction', 'STG0001', 7);

INSERT INTO tutorial_sequence_translations (tutorial_sequence_id, language_code, title, description)
VALUES
    (currval('tutorial_sequences_id_seq'), 'id', 'System', 'cat_deliver_order'),
    (currval('tutorial_sequences_id_seq'), 'en', 'System', 'cat_deliver_order');

-- Sequence 8: First customer served
INSERT INTO tutorial_sequences (tutorial_key, location, sequence)
VALUES ('gameplay_introduction', 'STG0001', 8);

INSERT INTO tutorial_sequence_translations (tutorial_sequence_id, language_code, title, description)
VALUES
    (currval('tutorial_sequences_id_seq'), 'id', 'Cat Manager', 'Wow kita baru aja melayani pelanggan pertama, keren!'),
    (currval('tutorial_sequences_id_seq'), 'en', 'Cat Manager', 'Wow we just served our first customer, awesome!');

-- Sequence 9: Upgrade reminder
INSERT INTO tutorial_sequences (tutorial_key, location, sequence)
VALUES ('gameplay_introduction', 'STG0001', 9);

INSERT INTO tutorial_sequence_translations (tutorial_sequence_id, language_code, title, description)
VALUES
    (currval('tutorial_sequences_id_seq'), 'id', 'Cat Manager', 'Mulai sekarang pasti lebih banyak pelanggan yang beli, jadi jangan lupa ningkatin level dapur kamu ya.'),
    (currval('tutorial_sequences_id_seq'), 'en', 'Cat Manager', 'From now on there will definitely be more customers buying, so don''t forget to level up your kitchen.');

-- Sequence 10: Highlight dimsum station
INSERT INTO tutorial_sequences (tutorial_key, location, sequence)
VALUES ('gameplay_introduction', 'STG0001', 10);

INSERT INTO tutorial_sequence_translations (tutorial_sequence_id, language_code, title, description)
VALUES
    (currval('tutorial_sequences_id_seq'), 'id', 'System', 'highlight_dimsum_station'),
    (currval('tutorial_sequences_id_seq'), 'en', 'System', 'highlight_dimsum_station');

-- Sequence 11: Tap dimsum station (input action)
INSERT INTO tutorial_sequences (tutorial_key, location, sequence)
VALUES ('gameplay_introduction', 'STG0001', 11);

INSERT INTO tutorial_sequence_translations (tutorial_sequence_id, language_code, title, description)
VALUES
    (currval('tutorial_sequences_id_seq'), 'id', 'Player', 'tap_dimsum_station'),
    (currval('tutorial_sequences_id_seq'), 'en', 'Player', 'tap_dimsum_station');

-- Sequence 12: Upgrade button explanation
INSERT INTO tutorial_sequences (tutorial_key, location, sequence)
VALUES ('gameplay_introduction', 'STG0001', 12);

INSERT INTO tutorial_sequence_translations (tutorial_sequence_id, language_code, title, description)
VALUES
    (currval('tutorial_sequences_id_seq'), 'id', 'Cat Manager', 'Pas duitnya cukup, kamu bisa pencet tombol ''Upgrade'' untuk naikin level dapur.'),
    (currval('tutorial_sequences_id_seq'), 'en', 'Cat Manager', 'When you have enough money, you can press the ''Upgrade'' button to level up the kitchen.');

-- Sequence 13: Station screen closed
INSERT INTO tutorial_sequences (tutorial_key, location, sequence)
VALUES ('gameplay_introduction', 'STG0001', 13);

INSERT INTO tutorial_sequence_translations (tutorial_sequence_id, language_code, title, description)
VALUES
    (currval('tutorial_sequences_id_seq'), 'id', 'System', 'close_station_screen'),
    (currval('tutorial_sequences_id_seq'), 'en', 'System', 'close_station_screen');

-- Sequence 14: Restaurant upgrades intro
INSERT INTO tutorial_sequences (tutorial_key, location, sequence)
VALUES ('gameplay_introduction', 'STG0001', 14);

INSERT INTO tutorial_sequence_translations (tutorial_sequence_id, language_code, title, description)
VALUES
    (currval('tutorial_sequences_id_seq'), 'id', 'Cat Manager', 'Terus kamu juga bisa beli ''Upgrade'' buat resto kamu.'),
    (currval('tutorial_sequences_id_seq'), 'en', 'Cat Manager', 'You can also buy ''Upgrades'' for your restaurant.');

-- Sequence 15: Highlight upgrade icon
INSERT INTO tutorial_sequences (tutorial_key, location, sequence)
VALUES ('gameplay_introduction', 'STG0001', 15);

INSERT INTO tutorial_sequence_translations (tutorial_sequence_id, language_code, title, description)
VALUES
    (currval('tutorial_sequences_id_seq'), 'id', 'System', 'highlight_upgrade_icon'),
    (currval('tutorial_sequences_id_seq'), 'en', 'System', 'highlight_upgrade_icon');

-- Sequence 16: Tap upgrade icon (input action)
INSERT INTO tutorial_sequences (tutorial_key, location, sequence)
VALUES ('gameplay_introduction', 'STG0001', 16);

INSERT INTO tutorial_sequence_translations (tutorial_sequence_id, language_code, title, description)
VALUES
    (currval('tutorial_sequences_id_seq'), 'id', 'Player', 'tap_upgrade_icon'),
    (currval('tutorial_sequences_id_seq'), 'en', 'Player', 'tap_upgrade_icon');

-- Sequence 17: Upgrade benefits
INSERT INTO tutorial_sequences (tutorial_key, location, sequence)
VALUES ('gameplay_introduction', 'STG0001', 17);

INSERT INTO tutorial_sequence_translations (tutorial_sequence_id, language_code, title, description)
VALUES
    (currval('tutorial_sequences_id_seq'), 'id', 'Cat Manager', 'Upgrades bakal bantu resto kamu biar dapur masaknya bisa lebih cepat.'),
    (currval('tutorial_sequences_id_seq'), 'en', 'Cat Manager', 'Upgrades will help your restaurant so the kitchen can cook faster.');

-- Sequence 18: Upgrade screen closed
INSERT INTO tutorial_sequences (tutorial_key, location, sequence)
VALUES ('gameplay_introduction', 'STG0001', 18);

INSERT INTO tutorial_sequence_translations (tutorial_sequence_id, language_code, title, description)
VALUES
    (currval('tutorial_sequences_id_seq'), 'id', 'System', 'close_upgrade_screen'),
    (currval('tutorial_sequences_id_seq'), 'en', 'System', 'close_upgrade_screen');

-- Sequence 19: Max level info
INSERT INTO tutorial_sequences (tutorial_key, location, sequence)
VALUES ('gameplay_introduction', 'STG0001', 19);

INSERT INTO tutorial_sequence_translations (tutorial_sequence_id, language_code, title, description)
VALUES
    (currval('tutorial_sequences_id_seq'), 'id', 'Cat Manager', 'Kalau dapur kamu udah level maksimal semua, kita bisa pindah ke resto yang lebih gede deh.'),
    (currval('tutorial_sequences_id_seq'), 'en', 'Cat Manager', 'If all your kitchens are at max level, we can move to a bigger restaurant.');

-- Sequence 20: Have fun
INSERT INTO tutorial_sequences (tutorial_key, location, sequence)
VALUES ('gameplay_introduction', 'STG0001', 20);

INSERT INTO tutorial_sequence_translations (tutorial_sequence_id, language_code, title, description)
VALUES
    (currval('tutorial_sequences_id_seq'), 'id', 'Cat Manager', 'Ayo kita bersenang-senang di resto ini!'),
    (currval('tutorial_sequences_id_seq'), 'en', 'Cat Manager', 'Let''s have fun at this restaurant!');

-- Sequence 21: Tutorial ends
INSERT INTO tutorial_sequences (tutorial_key, location, sequence)
VALUES ('gameplay_introduction', 'STG0001', 21);

INSERT INTO tutorial_sequence_translations (tutorial_sequence_id, language_code, title, description)
VALUES
    (currval('tutorial_sequences_id_seq'), 'id', 'System', 'tutorial_end'),
    (currval('tutorial_sequences_id_seq'), 'en', 'System', 'tutorial_end');
