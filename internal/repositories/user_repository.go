package repositories

import (
	"context"
	"database/sql"
	"errors"
	"github.com/winartodev/cat-cafe/internal/entities"
	"github.com/winartodev/cat-cafe/pkg/apperror"
)

type UserRepository interface {
	GetTx() *sql.Tx
	WithTx(tx *sql.Tx) UserRepository

	GetUserByIDDB(ctx context.Context, id int64) (res *entities.User, err error)
	GetUserBalanceByIDDB(ctx context.Context, id int64) (res *entities.UserBalance, err error)
	UpdateUserBalanceWithTx(ctx context.Context, userID int64, rewardType entities.RewardTypeSlug, amount int64) error
}

type userRepository struct {
	BaseRepository
}

func NewUserRepository(db *sql.DB) UserRepository {
	return &userRepository{
		BaseRepository{db: db, tx: nil},
	}
}

func (r *userRepository) WithTx(tx *sql.Tx) UserRepository {
	if tx == nil {
		return r
	}

	return &userRepository{BaseRepository{db: r.db, tx: tx}}
}

func (r *userRepository) GetUserByIDDB(ctx context.Context, id int64) (res *entities.User, err error) {
	var user entities.User
	// TODO: FIX THIS QUERY IMMEDIATELY
	err = r.db.QueryRowContext(ctx, getUserByIDDB, id).Scan(
		&user.ID,
		&user.ExternalID,
		&user.Username,
	)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		}

		return nil, err
	}

	return &user, err
}

func (r *userRepository) GetUserBalanceByIDDB(ctx context.Context, id int64) (res *entities.UserBalance, err error) {
	var userBalance entities.UserBalance

	err = r.db.QueryRowContext(ctx, getUserBalanceByIDQuery, id).Scan(
		&userBalance.Coin,
		&userBalance.Gem,
	)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		}

		return nil, err
	}

	return &userBalance, nil
}

func (r *userRepository) UpdateUserBalanceWithTx(ctx context.Context, userID int64, rewardType entities.RewardTypeSlug, amount int64) error {
	if r.tx == nil {
		return apperror.ErrRequiredActiveTx
	}

	var query string
	switch rewardType {
	case entities.RewardTypeCoin:
		query = `UPDATE users SET coin = coin + $1 WHERE id = $2`
	case entities.RewardTypeGem:
		query = `UPDATE users SET gem = gem + $1 WHERE id = $2`
	case entities.RewardTypeGoPayCoin:
		return nil
	default:
		return apperror.ErrUnknownRewardType
	}

	_, err := r.tx.ExecContext(ctx, query, amount, userID)

	return err
}
