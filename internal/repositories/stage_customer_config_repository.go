package repositories

import (
	"context"
	"database/sql"
	"github.com/winartodev/cat-cafe/internal/entities"
	"github.com/winartodev/cat-cafe/pkg/apperror"
	"github.com/winartodev/cat-cafe/pkg/database"
	"time"
)

type StageCustomerConfigRepository interface {
	WithTx(tx *sql.Tx) StageCustomerConfigRepository
	CreateCustomerConfigWithTxDB(ctx context.Context, stageID int64, data *entities.StageCustomerConfig) (*int64, error)
	UpdateCustomerConfigWithTxDB(ctx context.Context, stageID int64, data *entities.StageCustomerConfig) error
}

type stageCustomerConfigRepository struct {
	BaseRepository
}

func NewStageCustomerRepository(db *sql.DB) StageCustomerConfigRepository {
	return &stageCustomerConfigRepository{
		BaseRepository{
			db:   db,
			pool: db,
		},
	}
}

func (r *stageCustomerConfigRepository) WithTx(tx *sql.Tx) StageCustomerConfigRepository {
	if tx == nil {
		return r
	}

	return &stageCustomerConfigRepository{
		BaseRepository{
			db:   tx,
			pool: r.pool,
		},
	}
}

func (r *stageCustomerConfigRepository) CreateCustomerConfigWithTxDB(ctx context.Context, stageID int64, data *entities.StageCustomerConfig) (*int64, error) {
	var id int64
	now := time.Now()

	err := r.db.QueryRowContext(ctx, insertIntoCustomerConfig,
		stageID,
		data.CustomerSpawnTime,
		data.MaxCustomerOrderCount,
		data.MaxCustomerOrderVariant,
		data.StartingOrderTableCount,
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

func (r *stageCustomerConfigRepository) UpdateCustomerConfigWithTxDB(ctx context.Context, stageID int64, data *entities.StageCustomerConfig) error {
	now := time.Now()

	_, err := r.db.ExecContext(ctx, updateCustomerConfigQuery,
		stageID,
		data.CustomerSpawnTime,
		data.MaxCustomerOrderCount,
		data.MaxCustomerOrderVariant,
		data.StartingOrderTableCount,
		now,
	)
	if err != nil {
		return err
	}

	return nil
}
