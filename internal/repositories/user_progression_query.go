package repositories

const (
	getGameStageProgressionQuery = `
		SELECT 
			id,
			user_id,
			stage_id,
			is_complete,
			completed_at
		FROM user_stage_progress
		WHERE user_id = $1;
	`

	insertGameStageProgressionQuery = `
		INSERT INTO user_stage_progress (
			user_id,
			stage_id                             
		) VALUES (
			$1, $2 
		) RETURNING id
	`
)
