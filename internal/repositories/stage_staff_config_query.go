package repositories

const (
	insertStageStaffConfigQuery = `
		INSERT INTO stage_staff_configs (
			stage_id,
			starting_staff_manager,
			starting_staff_helper,
			created_at,
			updated_at
		) VALUES ($1, $2, $3, $4, $5) 
		  RETURNING id
	`

	updateStageStaffConfigQuery = `
		UPDATE stage_staff_configs
		SET
			starting_staff_manager = $2,
			starting_staff_helper = $3,
			updated_at = $4
		WHERE stage_id = $1;
	`
)
