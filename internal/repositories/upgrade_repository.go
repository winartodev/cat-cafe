package repositories

import (
	"context"
	"database/sql"

	"github.com/winartodev/cat-cafe/internal/entities"
)

type UpgradeRepository interface {
	WithTx(tx *sql.Tx) UpgradeRepository
	CreateUpgradeDB(ctx context.Context, data entities.Upgrade) (id *int64, err error)
	UpdateUpgradeDB(ctx context.Context, id int64, data entities.Upgrade) (err error)
	GetUpgradesDB(ctx context.Context, limit, offset int) (res []entities.Upgrade, totalRows int64, err error)
	GetUpgradeByIDDB(ctx context.Context, id int64) (res *entities.Upgrade, err error)
	GetUpgradeBySlugDB(ctx context.Context, slug string) (res *entities.Upgrade, err error)
	GetActiveUpgradesDB(ctx context.Context) (res []entities.Upgrade, err error)
}

type upgradeRepository struct {
	BaseRepository
}

func NewUpgradeRepository(db *sql.DB) UpgradeRepository {
	return &upgradeRepository{
		BaseRepository: BaseRepository{
			db:   db,
			pool: db,
		},
	}
}

func (r *upgradeRepository) WithTx(tx *sql.Tx) UpgradeRepository {
	if tx == nil {
		return r
	}

	return &upgradeRepository{
		BaseRepository: BaseRepository{
			db:   tx,
			pool: r.pool,
		},
	}
}

// CreateUpgradeDB implements [UpgradeRepository].
func (r *upgradeRepository) CreateUpgradeDB(ctx context.Context, data entities.Upgrade) (id *int64, err error) {
	panic("unimplemented")
}

// GetActiveUpgradesDB implements [UpgradeRepository].
func (r *upgradeRepository) GetActiveUpgradesDB(ctx context.Context) (res []entities.Upgrade, err error) {
	panic("unimplemented")
}

// GetUpgradeByIDDB implements [UpgradeRepository].
func (r *upgradeRepository) GetUpgradeByIDDB(ctx context.Context, id int64) (res *entities.Upgrade, err error) {
	panic("unimplemented")
}

// GetUpgradeBySlugDB implements [UpgradeRepository].
func (r *upgradeRepository) GetUpgradeBySlugDB(ctx context.Context, slug string) (res *entities.Upgrade, err error) {
	panic("unimplemented")
}

// GetUpgradesDB implements [UpgradeRepository].
func (r *upgradeRepository) GetUpgradesDB(ctx context.Context, limit int, offset int) (res []entities.Upgrade, totalRows int64, err error) {
	panic("unimplemented")
}

// UpdateUpgradeDB implements [UpgradeRepository].
func (r *upgradeRepository) UpdateUpgradeDB(ctx context.Context, id int64, data entities.Upgrade) (err error) {
	panic("unimplemented")
}
