package repositories

const (
	bulkInsertKitchenStationQuery = `
		INSERT INTO kitchen_stations (
			stage_id,
		    food_item_id, 
		    auto_unlock,
		    created_at, 
		    updated_at
		) VALUES 
	`

	getKitchenStationsQuery = `
		SELECT 
		    ks.stage_id,
		    ks.food_item_id,
		    ks.auto_unlock,
		    fi.slug,
		    fi.name,
		    fi.starting_price,
		    fi.starting_preparation
		FROM kitchen_stations AS ks
			JOIN food_items AS fi ON fi.id = ks.food_item_id
		WHERE ks.stage_id = $1
	`

	deleteKitchenStationQuery = `
		DELETE FROM kitchen_stations
		WHERE stage_id = $1;
	`

	getKitchenStationByFoodIDDB = `
		SELECT 
		    ks.stage_id,
		    ks.food_item_id,
		    ks.auto_unlock,
		    fi.slug,
		    fi.name,
		    fi.starting_price,
		    fi.starting_preparation
		FROM kitchen_stations AS ks
			JOIN food_items AS fi ON fi.id = ks.food_item_id
		WHERE ks.stage_id = $1 AND ks.food_item_id = $2
	`
)
