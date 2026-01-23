package repositories

const (
	insertFoodItemQuery = `
		INSERT INTO food_items(
		 	slug, 
			name, 
			starting_price, 
			starting_preparation, 
			created_at,
			updated_at              
		) VALUES ($1, $2, $3, $4, $5, $6)
		RETURNING id;
	`

	updateFoodItemQuery = `
		UPDATE food_items 
		SET 
			name = $1, 
			starting_price = $2, 
			starting_preparation = $3,
			updated_at = $4
		WHERE id = $5
		RETURNING id;
	`

	getFoodsQuery = `
		SELECT
			id,
			slug,
			name,
			starting_price,
			starting_preparation
		FROM food_items
		ORDER BY id
		LIMIT $1 OFFSET $2;
	`

	getFoodsBySlugQuery = `
		SELECT 
			id,
			slug,
			name,
			starting_price,
			starting_preparation
		FROM food_items
		WHERE slug = $1;
	`

	getFoodsByIDQuery = `
		SELECT 
			id,
			slug,
			name,
			starting_price,
			starting_preparation
		FROM food_items
		WHERE id = $1;
	`

	countFoodItemsQuery = `
		SELECT COUNT(*) 
		FROM food_items
	`
)
