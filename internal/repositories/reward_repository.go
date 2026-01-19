package repositories

import (
	"context"
	"database/sql"
	"errors"
	"github.com/redis/go-redis/v9"
	"github.com/winartodev/cat-cafe/internal/entities"
	"github.com/winartodev/cat-cafe/pkg/apperror"
	"github.com/winartodev/cat-cafe/pkg/database"
	"github.com/winartodev/cat-cafe/pkg/helper"
)

type RewardRepository interface {
	CreateRewardTypeDB(ctx context.Context, data entities.RewardType) (id *int64, err error)
	UpdateRewardTypesDB(ctx context.Context, id int64, data entities.RewardType) (err error)
	GetRewardTypeBySlugDB(ctx context.Context, slug string) (res *entities.RewardType, err error)
	GetRewardTypesDB(ctx context.Context) (res []entities.RewardType, err error)
	GetRewardTypeByIDDB(ctx context.Context, id int64) (res *entities.RewardType, err error)

	CreateRewardDB(ctx context.Context, data entities.Reward) (id *int64, err error)
	UpdateRewardDB(ctx context.Context, id int64, data entities.Reward) (err error)
	GetRewardsDB(ctx context.Context, limit, offset int) (res []entities.Reward, err error)
	GetRewardBySlugDB(ctx context.Context, slug string) (data *entities.Reward, err error)
	CountRewardDB(ctx context.Context) (count int64, err error)
}

type rewardRepository struct {
	db    *sql.DB
	redis *redis.Client
}

func NewRewardRepository(db *sql.DB, redis *redis.Client) RewardRepository {
	return &rewardRepository{
		db:    db,
		redis: redis,
	}
}

func (r *rewardRepository) CreateRewardTypeDB(ctx context.Context, data entities.RewardType) (id *int64, err error) {
	now := helper.NowUTC()
	var lastInsertId int64

	err = r.db.QueryRowContext(ctx,
		rewardTypeInsertQuery,
		data.Slug,
		data.Name,
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

func (r *rewardRepository) UpdateRewardTypesDB(ctx context.Context, id int64, data entities.RewardType) (err error) {
	now := helper.NowUTC()

	res, err := r.db.ExecContext(ctx,
		updateRewardTypeQuery,
		data.Name,
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

func (r *rewardRepository) GetRewardTypeBySlugDB(ctx context.Context, slug string) (res *entities.RewardType, err error) {
	row := r.db.QueryRowContext(ctx, getRewardTypeBySlugQuery, slug)
	return r.scanRewardTypeRow(row)
}

func (r *rewardRepository) GetRewardTypesDB(ctx context.Context) (res []entities.RewardType, err error) {
	rows, err := r.db.QueryContext(ctx, getRewardTypesQuery)
	if err != nil {
		return nil, err
	}

	defer rows.Close()

	for rows.Next() {
		var id int64
		var row entities.RewardType

		err := rows.Scan(
			&id,
			&row.Slug,
			&row.Name,
		)

		if err != nil {
			return nil, err
		}

		row.ID = &id
		res = append(res, row)
	}

	return res, err
}

func (r *rewardRepository) GetRewardTypeByIDDB(ctx context.Context, id int64) (res *entities.RewardType, err error) {
	row := r.db.QueryRowContext(ctx, getRewardTypeByIDQuery, id)
	return r.scanRewardTypeRow(row)
}

func (r *rewardRepository) CreateRewardDB(ctx context.Context, data entities.Reward) (id *int64, err error) {
	now := helper.NowUTC()
	var lastInsertId int64

	err = r.db.QueryRowContext(ctx,
		rewardInsertQuery,
		data.RewardType.ID,
		data.Slug,
		data.Name,
		data.Amount,
		data.IsActive,
		now,
		now,
	).Scan(&lastInsertId)

	if database.IsDuplicateError(err) {
		return nil, apperror.ErrConflict
	} else if err != nil {
		return nil, err
	}

	return &lastInsertId, err
}

func (r *rewardRepository) UpdateRewardDB(ctx context.Context, id int64, data entities.Reward) (err error) {
	//TODO implement me
	panic("implement me")
}

func (r *rewardRepository) GetRewardBySlugDB(ctx context.Context, slug string) (data *entities.Reward, err error) {
	row := r.db.QueryRowContext(ctx, getRewardBySlugQuery, slug)
	return r.scanRewardRow(row)
}

func (r *rewardRepository) GetRewardsDB(ctx context.Context, limit, offset int) (res []entities.Reward, err error) {
	rows, err := r.db.QueryContext(ctx, getRewardQuery, limit, offset)
	if err != nil {
		return nil, err
	}

	defer rows.Close()

	for rows.Next() {
		var reward entities.Reward
		var rewardType entities.RewardType

		err := rows.Scan(
			&reward.ID,
			&reward.Slug,
			&reward.Name,
			&reward.Amount,
			&reward.IsActive,
			&rewardType.Slug,
			&rewardType.Name,
		)
		if err != nil {
			return nil, err
		}

		reward.RewardType = &rewardType
		res = append(res, reward)
	}

	return res, nil
}

func (r *rewardRepository) CountRewardDB(ctx context.Context) (count int64, err error) {
	err = r.db.QueryRowContext(ctx, countRewardsQuery).Scan(&count)
	if err != nil {
		return 0, err
	}

	return count, nil
}

func (r *rewardRepository) scanRewardTypeRow(row *sql.Row) (*entities.RewardType, error) {
	var res entities.RewardType
	var id int64
	err := row.Scan(&id, &res.Slug, &res.Name)

	if errors.Is(err, sql.ErrNoRows) {
		return nil, nil
	} else if err != nil {
		return nil, err
	}

	res.ID = &id

	return &res, err
}

func (r *rewardRepository) scanRewardRow(row *sql.Row) (*entities.Reward, error) {
	var reward entities.Reward
	var rewardType entities.RewardType

	err := row.Scan(
		&reward.ID,
		&reward.Slug,
		&reward.Name,
		&reward.Amount,
		&reward.IsActive,
		&rewardType.Slug,
		&rewardType.Name,
	)

	if errors.Is(err, sql.ErrNoRows) {
		return nil, nil
	} else if err != nil {
		return nil, err
	}

	reward.RewardType = &rewardType

	return &reward, err
}
