package repositories

const (
	insertIntoGameStageQuery = `
		INSERT INTO game_stages (
		    slug,
			name,
			description,
			starting_coin,
			stage_prize,
			is_active,
			sequence,
			created_at,
			updated_at
		) VALUES ($1, $2, $3, $4, $5, $6, $7, $8,$9)
		RETURNING id;
	`

	updateGameStageQuery = `
		UPDATE game_stages
		SET
			name = $2,
			starting_coin = $3,
			stage_prize = $4,
			updated_at = $5,
			is_active = $6,
			sequence = $7,
			description = $8
		WHERE id = $1;
	`

	getGameStageByIDQuery = `
		SELECT 
		    id,
			slug,
			name,
			description,
			starting_coin,
			stage_prize,
			is_active,
			sequence
		FROM game_stages
		WHERE id = $1;
	`

	getGameStageBySlugQuery = `
		SELECT 
		    id,
			slug,
			name,
			description,
			starting_coin,
			stage_prize,
			is_active,
			sequence
		FROM game_stages
		WHERE slug = $1;
	`

	getGameStageQuery = `
		SELECT 
			id,
			slug,
			name,
			description,
			starting_coin,
			stage_prize,
			is_active,
			sequence
		FROM game_stages
		ORDER BY sequence, 
		         CASE WHEN is_active = true THEN 0 ELSE 1 END 
		LIMIT $1 OFFSET $2;
	`

	countGameStagesQuery = `
		SELECT COUNT(*) 
		FROM game_stages
	`

	getGameStageConfig = `
		SELECT COALESCE(scc.customer_spawn_time, 0)               	AS customer_spawn_time,
			   COALESCE(scc.max_customer_order_count, 0)          	AS max_customer_order_count,
			   COALESCE(scc.max_customer_order_variant, 0)        	AS max_customer_order_variant,
			   COALESCE(scc.starting_order_table_count, 0)        	AS starting_order_table_count,
			   COALESCE(ssc.starting_staff_manager, '')           	AS starting_staff_manager,
			   COALESCE(ssc.starting_staff_helper, '')            	AS starting_staff_helper,
			   COALESCE(skc.id, 0) 									AS kitchen_config_id,
			   COALESCE(skc.max_level, 0)                         	AS max_level,
			   COALESCE(skc.upgrade_profit_multiply, 0)           	AS upgrade_profit_multiply,
			   COALESCE(skc.upgrade_cost_multiply, 0)             	AS upgrade_cost_multiply,
			   COALESCE(skc.transition_phase_levels, '{}')        	AS transition_phase_levels,
			   COALESCE(skc.phase_profit_multipliers, '{}')       	AS phase_profit_multipliers,
			   COALESCE(skc.phase_upgrade_cost_multipliers, '{}') 	AS phase_upgrade_cost_multipliers,
			   COALESCE(skc.table_count_per_phases, '{}')         	AS table_count_per_phases,
			   COALESCE(rewards_agg.data, '[]')                   	AS phase_rewards,
			   COALESCE(stations_agg.data, '[]')                  	AS kitchen_stations,
			   COALESCE(scc2.zoom_size, 0.0)                      	AS zoom_size,
			   COALESCE(scc2.max_bound_x, 0.0)                    	AS max_bound_x,
			   COALESCE(scc2.min_bound_x, 0.0)                    	AS min_bound_x,
			   COALESCE(scc2.min_bound_y, 0.0)                    	AS min_bound_y,
			   COALESCE(scc2.max_bound_y, 0.0)                    	AS max_bound_y
		FROM game_stages gs
				 LEFT JOIN stage_customer_configs scc ON scc.stage_id = gs.id
				 LEFT JOIN stage_staff_configs ssc ON ssc.stage_id = gs.id
				 LEFT JOIN stage_kitchen_configs skc ON skc.stage_id = gs.id
				 LEFT JOIN stage_camera_configs scc2 ON scc2.stage_id = gs.id
				 LEFT JOIN LATERAL (
			SELECT json_agg(json_build_object(
									'phase_number', cr.phase_number,
									'reward_id', cr.reward_id,
									'reward_slug', r.slug,
									'reward_type', rt.slug
							) ORDER BY cr.phase_number) as data
			FROM kitchen_phase_completion_rewards AS cr
					 JOIN rewards r ON r.id = cr.reward_id
					 JOIN reward_types rt ON rt.id = r.reward_type_id
			WHERE cr.kitchen_config_id = skc.id
			) rewards_agg ON TRUE
				 LEFT JOIN LATERAL (
			SELECT json_agg(json_build_object(
					'food_item_slug', fi.slug,
					'food_name', fi.name,
					'starting_price', fi.starting_price,
					'starting_preparation', fi.starting_preparation,
					'auto_unlock', ks.auto_unlock
							)) as data
			FROM kitchen_stations AS ks
					 JOIN food_items fi ON fi.id = ks.food_item_id
			WHERE ks.stage_id = gs.id
			) stations_agg ON TRUE
		WHERE gs.id = $1;
	`

	getActiveGameStagesQuery = `
		SELECT
		    id,
			slug,
			name,
			description,
			sequence
		FROM game_stages
		WHERE is_active = true
		ORDER BY sequence
	`
)
