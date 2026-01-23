package repositories

import (
	"context"
	"database/sql"
	"fmt"
	"github.com/winartodev/cat-cafe/internal/entities"
	"github.com/winartodev/cat-cafe/pkg/helper"
	"strings"
)

type KitchenStationRepository interface {
	WithTX(tx *sql.Tx) KitchenStationRepository

	CreateKitchenStationsWithTxDB(ctx context.Context, stageID int64, items []entities.KitchenStation) ([]int64, error)
	GetKitchenStationsDB(ctx context.Context, stageID int64) (res []entities.KitchenStation, err error)
	DeleteKitchenStationDB(ctx context.Context, stageID int64) error
}

type kitchenStationRepository struct {
	BaseRepository
}

func NewKitchenStationRepository(db *sql.DB) KitchenStationRepository {
	return &kitchenStationRepository{
		BaseRepository{db: db, pool: db},
	}
}

func (r *kitchenStationRepository) WithTX(tx *sql.Tx) KitchenStationRepository {
	if tx == nil {
		return r
	}

	return &kitchenStationRepository{
		BaseRepository{
			db:   tx,
			pool: r.pool,
		},
	}
}

func (r *kitchenStationRepository) CreateKitchenStationsWithTxDB(ctx context.Context, stageID int64, items []entities.KitchenStation) ([]int64, error) {
	if len(items) == 0 {
		return nil, nil
	}

	numFields := 5
	placeholderCount := 1
	queryBuilder := strings.Builder{}
	queryBuilder.WriteString(bulkInsertKitchenStationQuery)

	args := make([]interface{}, 0, len(items)*numFields)
	for i, item := range items {
		if i > 0 {
			queryBuilder.WriteString(", ")
		}

		queryBuilder.WriteString("(")
		for j := 0; j < numFields; j++ {
			queryBuilder.WriteString(fmt.Sprintf("$%d", placeholderCount))
			if j < numFields-1 {
				queryBuilder.WriteString(", ")
			}
			placeholderCount++
		}
		queryBuilder.WriteString(")")

		now := helper.NowUTC()
		args = append(args, item.StageID, item.FoodItemID, item.AutoUnlock, now, now)
	}

	queryBuilder.WriteString(" RETURNING id;")

	rows, err := r.db.QueryContext(ctx, queryBuilder.String(), args...)
	if err != nil {
		return nil, err
	}

	defer rows.Close()

	var ids []int64
	for rows.Next() {
		var id int64
		if err := rows.Scan(&id); err != nil {
			return nil, err
		}
		ids = append(ids, id)
	}

	return ids, nil
}

func (r *kitchenStationRepository) GetKitchenStationsDB(ctx context.Context, stageID int64) (res []entities.KitchenStation, err error) {
	rows, err := r.db.QueryContext(ctx, getKitchenStationsQuery, stageID)
	if err != nil {
		return nil, err
	}

	defer rows.Close()

	for rows.Next() {
		var kitchenStation entities.KitchenStation
		var foodItem entities.FoodItem

		err := rows.Scan(
			&kitchenStation.StageID,
			&kitchenStation.FoodItemID,
			&kitchenStation.AutoUnlock,
			&foodItem.Slug,
			&foodItem.Name,
			&foodItem.StartingPrice,
			&foodItem.StartingPreparation,
		)
		if err != nil {
			return nil, err
		}

		kitchenStation.FoodItem = &foodItem
		res = append(res, kitchenStation)
	}

	return res, nil
}

func (r *kitchenStationRepository) DeleteKitchenStationDB(ctx context.Context, stageID int64) error {
	_, err := r.db.ExecContext(ctx, deleteKitchenStationQuery, stageID)
	if err != nil {
		return err
	}

	return nil
}
