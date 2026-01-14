package repositories

const (
	// TODO: FIX THIS QUERY IMMEDIATELY
	getUserByIDDB = `
		SELECT 
		    id, external_id, username
		FROM users WHERE id = $1
	`

	getUserDailyRewardQuery = `
		SELECT 
		    id,
		    user_id, 
		    current_streak,
		    last_claim_date 
		FROM user_daily_rewards
		WHERE user_id = $1
	`

	getUserBalanceByIDQuery = `SELECT coin, gem FROM users WHERE id = $1`

	upsertUserDailyRewardQuery = `
		INSERT INTO user_daily_rewards 
		    (
		     user_id,
		     current_streak,
		     last_claim_date,
		     updated_at
		    ) 
		VALUES ($1, $2, $3, $4)
		ON CONFLICT  (user_id)
		DO UPDATE SET 
		    current_streak = EXCLUDED.current_streak,
		    last_claim_date = EXCLUDED.last_claim_date,
			updated_at=EXCLUDED.updated_at
	`
)
