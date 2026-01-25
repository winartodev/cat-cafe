package repositories

const (
	rewardTypeInsertQuery = `
        INSERT INTO reward_types (slug, name, created_at, updated_at) 
        VALUES ($1, $2, $3, $4) 
        RETURNING id`

	getRewardTypesQuery = `
        SELECT id, slug, name FROM reward_types ORDER BY id
        `

	getRewardTypeByIDQuery = `
        SELECT id, slug, name FROM reward_types WHERE id=$1
        `

	getRewardTypeBySlugQuery = `
        SELECT id, slug, name FROM reward_types WHERE slug=$1
        `

	updateRewardTypeQuery = `
        UPDATE reward_types SET name = $1, updated_at = $2 WHERE id = $3
        `

	rewardInsertQuery = `
		INSERT INTO rewards ( 
		    reward_type_id,
		    slug,
		    name,
		    amount,
			is_active,
		    created_at,
		    updated_at
		) VALUES ($1, $2, $3, $4, $5, $6, $7)
		RETURNING id
	`

	getRewardQuery = `
		SELECT 
		    r.id, 
		    r.slug, 
		    r.name,
		    r.amount,
		    r.is_active,
		    rt.slug reward_type_slug,
		    rt.name reward_type_name
		FROM rewards AS r 
			JOIN reward_types AS rt ON r.reward_type_id = rt.id 
		ORDER BY r.id
		LIMIT $1 OFFSET $2
	`

	getRewardBySlugQuery = `
		SELECT 
		    r.id, 
		    r.slug, 
		    r.name,
		    r.amount,
		    r.is_active,
		    rt.slug reward_type_slug,
		    rt.name reward_type_name
		FROM rewards AS r 
			JOIN reward_types AS rt ON r.reward_type_id = rt.id 
		WHERE r.slug = $1
	`

	getRewardByIDQuery = `
		SELECT 
		    r.id, 
		    r.slug, 
		    r.name,
		    r.amount,
		    r.is_active,
		    rt.slug reward_type_slug,
		    rt.name reward_type_name
		FROM rewards AS r 
			JOIN reward_types AS rt ON r.reward_type_id = rt.id 
		WHERE r.id = $1
	`

	countRewardsQuery = `
		SELECT COUNT(*) 
		FROM rewards
	`
)
