BEGIN;

-- Delete all tutorial sequences and their translations
-- (translations will be deleted automatically due to foreign key)
DELETE FROM tutorial_sequences
WHERE tutorial_key IN ('restaurant_introduction', 'gameplay_introduction');

COMMIT;
