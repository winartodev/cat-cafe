package repositories

import (
	"context"
	"database/sql"
	"errors"

	"github.com/winartodev/cat-cafe/internal/entities"
	"github.com/winartodev/cat-cafe/pkg/apperror"
	"github.com/winartodev/cat-cafe/pkg/database"
	"github.com/winartodev/cat-cafe/pkg/helper"
)

type KitchenStationRepository interface {
	WithTX(tx *sql.Tx) KitchenStationRepository

	CreateKitchenStationsWithTxDB(ctx context.Context, stageID int64, items []entities.KitchenStation) ([]int64, error)
	GetKitchenStationsDB(ctx context.Context, stageID int64) (res []entities.KitchenStation, err error)
	DeleteKitchenStationDB(ctx context.Context, stageID int64) error
	GetKitchenStationByFoodIDDB(ctx context.Context, stageID int64, foodItemID int64) (res *entities.KitchenStation, err error)
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
	queryString := r.BuildBulkInsertQuery(bulkInsertKitchenStationQuery, len(items), numFields, "RETURNING id")

	args := make([]interface{}, 0, len(items)*numFields)
	now := helper.NowUTC()

	for _, item := range items {
		args = append(args,
			item.StageID,
			item.FoodItemID,
			item.AutoUnlock,
			now,
			now,
		)
	}

	rows, err := r.db.QueryContext(ctx, queryString, args...)
	if database.IsDuplicateError(err) {
		return nil, apperror.ErrConflict
	} else if err != nil {
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
			&foodItem.InitialProfit,
			&foodItem.CookingTime,
			&foodItem.InitialCost,
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

func (r *kitchenStationRepository) GetKitchenStationByFoodIDDB(ctx context.Context, stageID int64, foodItemID int64) (res *entities.KitchenStation, err error) {
	var kitchenStation entities.KitchenStation
	var foodItem entities.FoodItem

	err = r.db.QueryRowContext(
		ctx,
		getKitchenStationByFoodIDDB,
		stageID,
		foodItemID,
	).Scan(
		&kitchenStation.StageID,
		&kitchenStation.FoodItemID,
		&kitchenStation.AutoUnlock,
		&foodItem.Slug,
		&foodItem.Name,
		&foodItem.InitialProfit,
		&foodItem.CookingTime,
		&foodItem.InitialCost,
	)
	if errors.Is(err, sql.ErrNoRows) {
		return nil, nil
	} else if err != nil {
		return nil, err
	}

	kitchenStation.FoodItem = &foodItem

	return &kitchenStation, nil
}
