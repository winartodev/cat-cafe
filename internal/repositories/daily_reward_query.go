package repositories

const (
	insertDailyRewardQuery = `
        INSERT INTO daily_rewards 
		(
			 reward_id,
			 day_number, 
			 is_active,
			 description,
			 created_at,
			 updated_at
		 )
        VALUES ($1, $2, $3, $4, $5, $6)
        RETURNING id
	`

	getDailyRewardsQuery = `
        SELECT
                dr.id,
                dr.reward_id,
                dr.day_number,
                r.slug,
                r.name,
                r.amount,
                r.is_active,
                dr.is_active,
                dr.description,
                rt.slug as reward_type_slug,
                rt.name as reward_type_name
        FROM daily_rewards AS dr
       		JOIN rewards AS r on r.id = dr.reward_id
        	JOIN reward_types AS rt on r.reward_type_id = rt.id
        order by dr.id
	`

	getDailyRewardsWithPaginationQuery = `
		SELECT
			dr.id,
			dr.reward_id,
			dr.day_number,
			r.slug,
			r.name,
			r.amount,
			r.is_active,
			dr.is_active,
			dr.description,
			rt.slug as reward_type_slug,
			rt.name as reward_type_name
		FROM daily_rewards AS dr
			JOIN rewards AS r ON r.id = dr.reward_id
			JOIN reward_types AS rt ON r.reward_type_id = rt.id
		ORDER BY dr.id
		LIMIT $1 OFFSET $2
	`

	countDailyRewardsQuery = `
		SELECT COUNT(*) 
		FROM daily_rewards
	`

	getDailyRewardByIDQuery = `
		SELECT
                dr.id,
                dr.reward_id,
                rt.slug,
                rt.name,
                dr.day_number,
                r.amount,
                dr.is_active,
                dr.description
		FROM daily_rewards AS dr
       		JOIN rewards AS r on r.id = dr.reward_id
        	JOIN reward_types AS rt on r.reward_type_id = rt.id 
		WHERE dr.id = $1 
	`

	updateDailyRewardQuery = `
		UPDATE daily_rewards
		SET
			reward_id = $1,
			day_number = $2,
			is_active = $3,
			description = $4,
			updated_at = $5
		WHERE id = $6
	`
)
