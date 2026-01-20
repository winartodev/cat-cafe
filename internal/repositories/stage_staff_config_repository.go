package repositories

import (
	"context"
	"database/sql"
	"github.com/winartodev/cat-cafe/internal/entities"
	"github.com/winartodev/cat-cafe/pkg/apperror"
	"github.com/winartodev/cat-cafe/pkg/database"
	"github.com/winartodev/cat-cafe/pkg/helper"
)

type StageStaffConfigRepository interface {
	WithTx(tx *sql.Tx) StageStaffConfigRepository
	CreateStaffConfigWithTxDB(ctx context.Context, stageID int64, data *entities.StageStaffConfig) (*int64, error)
	UpdateStaffConfigWithTxDB(ctx context.Context, stageID int64, data *entities.StageStaffConfig) error
}

type stageStaffConfigRepository struct {
	BaseRepository
}

func NewStageStaffConfigRepository(db *sql.DB) StageStaffConfigRepository {
	return &stageStaffConfigRepository{
		BaseRepository{
			db:   db,
			pool: db,
		},
	}
}

func (r *stageStaffConfigRepository) WithTx(tx *sql.Tx) StageStaffConfigRepository {
	if tx == nil {
		return r
	}

	return &stageStaffConfigRepository{
		BaseRepository{
			db:   tx,
			pool: r.pool,
		},
	}
}

func (r *stageStaffConfigRepository) CreateStaffConfigWithTxDB(ctx context.Context, stageID int64, data *entities.StageStaffConfig) (*int64, error) {
	now := helper.NowUTC()
	var id int64

	err := r.db.QueryRowContext(ctx,
		insertStageStaffConfigQuery,
		stageID,
		data.StartingStaffManager,
		data.StartingStaffManager,
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

func (r *stageStaffConfigRepository) UpdateStaffConfigWithTxDB(ctx context.Context, stageID int64, data *entities.StageStaffConfig) error {
	now := helper.NowUTC()

	_, err := r.db.ExecContext(ctx,
		updateStageStaffConfigQuery,
		stageID,
		data.StartingStaffManager,
		data.StartingStaffManager,
		now,
	)
	if err != nil {
		return err
	}

	return nil
}
