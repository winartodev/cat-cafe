package repositories

import (
	"context"
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/redis/go-redis/v9"
	"github.com/winartodev/cat-cafe/internal/entities"
	"github.com/winartodev/cat-cafe/pkg/apperror"
	"time"
)

const (
	userDailyRewardRedisKey = "user_daily_reward:%d"
)

type UserDailyRewardRepository interface {
	GetTx() *sql.Tx
	WithTx(tx *sql.Tx) UserDailyRewardRepository

	GetUserDailyRewardByIDDB(ctx context.Context, id int64) (res *entities.UserDailyReward, err error)
	UpsertUserProgressionWithTx(ctx context.Context, userID int64, longestStreak int64, currentStreak int64, lastClaim time.Time) (err error)

	//GetUserDailyRewardRedis(ctx context.Context, userID int64) (res *entities.UserDailyReward, err error)
	//SetUserDailyRewardRedis(ctx context.Context, userID int64, progress *entities.UserDailyReward, ttl time.Duration) (err error)
	//DeleteUserDailyRewardRedis(ctx context.Context, userID int64) error
}

type userDailyRewardRepository struct {
	BaseRepository
}

func NewUserDailyRewardRepository(db *sql.DB, redis *redis.Client) UserDailyRewardRepository {
	return &userDailyRewardRepository{BaseRepository{db: db, tx: nil, redis: redis}}
}

func (r *userDailyRewardRepository) WithTx(tx *sql.Tx) UserDailyRewardRepository {
	if tx == nil {
		return r
	}

	return &userDailyRewardRepository{BaseRepository{db: r.db, tx: tx}}
}

func (r *userDailyRewardRepository) GetUserDailyRewardByIDDB(ctx context.Context, id int64) (res *entities.UserDailyReward, err error) {
	var data entities.UserDailyReward
	err = r.db.QueryRowContext(ctx, getUserDailyRewardProgressQuery, id).Scan(
		&data.ID,
		&data.UserID,
		&data.LongestStreak,
		&data.CurrentStreak,
		&data.LastClaimDate,
	)

	if errors.Is(err, sql.ErrNoRows) {
		return nil, nil
	} else if err != nil {
		return nil, err
	}

	return &data, err
}

func (r *userDailyRewardRepository) UpsertUserProgressionWithTx(ctx context.Context, userID int64, longestStreak int64, currentStreak int64, lastClaim time.Time) (err error) {
	if r.tx == nil {
		return apperror.ErrRequiredActiveTx
	}

	now := time.Now()

	_, err = r.tx.ExecContext(ctx, upsertUserDailyRewardProgressQuery,
		userID,
		longestStreak,
		currentStreak,
		lastClaim,
		now,
	)

	return err
}

func (r *userDailyRewardRepository) GetUserDailyRewardRedis(ctx context.Context, userID int64) (res *entities.UserDailyReward, err error) {
	key := r.userDailyRewardKey(userID)
	val, err := r.redis.Get(ctx, key).Result()

	if errors.Is(err, redis.Nil) {
		return nil, nil
	} else if err != nil {
		return nil, err
	}

	var progress entities.UserDailyReward
	err = json.Unmarshal([]byte(val), &progress)
	if err != nil {
		return nil, err
	}

	return &progress, nil
}

func (r *userDailyRewardRepository) SetUserDailyRewardRedis(ctx context.Context, userID int64, progress *entities.UserDailyReward, ttl time.Duration) (err error) {
	key := r.userDailyRewardKey(userID)
	data, err := json.Marshal(progress)
	if err != nil {
		return err
	}

	return r.redis.Set(ctx, key, data, ttl).Err()
}

func (r *userDailyRewardRepository) DeleteUserDailyRewardRedis(ctx context.Context, userID int64) error {
	key := r.userDailyRewardKey(userID)
	return r.redis.Del(ctx, key).Err()
}

func (r *userDailyRewardRepository) userDailyRewardKey(userID int64) string {
	return fmt.Sprintf(userDailyRewardRedisKey, userID)
}
