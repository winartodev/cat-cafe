package repositories

import (
	"context"
	"database/sql"
	"errors"
	"github.com/winartodev/cat-cafe/internal/entities"
	"github.com/winartodev/cat-cafe/pkg/apperror"
	"time"
)

type UserDailyRewardRepository interface {
	GetTx() *sql.Tx
	WithTx(tx *sql.Tx) UserDailyRewardRepository

	GetUserDailyRewardByIDDB(ctx context.Context, id int64) (res *entities.UserDailyReward, err error)
	UpsertUserProgressionWithTx(ctx context.Context, userID int64, streak int64, lastClaim time.Time) (err error)
}

type userDailyRewardRepository struct {
	BaseRepository
}

func NewUserDailyRewardRepository(db *sql.DB) UserDailyRewardRepository {
	return &userDailyRewardRepository{
		BaseRepository{db: db, tx: nil},
	}
}

func (r *userDailyRewardRepository) WithTx(tx *sql.Tx) UserDailyRewardRepository {
	if tx == nil {
		return r
	}

	return &userDailyRewardRepository{BaseRepository{db: r.db, tx: tx}}
}

func (r *userDailyRewardRepository) GetUserDailyRewardByIDDB(ctx context.Context, id int64) (res *entities.UserDailyReward, err error) {
	var data entities.UserDailyReward
	err = r.db.QueryRowContext(ctx, getUserDailyRewardQuery, id).Scan(
		&data.ID,
		&data.UserID,
		&data.CurrentStreak,
		&data.LastClaimDate,
	)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		}

		return nil, err
	}

	return &data, err
}

func (r *userDailyRewardRepository) UpsertUserProgressionWithTx(ctx context.Context, userID int64, streak int64, lastClaim time.Time) (err error) {
	if r.tx == nil {
		return apperror.ErrRequiredActiveTx
	}

	now := time.Now()

	_, err = r.tx.ExecContext(ctx, upsertUserDailyRewardQuery,
		userID,
		streak,
		lastClaim,
		now,
	)

	return err
}
