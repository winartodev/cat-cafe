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

	insertDailyRewardQuery = `
        INSERT INTO daily_rewards 
		(
			 reward_type_id,
			 day_number,
			 reward_amount, 
			 is_active,
			 description,
			 created_at,
			 updated_at
		 )
        VALUES ($1, $2, $3, $4, $5, $6, $7)
        RETURNING id`

	getDailyRewardsQuery = `
        SELECT
                dr.id,
                dr.reward_type_id,
                rt.slug,
                rt.name,
                dr.day_number,
                dr.reward_amount,
                dr.is_active,
                dr.description
        FROM daily_rewards AS dr
        JOIN reward_types AS rt on dr.reward_type_id = rt.id  order by dr.id
	`

	getDailyRewardByIDQuery = `
		SELECT
                dr.id,
                dr.reward_type_id,
                rt.slug,
                rt.name,
                dr.day_number,
                dr.reward_amount,
                dr.is_active,
                dr.description
        FROM daily_rewards AS dr
        JOIN reward_types AS rt on dr.reward_type_id = rt.id WHERE dr.id = $1 
	`

	getDailyRewardByDayQuery = `
		SELECT
               dr.id,
                dr.reward_type_id,
                rt.slug,
                rt.name,
                dr.day_number,
                dr.reward_amount,
                dr.is_active,
                dr.description
        FROM daily_rewards AS dr
        JOIN reward_types AS rt on dr.reward_type_id = rt.id  WHERE dr.day_number = $1 
	`

	updateDailyRewardQuery = `
		UPDATE daily_rewards
		SET
			reward_type_id = $1,
			day_number = $2,
			reward_amount = $3,
			is_active = $4,
			description = $5,
			updated_at = $6
		WHERE id = $7
	`
)
