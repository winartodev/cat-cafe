package repositories

const (
	insertFoodItemQuery = `
		INSERT INTO food_items(
		 	slug, 
			name, 
		   	initial_cost,
			initial_profit, 
			cooking_time, 
			created_at,
			updated_at              
		) VALUES ($1, $2, $3, $4, $5, $6,$7)
		RETURNING id;
	`

	updateFoodItemQuery = `
		UPDATE food_items 
		SET 
			name = $1, 
			initial_profit = $2, 
			cooking_time = $3,
			initial_cost = $4,
			updated_at = $5
		WHERE id = $6
		RETURNING id;
	`

	getFoodsQuery = `
		SELECT
			id,
			slug,
			name,
			initial_cost,
			initial_profit,
			cooking_time
		FROM food_items
		ORDER BY id
		LIMIT $1 OFFSET $2;
	`

	getFoodsBySlugQuery = `
		SELECT 
			id,
			slug,
			name,
			initial_cost,
			initial_profit,
			cooking_time
		FROM food_items
		WHERE slug = $1;
	`

	getFoodsByIDQuery = `
		SELECT 
			id,
			slug,
			name,
			initial_cost,
			initial_profit,
			cooking_time
		FROM food_items
		WHERE id = $1;
	`

	countFoodItemsQuery = `
		SELECT COUNT(*) 
		FROM food_items
	`

	insertFoodItemOverrideLevelQuery = `
		INSERT INTO food_level_overrides(
			food_item_id,
			level,
			cost,
			profit,
			preparation_time,
			created_at,
			updated_at
		) VALUES 
	`

	deleteFoodItemOverrideLevelQuery = `
		DELETE FROM food_level_overrides
		WHERE food_item_id = $1
	`

	getOverrideLevelQuery = `
		SELECT
			food_item_id,
			level,
			cost,
			profit,
			preparation_time
		FROM food_level_overrides
		WHERE food_item_id = $1;
	`

	getOverrideLevelByFoodItemIDAndLevelQuery = `
		SELECT
			food_item_id,
			level,
			cost,
			profit,
			preparation_time
		FROM food_level_overrides
		WHERE food_item_id = $1 AND level = $2;
	`
)
