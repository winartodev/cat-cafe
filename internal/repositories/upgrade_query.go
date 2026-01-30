package repositories

const (
	insertUpgradeQuery = `
		INSERT INTO upgrades (
			slug,
			name,
			description,
			cost,
			cost_type,
			effect_type,
			effect_value,
			effect_unit,
			effect_target,
			effect_target_id,
			is_active,
			sequence,
			created_at,
			updated_at
		) VALUES (
			$1,
			$2,
			$3,
			$4,
			$5,
			$6,
			$7,
			$8,
			$9,
			$10,
			$11,
			$12,
			$13,
			$14
		) RETURNING id`

	getUpgradesQuery = `
		SELECT
			id,
			slug,
			name,
			description,
			is_active,
			sequence
		FROM upgrades
		ORDER BY sequence ASC
		LIMIT $1 OFFSET $2
	`

	countUpgradesQuery = `
		SELECT 
			COUNT(*) 
		FROM upgrades 
		WHERE is_active = true
	`

	getUpgradeByIDQuery = `
		SELECT
			u.id,
			u.slug,
			u.name,
			u.description,
			u.is_active,
			u.sequence,
			u.cost,
			u.cost_type,
			u.effect_type,
			u.effect_value,
			u.effect_unit,
			u.effect_target,
			u.effect_target_id,
			COALESCE(fi.slug, '') as effect_target_name
		FROM upgrades u	
		LEFT JOIN food_items fi ON fi.id = u.effect_target_id AND u.effect_target = 'food'
		WHERE u.id = $1
	`

	getUpgradeBySlugQuery = `
		SELECT
			u.id,
			u.slug,
			u.name,
			u.description,
			u.is_active,
			u.sequence,
			u.cost,
			u.cost_type,
			u.effect_type,
			u.effect_value,
			u.effect_unit,
			u.effect_target,
			u.effect_target_id,
			COALESCE(fi.slug, '') as effect_target_name
		FROM upgrades u	
		LEFT JOIN food_items fi ON fi.id = u.effect_target_id AND u.effect_target = 'food'
		WHERE u.slug = $1 AND u.is_active = true
	`

	getActiveUpgradesQuery = `
		SELECT
			u.id,
			u.slug,
			u.name,
			u.description,
			u.is_active,
			u.sequence,
			u.cost,
			u.cost_type,
			u.effect_type,
			u.effect_value,
			u.effect_unit,
			u.effect_target,
			u.effect_target_id,
			COALESCE(fi.slug, '') as effect_target_name
		FROM upgrades u
		LEFT JOIN food_items fi ON fi.id = u.effect_target_id AND u.effect_target = 'food'
		WHERE u.is_active = true
		ORDER BY u.sequence ASC
	`

	updateUpgradeQuery = `
		UPDATE upgrades
		SET
			slug = $1,
			name = $2,
			description = $3,
			cost = $4,
			cost_type = $5,
			effect_type = $6,
			effect_value = $7,
			effect_unit = $8,
			effect_target = $9,
			effect_target_id = $10,
			is_active = $11,
			sequence = $12,
			updated_at = $13
		WHERE id = $14
	`
)
