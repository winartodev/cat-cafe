package repositories

const (
	insertIntoCustomerConfig = `
		INSERT INTO stage_customer_configs (
			stage_id,
			customer_spawn_time,
			max_customer_order_count,
			max_customer_order_variant,
			starting_order_table_count,
			created_at,
			updated_at
		) VALUES ($1, $2, $3, $4, $5, $6, $7)
		RETURNING id
	`

	updateCustomerConfigQuery = `
		UPDATE stage_customer_configs
		SET
			customer_spawn_time = $2,
			max_customer_order_count = $3,
			max_customer_order_variant = $4,
			starting_order_table_count = $5,
			updated_at = $6
		WHERE stage_id = $1;
	`
)
