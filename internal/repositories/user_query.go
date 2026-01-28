package repositories

const (
	insertUserQuery = `
		INSERT INTO users 
		    (
		     external_id,
		     username,
		     email,
		     created_at,
		     updated_at
		     ) 
		VALUES (
		        $1, $2, $3, $4, $5
		) RETURNING id
	`

	// TODO: FIX THIS QUERY IMMEDIATELY
	getUserByIDQuery = `
		SELECT 
		    id, external_id, username, email, gem, coin
		FROM users WHERE id = $1
	`

	getUserByIDForUpdateQuery = `
		SELECT 
		    id, external_id, username, email, gem, coin
		FROM users WHERE id = $1 FOR UPDATE
	`

	getUserByEmailQuery = `
		SELECT 
		    id, external_id, username, email, gem, coin
		FROM users WHERE email = $1
	`

	getUserDailyRewardProgressQuery = `
		SELECT 
		    id,
		    user_id, 
		    longest_streak,
		    current_streak,
		    last_claim_date 
		FROM user_daily_reward_progress
		WHERE user_id = $1
	`

	getUserBalanceByIDQuery = `SELECT coin, gem FROM users WHERE id = $1`

	upsertUserDailyRewardProgressQuery = `
		INSERT INTO user_daily_reward_progress 
		    (
		     user_id,
		     longest_streak,
		     current_streak,
		     last_claim_date,
		     updated_at
		    ) 
		VALUES ($1, $2, $3, $4, $5)
		ON CONFLICT  (user_id)
		DO UPDATE SET 
		    longest_streak = EXCLUDED.longest_streak,
		    current_streak = EXCLUDED.current_streak,
		    last_claim_date = EXCLUDED.last_claim_date,
			updated_at=EXCLUDED.updated_at
	`

	updateLastSyncBalanceQuery = `UPDATE users SET last_sync_balance_at = $1, updated_at = $2 WHERE id = $3`
)
