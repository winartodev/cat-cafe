package repositories

import (
	"context"
	"database/sql"
	"github.com/winartodev/cat-cafe/internal/entities"
	"github.com/winartodev/cat-cafe/pkg/apperror"
	"github.com/winartodev/cat-cafe/pkg/database"
	"github.com/winartodev/cat-cafe/pkg/helper"
)

type StageCameraConfigRepository interface {
	WithTx(tx *sql.Tx) StageCameraConfigRepository
	CreateStageCameraDB(ctx context.Context, stageID int64, data *entities.StageCameraConfig) (*int64, error)
	UpdateStageCameraDB(ctx context.Context, stageID int64, data *entities.StageCameraConfig) error
}

type stageCameraConfigRepository struct {
	BaseRepository
}

func NewStageCameraConfigRepository(db *sql.DB) StageCameraConfigRepository {
	return &stageCameraConfigRepository{
		BaseRepository{
			db:   db,
			pool: db,
		},
	}
}

func (r *stageCameraConfigRepository) WithTx(tx *sql.Tx) StageCameraConfigRepository {
	if tx == nil {
		return r
	}

	return &stageCameraConfigRepository{
		BaseRepository{
			db:   tx,
			pool: r.pool,
		},
	}
}

func (r *stageCameraConfigRepository) CreateStageCameraDB(ctx context.Context, stageID int64, data *entities.StageCameraConfig) (*int64, error) {
	now := helper.NowUTC()

	var id int64
	err := r.db.QueryRowContext(ctx, insertIntoStageCameraConfigQuery,
		stageID,
		&data.ZoomSize,
		&data.MinBoundX,
		&data.MinBoundY,
		&data.MaxBoundX,
		&data.MaxBoundY,
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

func (r *stageCameraConfigRepository) UpdateStageCameraDB(ctx context.Context, stageID int64, data *entities.StageCameraConfig) error {
	now := helper.NowUTC()

	_, err := r.db.ExecContext(ctx, updateStageCameraConfigQuery,
		stageID,
		&data.ZoomSize,
		&data.MinBoundX,
		&data.MinBoundY,
		&data.MaxBoundX,
		&data.MaxBoundY,
		now,
	)
	if err != nil {
		return err
	}

	return nil
}
