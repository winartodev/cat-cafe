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
	"github.com/winartodev/cat-cafe/pkg/helper"
	"time"
)

const (
	dailyRewardMasterRedisKey = "master:daily_rewards"
)

type DailyRewardRepository interface {
	GetTx() *sql.Tx
	WithTx(tx *sql.Tx) DailyRewardRepository

	CreateDailyRewardDB(ctx context.Context, data entities.DailyReward) (id *int64, err error)
	GetDailyRewardsWithPaginationDB(ctx context.Context, limit, offset int) (res []entities.DailyReward, err error)
	GetDailyRewardsDB(ctx context.Context) (res []entities.DailyReward, err error)
	GetDailyRewardByIDDB(ctx context.Context, id int64) (res *entities.DailyReward, err error)
	UpdateDailyRewardDB(ctx context.Context, id int64, data entities.DailyReward) (err error)
	CountDailyRewardsDB(ctx context.Context) (count int64, err error)

	DailyRewardWithTx(ctx context.Context, fn func(txRepo DailyRewardRepository) error) (err error)

	//GetDailyRewardsRedis(ctx context.Context) (res []entities.DailyReward, err error)
	//SetDailyRewardsRedis(ctx context.Context, data []entities.DailyReward) (err error)
	//DeleteDailyRewardsRedis(ctx context.Context) error
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

func (r *dailyRewardRepository) CreateDailyRewardDB(ctx context.Context, data entities.DailyReward) (id *int64, err error) {
	now := helper.NowUTC()
	var lastInsertId int64

	err = r.db.QueryRowContext(ctx, insertDailyRewardQuery,
		data.Reward.ID,
		data.DayNumber,
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
		var reward entities.Reward
		var rewardType entities.RewardType

		err := rows.Scan(
			&row.ID,
			&reward.ID,
			&row.DayNumber,
			&reward.Slug,
			&reward.Name,
			&reward.Amount,
			&reward.IsActive,
			&row.IsActive,
			&row.Description,
			&rewardType.Slug,
			&rewardType.Name,
		)

		if err != nil {
			return nil, err
		}

		reward.RewardType = &rewardType
		row.Reward = &reward

		res = append(res, row)
	}

	return res, err
}

func (r *dailyRewardRepository) GetDailyRewardsWithPaginationDB(ctx context.Context, limit, offset int) (res []entities.DailyReward, err error) {
	var dailyRewards []entities.DailyReward

	rows, err := r.db.QueryContext(ctx, getDailyRewardsWithPaginationQuery, limit, offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var dailyReward entities.DailyReward
		var reward entities.Reward
		var rewardType entities.RewardType

		err := rows.Scan(
			&dailyReward.ID,
			&reward.ID,
			&dailyReward.DayNumber,
			&reward.Slug,
			&reward.Name,
			&reward.Amount,
			&reward.IsActive,
			&dailyReward.IsActive,
			&dailyReward.Description,
			&rewardType.Slug,
			&rewardType.Name,
		)
		if err != nil {
			return nil, err
		}

		reward.RewardType = &rewardType
		dailyReward.Reward = &reward

		dailyRewards = append(dailyRewards, dailyReward)
	}

	return dailyRewards, nil

}

func (r *dailyRewardRepository) GetDailyRewardByIDDB(ctx context.Context, id int64) (res *entities.DailyReward, err error) {
	row := r.db.QueryRowContext(ctx, getDailyRewardByIDQuery, id)
	return r.scanDailyRewardTypeRow(row)
}

func (r *dailyRewardRepository) UpdateDailyRewardDB(ctx context.Context, id int64, data entities.DailyReward) (err error) {
	now := helper.NowUTC()

	res, err := r.db.ExecContext(ctx,
		updateDailyRewardQuery,
		data.Reward.ID,
		data.DayNumber,
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

func (r *dailyRewardRepository) CountDailyRewardsDB(ctx context.Context) (count int64, err error) {
	err = r.db.QueryRowContext(ctx, countDailyRewardsQuery).Scan(&count)
	if err != nil {
		return 0, err
	}

	return count, nil
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
	var reward entities.Reward

	err := row.Scan(
		&dailyReward.ID,
		&reward.ID,
		&reward.Slug,
		&reward.Name,
		&dailyReward.DayNumber,
		&reward.Amount,
		&dailyReward.IsActive,
		&dailyReward.Description,
	)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		}
		return nil, err
	}

	dailyReward.Reward = &reward

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
