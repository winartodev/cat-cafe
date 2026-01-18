package repositories

import (
	"context"
	"database/sql"
	"encoding/json"
	"errors"
	_ "github.com/lib/pq"
	"github.com/redis/go-redis/v9"
	"github.com/winartodev/cat-cafe/internal/entities"
	"github.com/winartodev/cat-cafe/pkg/apperror"
	"github.com/winartodev/cat-cafe/pkg/database"
	"github.com/winartodev/cat-cafe/pkg/helper"
	"time"
)

const (
	dailyRewardMasterRedisKey = "master:daily_rewards"
)

type DailyRewardRepository interface {
	GetTx() *sql.Tx
	WithTx(tx *sql.Tx) DailyRewardRepository

	CreateRewardTypeDB(ctx context.Context, data entities.RewardType) (id *int64, err error)
	UpdateRewardTypesDB(ctx context.Context, id int64, data entities.RewardType) (err error)
	GetRewardTypesDB(ctx context.Context) (res []entities.RewardType, err error)
	GetRewardTypeByIDDB(ctx context.Context, id int64) (res *entities.RewardType, err error)
	GetRewardTypeBySlugDB(ctx context.Context, slug string) (res *entities.RewardType, err error)

	CreateDailyRewardDB(ctx context.Context, data entities.DailyReward) (id *int64, err error)
	GetDailyRewardsDB(ctx context.Context) (res []entities.DailyReward, err error)
	GetDailyRewardByIDDB(ctx context.Context, id int64) (res *entities.DailyReward, err error)
	UpdateDailyRewardDB(ctx context.Context, id int64, data entities.DailyReward) (err error)

	DailyRewardWithTx(ctx context.Context, fn func(txRepo DailyRewardRepository) error) (err error)

	GetDailyRewardsRedis(ctx context.Context) (res []entities.DailyReward, err error)
	SetDailyRewardsRedis(ctx context.Context, data []entities.DailyReward) (err error)
	DeleteDailyRewardsRedis(ctx context.Context) error
}

type dailyRewardRepository struct {
	BaseRepository
}

func NewDailyRewardsRepository(db *sql.DB, redis *redis.Client) DailyRewardRepository {
	return &dailyRewardRepository{BaseRepository{db: db, tx: nil, redis: redis}}
}

func (r *dailyRewardRepository) WithTx(tx *sql.Tx) DailyRewardRepository {
	if tx == nil {
		return r
	}

	return &dailyRewardRepository{BaseRepository{db: r.db, tx: tx}}
}

func (r *dailyRewardRepository) CreateRewardTypeDB(ctx context.Context, data entities.RewardType) (id *int64, err error) {
	now := helper.TimeUTC()
	var lastInsertId int64

	err = r.db.QueryRowContext(ctx, rewardTypeInsertQuery, data.Slug, data.Name, now, now).Scan(&id)
	if err != nil {
		if database.IsDuplicateError(err) {
			return nil, apperror.ErrConflict
		}

		return nil, err
	}

	return &lastInsertId, err
}

func (r *dailyRewardRepository) GetRewardTypesDB(ctx context.Context) (res []entities.RewardType, err error) {
	rows, err := r.db.QueryContext(ctx, getRewardTypesQuery)
	if err != nil {
		return nil, err
	}

	defer rows.Close()

	for rows.Next() {
		var row entities.RewardType

		err := rows.Scan(
			&row.ID,
			&row.Slug,
			&row.Name,
		)

		if err != nil {
			return nil, err
		}

		res = append(res, row)
	}

	return res, err
}

func (r *dailyRewardRepository) UpdateRewardTypesDB(ctx context.Context, id int64, data entities.RewardType) (err error) {
	now := helper.TimeUTC()

	res, err := r.db.ExecContext(ctx, updateRewardTypeQuery, data.Name, now, id)
	if err != nil {
		return err
	}

	rows, _ := res.RowsAffected()
	if rows == 0 {
		return apperror.ErrNoUpdateRecord
	}

	return err
}

func (r *dailyRewardRepository) GetRewardTypeByIDDB(ctx context.Context, id int64) (res *entities.RewardType, err error) {
	row := r.db.QueryRowContext(ctx, getRewardTypeByIDQuery, id)
	return r.scanRewardTypeRow(row)
}

func (r *dailyRewardRepository) GetRewardTypeBySlugDB(ctx context.Context, slug string) (res *entities.RewardType, err error) {
	row := r.db.QueryRowContext(ctx, getRewardTypeBySlugQuery, slug)
	return r.scanRewardTypeRow(row)
}

func (r *dailyRewardRepository) CreateDailyRewardDB(ctx context.Context, data entities.DailyReward) (id *int64, err error) {
	now := helper.TimeUTC()
	var lastInsertId int64

	err = r.db.QueryRowContext(ctx, insertDailyRewardQuery,
		data.RewardType.ID,
		data.DayNumber,
		data.RewardAmount,
		data.IsActive,
		data.Description,
		now,
		now,
	).Scan(&lastInsertId)
	if err != nil {
		return nil, err
	}

	return &lastInsertId, err
}

func (r *dailyRewardRepository) GetDailyRewardsDB(ctx context.Context) (res []entities.DailyReward, err error) {
	rows, err := r.db.QueryContext(ctx, getDailyRewardsQuery)
	if err != nil {
		return nil, err
	}

	defer rows.Close()

	for rows.Next() {
		var row entities.DailyReward
		var rewardType entities.RewardType
		err := rows.Scan(
			&row.ID,
			&rewardType.ID,
			&rewardType.Slug,
			&rewardType.Name,
			&row.DayNumber,
			&row.RewardAmount,
			&row.IsActive,
			&row.Description,
		)

		if err != nil {
			return nil, err
		}

		row.RewardType = &rewardType

		res = append(res, row)
	}

	return res, err
}

func (r *dailyRewardRepository) GetDailyRewardByIDDB(ctx context.Context, id int64) (res *entities.DailyReward, err error) {
	row := r.db.QueryRowContext(ctx, getDailyRewardByIDQuery, id)
	return r.scanDailyRewardTypeRow(row)
}

func (r *dailyRewardRepository) UpdateDailyRewardDB(ctx context.Context, id int64, data entities.DailyReward) (err error) {
	now := helper.TimeUTC()

	res, err := r.db.ExecContext(ctx,
		updateDailyRewardQuery,
		data.RewardType.ID,
		data.DayNumber,
		data.RewardAmount,
		data.IsActive,
		data.Description,
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

func (r *dailyRewardRepository) DailyRewardWithTx(ctx context.Context, fn func(txRepo DailyRewardRepository) error) error {
	tx, err := r.db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}

	txRepo := r.WithTx(tx)

	err = fn(txRepo)
	if err != nil {
		return err
	}

	return tx.Commit()
}

func (r *dailyRewardRepository) scanRewardTypeRow(row *sql.Row) (*entities.RewardType, error) {
	var res entities.RewardType

	err := row.Scan(&res.ID, &res.Slug, &res.Name)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		}
		return nil, err
	}

	return &res, err
}

func (r *dailyRewardRepository) scanDailyRewardTypeRow(row *sql.Row) (*entities.DailyReward, error) {
	var dailyReward entities.DailyReward
	var rewardType entities.RewardType

	err := row.Scan(
		&dailyReward.ID,
		&rewardType.ID,
		&rewardType.Slug,
		&rewardType.Name,
		&dailyReward.DayNumber,
		&dailyReward.RewardAmount,
		&dailyReward.IsActive,
		&dailyReward.Description,
	)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		}
		return nil, err
	}

	dailyReward.RewardType = &rewardType

	return &dailyReward, err
}

func (r *dailyRewardRepository) GetDailyRewardsRedis(ctx context.Context) (res []entities.DailyReward, err error) {
	val, err := r.redis.Get(ctx, dailyRewardMasterRedisKey).Result()
	if errors.Is(err, redis.Nil) {
		return nil, nil
	} else if err != nil {
		return nil, err
	}

	err = json.Unmarshal([]byte(val), &res)
	if err != nil {
		return nil, err
	}

	return res, nil
}

func (r *dailyRewardRepository) SetDailyRewardsRedis(ctx context.Context, data []entities.DailyReward) (err error) {
	jsonData, err := json.Marshal(data)
	if err != nil {
		return err
	}

	return r.redis.Set(ctx, dailyRewardMasterRedisKey, jsonData, 168*time.Hour).Err()
}

func (r *dailyRewardRepository) DeleteDailyRewardsRedis(ctx context.Context) error {
	return r.redis.Del(ctx, dailyRewardMasterRedisKey).Err()
}
