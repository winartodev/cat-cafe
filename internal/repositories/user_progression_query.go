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
			completed_at,
			last_started_at
		FROM user_stage_progress
		WHERE user_id = $1
		ORDER BY last_started_at DESC
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
			unlocked_stations,
			station_upgrades
		FROM user_kitchen_progress
		WHERE user_id = $1 and stage_id =$2
	`

	insertUserKitchenProgressQuery = `
	 	INSERT INTO user_kitchen_progress (
        	user_id, 
	 	    stage_id,
	 	    station_levels,
	 	    unlocked_stations,
	 		station_upgrades
	 	) VALUES ($1, $2, $3, $4, $5)
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

	markStageAsStartedQuery = `
		UPDATE user_stage_progress 
			SET 
			    last_started_at = $1
		WHERE user_id = $2 AND stage_id = $3
	`

	insertUserUpgradeStageProgressionQuery = `
		INSERT INTO user_stage_upgrades (
		user_id,
		game_stage_id,
		game_stage_upgrade_id,
		purchased_at,
		created_at,
		updated_at
		) VALUES (
		        $1, $2, $3, $4, $5, $6  
		)
	`

	updateStationUpgradeQuery = `
		UPDATE user_kitchen_progress  
			SET station_upgrades = $1,
			    updated_at = $2
		WHERE user_id = $3 AND stage_id = $4
	`

	stageUpgradeAlreadyPurchaseQuery = `
		SELECT
		    usu.id,
			u.slug, u.name, u.cost, 
			u.cost_type, u.effect_type, u.effect_value, 
			u.effect_unit, u.effect_target, u.effect_target_id, 
			COALESCE(fi.slug, '') AS effect_target_name
		FROM user_stage_upgrades usu
		JOIN game_stage_upgrades gsu on usu.game_stage_upgrade_id = gsu.id
		JOIN upgrades u ON gsu.upgrade_id = u.id
			LEFT JOIN food_items fi ON fi.id = u.effect_target_id AND u.effect_target = 'food'
		where usu.user_id = $1 AND usu.game_stage_id = $2
	`

	currentStageUpgradeQuery = `
		SELECT
			u.slug, u.name,  u.description, u.cost, 
			u.cost_type, u.effect_type, u.effect_value, 
			u.effect_unit, u.effect_target, u.effect_target_id,
			COALESCE(fi.slug, '') AS effect_target_name,
			CASE
				WHEN usp.id IS NULL THEN false
				ELSE true
			END AS is_purchased,
			usp.purchased_at
		FROM game_stage_upgrades gsu
		JOIN upgrades u ON gsu.upgrade_id = u.id
		    LEFT JOIN food_items fi ON fi.id = u.effect_target_id AND u.effect_target = 'food'
		LEFT JOIN user_stage_upgrades usp ON usp.game_stage_upgrade_id = gsu.id                        
			AND usp.user_id = $1
		WHERE gsu.game_stage_id = $2
		ORDER BY is_purchased, u.sequence;
	`
)
