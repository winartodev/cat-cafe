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

	getUserKitchenProgressQuery = `
		SELECT 
			id,
			user_id,
			stage_id,
			station_levels,
			unlocked_stations
		FROM user_kitchen_progress
		WHERE user_id = $1 and stage_id =$2
	`

	insertUserKitchenProgressQuery = `
	 	INSERT INTO user_kitchen_progress (
        	user_id, 
	 	    stage_id,
	 	    station_levels,
	 	    unlocked_stations
	 	) VALUES ($1, $2, $3, $4)
        ON CONFLICT (user_id, stage_id) DO NOTHING
	`

	updateUserKitchenProgressQuery = `
		UPDATE user_kitchen_progress
        SET 
            station_levels = $3,
            unlocked_stations = $4,
            updated_at = $5
        WHERE user_id = $1 AND stage_id = $2
	`

	insertUserKitchenPhaseProgressQuery = `
		INSERT INTO user_kitchen_phase_progress (
			user_id,
		    kitchen_config_id, 
		    current_phase, 
		    completed_phases
		) VALUES ($1, $2, $3, $4)
        ON CONFLICT (user_id, kitchen_config_id) DO NOTHING
	`

	updateUserKitchenPhaseProgressQuery = `
		UPDATE user_kitchen_phase_progress
        SET 
            current_phase = $3,
            completed_phases = $4,
            updated_at = $5
        WHERE user_id = $1 AND kitchen_config_id = $2
	`

	getUserKitchenPhaseProgressQuery = `
		SELECT 
			user_id,
		    kitchen_config_id, 
		    current_phase, 
		    completed_phases
		FROM user_kitchen_phase_progress
		WHERE user_id = $1 and kitchen_config_id =$2
	`

	getUserPhaseRewardClaimedQuery = `
		SELECT 1 as claimed
		FROM user_kitchen_phase_reward_claims
		WHERE 
		    user_id = $1 AND 
		    kitchen_config_id = $2 AND 
		    phase_number = $3 AND 
		    reward_id = $4 
		LIMIT 1
	`

	createUserKitchenPhaseClaimReward = `
		INSERT INTO user_kitchen_phase_reward_claims
		(
		 	user_id,
		 	kitchen_config_id,
		 	reward_id,
		 	phase_number,
		 	claimed_at
		) VALUES ($1, $2, $3, $4,$5)
	`
)
