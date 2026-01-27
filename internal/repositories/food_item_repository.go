package repositories

import (
	"context"
	"database/sql"
	"errors"

	"github.com/lib/pq"

	"github.com/winartodev/cat-cafe/internal/entities"
	"github.com/winartodev/cat-cafe/pkg/apperror"
	"github.com/winartodev/cat-cafe/pkg/database"
	"github.com/winartodev/cat-cafe/pkg/helper"
)

type FoodItemRepository interface {
	WithTx(tx *sql.Tx) FoodItemRepository
	FoodItemWithTx(ctx context.Context, fn func(tx *sql.Tx) error) error

	CreateFoodDB(ctx context.Context, data entities.FoodItem) (id *int64, err error)
	UpdateFoodDB(ctx context.Context, id int64, data entities.FoodItem) (err error)
	GetFoodBySlugDB(ctx context.Context, slug string) (*entities.FoodItem, error)
	GetFoodByIDDB(ctx context.Context, id int64) (*entities.FoodItem, error)
	GetFoodsDB(ctx context.Context, limit, offset int) ([]entities.FoodItem, error)
	CountFoodItemDB(ctx context.Context) (count int64, err error)

	GetFoodItemIDsBySlugsDB(ctx context.Context, slugs []string) (map[string]int64, error)

	CreateOverrideLevelDB(ctx context.Context, foodItemID int64, data []entities.FoodItemOverrideLevel) (err error)
	GetOverrideLevelDB(ctx context.Context, foodItemID int64) ([]entities.FoodItemOverrideLevel, error)
	GetOverrideLevelByFoodItemIDAndLevelDB(ctx context.Context, foodItemID int64, level int) (*entities.FoodItemOverrideLevel, error)
	DeleteOverrideLevelDB(ctx context.Context, foodItemID int64) (err error)
}

type foodItemRepository struct {
	BaseRepository
}

func NewFoodItemRepository(db *sql.DB) FoodItemRepository {
	return &foodItemRepository{
		BaseRepository{
			db:   db,
			pool: db,
		},
	}
}

func (r *foodItemRepository) WithTx(tx *sql.Tx) FoodItemRepository {
	if tx == nil {
		return r
	}

	return &foodItemRepository{
		BaseRepository{
			db:   tx,
			pool: r.pool,
		},
	}
}

func (r *foodItemRepository) FoodItemWithTx(ctx context.Context, fn func(tx *sql.Tx) error) error {
	tx, err := r.pool.BeginTx(ctx, nil)
	if err != nil {
		return err
	}

	defer tx.Rollback()

	if err := fn(tx); err != nil {
		return err
	}

	return tx.Commit()
}

func (r *foodItemRepository) CreateFoodDB(ctx context.Context, data entities.FoodItem) (id *int64, err error) {
	now := helper.NowUTC()
	var lastInsertId int64

	err = r.db.QueryRowContext(ctx,
		insertFoodItemQuery,
		data.Slug,
		data.Name,
		data.InitialCost,
		data.InitialProfit,
		data.CookingTime,
		now,
		now,
	).Scan(&lastInsertId)
	if err != nil {
		if database.IsDuplicateError(err) {
			return nil, apperror.ErrConflict
		}

		return nil, err
	}

	return &lastInsertId, err
}

func (r *foodItemRepository) UpdateFoodDB(ctx context.Context, id int64, data entities.FoodItem) (err error) {
	now := helper.NowUTC()

	res, err := r.db.ExecContext(ctx,
		updateFoodItemQuery,
		data.Name,
		data.InitialProfit,
		data.CookingTime,
		data.InitialCost,
		now,
		id,
	)
	if err != nil {
		return err
	}

	rows, _ := res.RowsAffected()
	if rows == 0 {
		return apperror.ErrNoUpdateRecord
	}

	return err
}

func (r *foodItemRepository) GetFoodBySlugDB(ctx context.Context, slug string) (res *entities.FoodItem, err error) {
	row := r.db.QueryRowContext(ctx, getFoodsBySlugQuery, slug)
	return r.scanFoodItemRow(row)
}

func (r *foodItemRepository) GetFoodByIDDB(ctx context.Context, id int64) (res *entities.FoodItem, err error) {
	row := r.db.QueryRowContext(ctx, getFoodsByIDQuery, id)
	return r.scanFoodItemRow(row)
}

func (r *foodItemRepository) GetFoodsDB(ctx context.Context, limit, offset int) (res []entities.FoodItem, err error) {
	rows, err := r.db.QueryContext(ctx, getFoodsQuery, limit, offset)
	if err != nil {
		return nil, err
	}

	defer rows.Close()

	for rows.Next() {
		var data entities.FoodItem

		err := rows.Scan(
			&data.ID,
			&data.Slug,
			&data.Name,
			&data.InitialCost,
			&data.InitialProfit,
			&data.CookingTime,
		)
		if err != nil {
			return nil, err
		}

		res = append(res, data)
	}

	return res, nil
}

func (r *foodItemRepository) scanFoodItemRow(row *sql.Row) (*entities.FoodItem, error) {
	var data entities.FoodItem

	err := row.Scan(
		&data.ID,
		&data.Slug,
		&data.Name,
		&data.InitialCost,
		&data.InitialProfit,
		&data.CookingTime,
	)

	if errors.Is(err, sql.ErrNoRows) {
		return nil, nil
	} else if err != nil {
		return nil, err
	}

	return &data, err
}

func (r *foodItemRepository) CountFoodItemDB(ctx context.Context) (count int64, err error) {
	err = r.db.QueryRowContext(ctx, countFoodItemsQuery).Scan(&count)
	if err != nil {
		return 0, err
	}

	return count, nil
}

func (r *foodItemRepository) GetFoodItemIDsBySlugsDB(ctx context.Context, slugs []string) (map[string]int64, error) {
	query := `SELECT id, slug FROM food_items WHERE slug = ANY($1)`
	rows, err := r.db.QueryContext(ctx, query, pq.Array(slugs))
	if err != nil {
		return nil, err
	}

	defer rows.Close()

	result := make(map[string]int64)
	for rows.Next() {
		var id int64
		var slug string
		if err := rows.Scan(&id, &slug); err != nil {
			return nil, err
		}
		result[slug] = id
	}
	return result, nil
}

func (r *foodItemRepository) GetOverrideLevelDB(ctx context.Context, foodItemID int64) (res []entities.FoodItemOverrideLevel, err error) {
	rows, err := r.db.QueryContext(ctx, getOverrideLevelQuery, foodItemID)
	if err != nil {
		return nil, err
	}

	defer rows.Close()

	for rows.Next() {
		var data entities.FoodItemOverrideLevel
		err := rows.Scan(
			&data.FoodItemID,
			&data.Level,
			&data.Cost,
			&data.Profit,
			&data.PreparationTime,
		)
		if err != nil {
			return nil, err
		}

		res = append(res, data)
	}

	return res, nil
}

func (r *foodItemRepository) GetOverrideLevelByFoodItemIDAndLevelDB(ctx context.Context, foodItemID int64, level int) (*entities.FoodItemOverrideLevel, error) {
	var data entities.FoodItemOverrideLevel
	err := r.db.QueryRowContext(ctx, getOverrideLevelByFoodItemIDAndLevelQuery, foodItemID, level).Scan(
		&data.FoodItemID,
		&data.Level,
		&data.Cost,
		&data.Profit,
		&data.PreparationTime,
	)
	if errors.Is(err, sql.ErrNoRows) {
		return nil, nil
	} else if err != nil {
		return nil, err
	}

	return &data, nil
}

func (r *foodItemRepository) CreateOverrideLevelDB(ctx context.Context, foodItemID int64, data []entities.FoodItemOverrideLevel) (err error) {
	if len(data) == 0 {
		return nil
	}

	numFields := 7
	query := r.BuildBulkInsertQuery(insertFoodItemOverrideLevelQuery, len(data), numFields, "")
	args := make([]interface{}, 0, len(data)*numFields)
	now := helper.NowUTC()

	for _, item := range data {
		args = append(args, foodItemID, item.Level, item.Cost, item.Profit, item.PreparationTime, now, now)
	}

	_, err = r.db.ExecContext(ctx, query, args...)
	if database.IsDuplicateError(err) {
		return apperror.ErrConflict
	} else if err != nil {
		return err
	}

	return nil
}

func (r *foodItemRepository) DeleteOverrideLevelDB(ctx context.Context, foodItemID int64) (err error) {
	_, err = r.db.ExecContext(ctx,
		deleteFoodItemOverrideLevelQuery,
		foodItemID,
	)
	if err != nil {
		return err
	}

	return nil
}
