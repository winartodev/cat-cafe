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
		WHERE user_id = $1 and stage_id =$2
		ORDER BY created_at DESC, id DESC  
    	LIMIT 1;
	`

	getLatestGameStageProgressionQuery = `
		SELECT 
			id,
			user_id,
			stage_id,
			is_complete,
			completed_at
		FROM user_stage_progress
		WHERE user_id = $1
		ORDER BY id DESC
    	LIMIT 1;
	`

	insertGameStageProgressionQuery = `
		INSERT INTO user_stage_progress (
			user_id,
			stage_id                             
		) VALUES (
			$1, $2 
		) RETURNING id
	`

	checkStageProgressionExistsQuery = `
  		SELECT EXISTS(
            SELECT 1 
            FROM user_stage_progress 
            WHERE user_id = $1 AND stage_id = $2
		)
	`

	markStageAsComplete = `
		UPDATE user_stage_progress
		SET 
			is_complete = true,
			completed_at = $1
		WHERE user_id = $2 AND stage_id = $3
	`
)
