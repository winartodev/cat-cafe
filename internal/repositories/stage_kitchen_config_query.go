package repositories

const (
	insertStageKitchenConfigQuery = `
		INSERT INTO stage_kitchen_configs (
			stage_id,
		   	max_level,
		   	upgrade_profit_multiply,
		   	upgrade_cost_multiply,
		   	transition_phase_levels,
		   	phase_profit_multipliers,
		   	phase_upgrade_cost_multipliers,
		   	table_count_per_phases,
		   	created_at,
		   	updated_at
		) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10)
		RETURNING id;
	`

	updateStageKitchenConfigQuery = `
		UPDATE stage_kitchen_configs
		SET
		   	max_level = $2,
		   	upgrade_profit_multiply = $3,
		   	upgrade_cost_multiply = $4,
		   	transition_phase_levels = $5,
		   	phase_profit_multipliers = $6,
		   	phase_upgrade_cost_multipliers = $7,
		   	table_count_per_phases = $8,
		   	updated_at = $9
		WHERE stage_id = $1
		RETURNING id;
	`

	insertKitchenPhaseCompletionRewardsQuery = `
		INSERT INTO kitchen_phase_completion_rewards (
			kitchen_config_id,
			phase_number,
			reward_id,
			created_at,
			updated_at
		) VALUES ($1, $2, $3, $4, $5)
		RETURNING id;
	`

	deleteKitchenPhaseCompletionRewardsQuery = `
		DELETE FROM kitchen_phase_completion_rewards
		WHERE kitchen_config_id = $1;
	`
)
