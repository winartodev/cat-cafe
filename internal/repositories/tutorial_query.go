package repositories

const (
	createTutorialSequenceQuery = `
		INSERT INTO tutorial_sequences (
			tutorial_key,
		    location,
		    sequence,
		    created_at,
		    updated_at
		) VALUES (
			$1, $2, $3, $4, $5
		) RETURNING id
	`

	createTutorialTranslations = `
		INSERT INTO tutorial_sequence_translations (
			tutorial_sequence_id,
		    language_code,
		    title,
		    description,
		    created_at,
		    updated_at
		) VALUES 
	`

	getTutorialsQuery = `
		SELECT
			tutorial_key,
			location
		FROM tutorial_sequences
		GROUP BY (tutorial_key, location)
		LIMIT $1 OFFSET $2;
	`

	getDetailTutorialsQuery = `
		SELECT
		    ts.id,
			ts.tutorial_key,
			ts.location,
			ts.sequence,
			tst.language_code,
			tst.title,
			tst.description
		FROM tutorial_sequences ts
		JOIN tutorial_sequence_translations tst ON ts.id = tst.tutorial_sequence_id
		WHERE ts.tutorial_key = $1
		LIMIT $2 OFFSET $3;
	`

	updateTutorialSequenceQuery = `
        UPDATE tutorial_sequences 
        SET location = $1, 
            sequence = $2, 
            updated_at = $3 
        WHERE id = $4
    `

	deleteTutorialTranslationsQuery = `
        DELETE FROM tutorial_sequence_translations 
        WHERE tutorial_sequence_id = $1
    `
)
