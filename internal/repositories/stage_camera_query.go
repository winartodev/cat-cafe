package repositories

const (
	insertIntoStageCameraConfigQuery = `
		INSERT INTO stage_camera_configs (
		    stage_id,
			zoom_size,
			min_bound_x,
			min_bound_y, 
			max_bound_x, 
			max_bound_y, 
			created_at,
			updated_at
		) VALUES ($1, $2, $3, $4, $5, $6, $7, $8)
		RETURNING id
	`

	updateStageCameraConfigQuery = `
		UPDATE stage_camera_configs
		SET
			zoom_size = $2,
			min_bound_x = $3,
			min_bound_y = $4, 
			max_bound_x = $5, 
			max_bound_y = $6, 
			updated_at = $7
		WHERE stage_id = $1;
	`
)
