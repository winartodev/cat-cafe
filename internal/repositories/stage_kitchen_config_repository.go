package repositories

import (
	"context"
	"database/sql"
	"github.com/lib/pq"
	"github.com/winartodev/cat-cafe/internal/entities"
	"github.com/winartodev/cat-cafe/pkg/apperror"
	"github.com/winartodev/cat-cafe/pkg/database"
	"github.com/winartodev/cat-cafe/pkg/helper"
)

type StageKitchenConfigRepository interface {
	WithTx(tx *sql.Tx) StageKitchenConfigRepository
	CreateKitchenConfigWithTxDB(ctx context.Context, stageID int64, data *entities.StageKitchenConfig) (id *int64, err error)
	UpdateKitchenConfigWithTxDB(ctx context.Context, stageID int64, data *entities.StageKitchenConfig) (id *int64, err error)
	GetKitchenConfigByStageIDDB(ctx context.Context, stageID int64) (res *entities.StageKitchenConfig, err error)

	GetKitchenCompletionRewardsDB(ctx context.Context, kitchenConfigID int64) (data []entities.KitchenPhaseCompletionRewards, err error)
	CreateKitchenCompletionRewardDB(ctx context.Context, kitchenConfigID int64, data *entities.KitchenPhaseCompletionRewards) (id *int64, err error)
	DeleteKitchenCompletionRewardDB(ctx context.Context, kitchenConfigID int64) error
}

type stageKitchenConfigRepository struct {
	BaseRepository
}

func NewStageKitchenConfigRepository(db *sql.DB) StageKitchenConfigRepository {
	return &stageKitchenConfigRepository{
		BaseRepository{
			db:   db,
			pool: db,
		},
	}
}

func (r *stageKitchenConfigRepository) WithTx(tx *sql.Tx) StageKitchenConfigRepository {
	if tx == nil {
		return r
	}

	return &stageKitchenConfigRepository{
		BaseRepository{
			db:   tx,
			pool: r.pool,
		},
	}
}

func (r *stageKitchenConfigRepository) CreateKitchenConfigWithTxDB(ctx context.Context, stageID int64, data *entities.StageKitchenConfig) (*int64, error) {
	var id int64
	now := helper.NowUTC()
	err := r.db.QueryRowContext(ctx, insertStageKitchenConfigQuery,
		stageID,
		&data.MaxLevel,
		&data.UpgradeProfitMultiply,
		&data.UpgradeCostMultiply,
		pq.Array(&data.TransitionPhaseLevels),
		pq.Array(&data.PhaseProfitMultipliers),
		pq.Array(&data.PhaseUpgradeCostMultipliers),
		pq.Array(&data.TableCountPerPhases),
		now,
		now,
	).Scan(&id)
	if database.IsDuplicateError(err) {
		return nil, apperror.ErrConflict
	} else if err != nil {
		return nil, err
	}

	return &id, nil
}

func (r *stageKitchenConfigRepository) UpdateKitchenConfigWithTxDB(ctx context.Context, stageID int64, data *entities.StageKitchenConfig) (*int64, error) {
	var id int64
	now := helper.NowUTC()
	err := r.db.QueryRowContext(ctx, updateStageKitchenConfigQuery,
		stageID,
		&data.MaxLevel,
		&data.UpgradeProfitMultiply,
		&data.UpgradeCostMultiply,
		pq.Array(&data.TransitionPhaseLevels),
		pq.Array(&data.PhaseProfitMultipliers),
		pq.Array(&data.PhaseUpgradeCostMultipliers),
		pq.Array(&data.TableCountPerPhases),
		now,
	).Scan(&id)
	if database.IsDuplicateError(err) {
		return nil, apperror.ErrConflict
	} else if err != nil {
		return nil, err
	}

	return &id, nil
}

func (r *stageKitchenConfigRepository) GetKitchenConfigByStageIDDB(ctx context.Context, stageID int64) (res *entities.StageKitchenConfig, err error) {
	var data entities.StageKitchenConfig
	err = r.db.QueryRowContext(ctx, getStageKitchenConfigQuery, stageID).Scan(
		&data.ID,
		&data.MaxLevel,
		&data.UpgradeProfitMultiply,
		&data.UpgradeCostMultiply,
		pq.Array(&data.TransitionPhaseLevels),
		pq.Array(&data.PhaseProfitMultipliers),
		pq.Array(&data.PhaseUpgradeCostMultipliers),
		pq.Array(&data.TableCountPerPhases),
	)
	if err != nil {
		return nil, err
	}

	return &data, nil
}

func (r *stageKitchenConfigRepository) CreateKitchenCompletionRewardDB(ctx context.Context, kitchenConfigID int64, data *entities.KitchenPhaseCompletionRewards) (*int64, error) {
	var id int64
	now := helper.NowUTC()
	err := r.db.QueryRowContext(ctx,
		insertKitchenPhaseCompletionRewardsQuery,
		kitchenConfigID,
		data.PhaseNumber,
		data.RewardID,
		now,
		now,
	).Scan(&id)
	if database.IsDuplicateError(err) {
		return nil, apperror.ErrConflict
	} else if err != nil {
		return nil, err
	}

	return &id, nil
}

func (r *stageKitchenConfigRepository) DeleteKitchenCompletionRewardDB(ctx context.Context, kitchenConfigID int64) error {
	_, err := r.db.ExecContext(ctx, deleteKitchenPhaseCompletionRewardsQuery, kitchenConfigID)
	if err != nil {
		return err
	}

	return nil
}

func (r *stageKitchenConfigRepository) GetKitchenCompletionRewardsDB(ctx context.Context, kitchenConfigID int64) (res []entities.KitchenPhaseCompletionRewards, err error) {
	rows, err := r.db.QueryContext(ctx, getKitchenPhaseCompletionRewardsQuery, kitchenConfigID)
	if err != nil {
		return nil, err
	}

	for rows.Next() {
		var data entities.KitchenPhaseCompletionRewards
		err = rows.Scan(
			&data.KitchenConfigID,
			&data.PhaseNumber,
			&data.RewardID,
		)
		if err != nil {
			return nil, err
		}

		res = append(res, data)
	}

	return res, nil
}
