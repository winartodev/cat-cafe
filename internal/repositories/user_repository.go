package repositories

import (
	"context"
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"time"

	"github.com/redis/go-redis/v9"
	"github.com/winartodev/cat-cafe/internal/entities"
	"github.com/winartodev/cat-cafe/pkg/apperror"
	"github.com/winartodev/cat-cafe/pkg/helper"
)

const (
	userTokenBlacklistRedisKey = "user:token:blacklist:%s"
	userIdRedisKey             = "user:id:%d"
)

type UserRepository interface {
	GetTx() *sql.Tx
	WithTx(tx *sql.Tx) UserRepository

	CreateUserDB(ctx context.Context, data *entities.User) (*int64, error)
	GetUserByIDDB(ctx context.Context, id int64) (res *entities.User, err error)
	GetUserByEmailDB(ctx context.Context, email string) (res *entities.User, err error)
	GetUserBalanceByIDDB(ctx context.Context, id int64) (res *entities.UserBalance, err error)

	BalanceWithTx(ctx context.Context, fn func(txRepo UserRepository) error) error
	UpdateUserBalanceWithTx(ctx context.Context, userID int64, rewardType entities.UserBalanceType, amount int64) (err error)
	UpdateLastSyncBalanceWithTx(ctx context.Context, userID int64, lastSyncTime time.Time) (err error)

	SetUserRedis(ctx context.Context, userID int64, data *entities.UserCache, exp time.Duration) (err error)
	GetUserRedis(ctx context.Context, userID int64) (res *entities.UserCache, err error)
	DeleteUserRedis(ctx context.Context, userID int64) (err error)

	BlacklistToken(ctx context.Context, token string, expiration time.Duration) error
	IsTokenBlacklisted(ctx context.Context, token string) bool
}

type userRepository struct {
	BaseRepository
}

func NewUserRepository(db *sql.DB, redis *redis.Client) UserRepository {
	return &userRepository{BaseRepository{db: db, tx: nil, redis: redis}}
}

func (r *userRepository) WithTx(tx *sql.Tx) UserRepository {
	if tx == nil {
		return r
	}

	return &userRepository{BaseRepository{db: r.db, tx: tx}}
}

func (r *userRepository) CreateUserDB(ctx context.Context, data *entities.User) (*int64, error) {
	now := helper.NowUTC()
	var id *int64
	err := r.db.QueryRowContext(
		ctx,
		insertUserQuery,
		data.ExternalID,
		data.Username,
		data.Email,
		now,
		now,
	).Scan(&id)
	if err != nil {
		return nil, err
	}

	return id, nil
}

func (r *userRepository) GetUserByIDDB(ctx context.Context, id int64) (res *entities.User, err error) {
	row := r.db.QueryRowContext(ctx, getUserByIDQuery, id)
	return r.scanUserRow(row)
}

func (r *userRepository) GetUserByEmailDB(ctx context.Context, email string) (res *entities.User, err error) {
	row := r.db.QueryRowContext(ctx, getUserByEmailQuery, email)
	return r.scanUserRow(row)
}

func (r *userRepository) GetUserBalanceByIDDB(ctx context.Context, id int64) (res *entities.UserBalance, err error) {
	var userBalance entities.UserBalance

	err = r.db.QueryRowContext(ctx, getUserBalanceByIDQuery, id).Scan(
		&userBalance.Coin,
		&userBalance.Gem,
	)

	if errors.Is(err, sql.ErrNoRows) {
		return nil, nil
	} else if err != nil {
		return nil, err
	}

	return &userBalance, nil
}

func (r *userRepository) BalanceWithTx(ctx context.Context, fn func(txRepo UserRepository) error) error {
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

func (r *userRepository) UpdateUserBalanceWithTx(ctx context.Context, userID int64, balanceType entities.UserBalanceType, amount int64) error {
	if r.tx == nil {
		return apperror.ErrRequiredActiveTx
	}

	var query string
	switch balanceType {
	case entities.BalanceTypeCoin:
		query = `UPDATE users SET coin = coin + $1 WHERE id = $2`
	case entities.BalanceTypeGem:
		query = `UPDATE users SET gem = gem + $1 WHERE id = $2`
	default:
		return apperror.ErrUnknownRewardType
	}

	_, err := r.tx.ExecContext(ctx, query, amount, userID)

	return err
}

func (r *userRepository) UpdateLastSyncBalanceWithTx(ctx context.Context, userID int64, lastSyncTime time.Time) (err error) {
	if r.tx == nil {
		return apperror.ErrRequiredActiveTx
	}

	now := helper.NowUTC()
	_, err = r.tx.ExecContext(ctx, updateLastSyncBalanceQuery, lastSyncTime, now, userID)
	if err != nil {
		return err
	}

	return nil

}

func (r *userRepository) scanUserRow(row *sql.Row) (*entities.User, error) {
	var user entities.User
	var userBalance entities.UserBalance
	err := row.Scan(
		&user.ID,
		&user.ExternalID,
		&user.Username,
		&user.Email,
		&userBalance.Gem,
		&userBalance.Coin,
	)
	if errors.Is(err, sql.ErrNoRows) {
		return nil, nil
	} else if err != nil {
		return nil, err
	}

	user.UserBalance = &userBalance

	return &user, nil
}

func (r *userRepository) SetUserRedis(ctx context.Context, userID int64, data *entities.UserCache, exp time.Duration) error {
	key := r.userIDKey(userID)
	jsonData, err := json.Marshal(data)
	if err != nil {
		return err
	}

	return r.redis.Set(ctx, key, jsonData, exp).Err()
}

func (r *userRepository) GetUserRedis(ctx context.Context, userID int64) (*entities.UserCache, error) {
	key := r.userIDKey(userID)
	val, err := r.redis.Get(ctx, key).Result()
	if err != nil {
		return nil, err
	}

	var user entities.UserCache
	err = json.Unmarshal([]byte(val), &user)
	if err != nil {
		return nil, err
	}

	return &user, nil
}

func (r *userRepository) DeleteUserRedis(ctx context.Context, userID int64) error {
	key := r.userIDKey(userID)
	return r.redis.Del(ctx, key).Err()
}

func (r *userRepository) BlacklistToken(ctx context.Context, token string, expiration time.Duration) error {
	key := r.blacklistKey(token)
	return r.redis.Set(ctx, key, "1", expiration).Err()
}

func (r *userRepository) IsTokenBlacklisted(ctx context.Context, token string) bool {
	key := r.blacklistKey(token)

	_, err := r.redis.Get(ctx, key).Result()
	return err == nil
}

func (r *userRepository) blacklistKey(token string) string {
	return fmt.Sprintf(userTokenBlacklistRedisKey, token)
}

func (r *userRepository) userIDKey(userId int64) string {
	return fmt.Sprintf(userIdRedisKey, userId)
}
