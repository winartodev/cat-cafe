package repositories

const (
	insertStageUpgradeQuery = `
		INSERT INTO game_stage_upgrades (
			game_stage_id,
			upgrade_id,
			created_at,
			updated_at
		) VALUES
	`

	getStageUpgradeQuery = `
		SELECT
			gs.slug,
			u.*
		FROM game_stage_upgrades gsu
		JOIN game_stages gs ON gsu.game_stage_id =  gs.id
		JOIN upgrades u  ON gsu.upgrade_id = u.id
		WHERE gsu.game_stage_id = $1
	`

	getStageUpgradeCountQuery = `
		SELECT COUNT(u.*)
		FROM game_stage_upgrades gsu
				 JOIN game_stages gs ON gsu.game_stage_id =  gs.id
				 JOIN upgrades u  ON gsu.upgrade_id = u.id
		WHERE gsu.game_stage_id = $1
	`

	deleteStageUpgradeQuery = `
		DELETE FROM game_stage_upgrades WHERE game_stage_id = $1
	`

	getUpgradeByStageIDAndSlugQuery = `
		SELECT 
		    gsu.id,
		    u.slug, u.name, u.cost,
    		u.cost_type, u.effect_value, u.effect_unit,
    		u.effect_type, u.effect_target, u.effect_target_id, 
    		COALESCE(fi.slug, '') AS target_name
		FROM 
		    game_stage_upgrades gsu
		JOIN upgrades u ON gsu.upgrade_id = u.id
		LEFT JOIN food_items fi ON fi.id = u.effect_target_id AND u.effect_target = 'food' 
		WHERE 
		    gsu.game_stage_id = $1 AND
		    u.slug = $2
	`
)
