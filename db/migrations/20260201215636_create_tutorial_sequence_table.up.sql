BEGIN;

CREATE TABLE tutorial_sequences (
    id SERIAL PRIMARY KEY,
    tutorial_key VARCHAR(50) NOT NULL,
    location VARCHAR(100) NOT NULL,
    sequence INT DEFAULT 0 NOT NULL,
    created_at TIMESTAMPTZ DEFAULT NOW() NOT NULL,
    updated_at timestamptz DEFAULT NOW() NOT NULL,
    UNIQUE(tutorial_key, sequence)
);

CREATE TABLE tutorial_sequence_translations (
    id SERIAL PRIMARY KEY,
    tutorial_sequence_id BIGINT REFERENCES tutorial_sequences(id) NOT NULL,
    language_code VARCHAR(5) NOT NULL,
    title VARCHAR(255) NOT NULL,
    description TEXT,
    UNIQUE(tutorial_sequence_id, language_code)
);

CREATE INDEX idx_tutorial_lang ON tutorial_sequence_translations(language_code);

COMMIT;
