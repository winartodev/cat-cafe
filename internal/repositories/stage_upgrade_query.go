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
)
