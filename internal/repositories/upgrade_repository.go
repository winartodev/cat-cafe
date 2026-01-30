package repositories

import (
	"context"
	"database/sql"

	"github.com/winartodev/cat-cafe/internal/entities"
	"github.com/winartodev/cat-cafe/pkg/apperror"
	"github.com/winartodev/cat-cafe/pkg/database"
	"github.com/winartodev/cat-cafe/pkg/helper"
)

type UpgradeRepository interface {
	WithTx(tx *sql.Tx) UpgradeRepository
	CreateUpgradeDB(ctx context.Context, data entities.Upgrade) (id *int64, err error)
	UpdateUpgradeDB(ctx context.Context, id int64, data entities.Upgrade) (err error)
	GetUpgradesDB(ctx context.Context, limit, offset int) (res []entities.Upgrade, err error)
	GetUpgradeByIDDB(ctx context.Context, id int64) (res *entities.Upgrade, err error)
	GetUpgradeBySlugDB(ctx context.Context, slug string) (res *entities.Upgrade, err error)
	GetActiveUpgradesDB(ctx context.Context) (res []entities.Upgrade, err error)
	CountUpgradesDB(ctx context.Context) (totalRows int64, err error)
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
	now := helper.NowUTC()

	var lastInsertId int64

	err = r.db.QueryRowContext(ctx, insertUpgradeQuery,
		data.Slug,
		data.Name,
		data.Description,
		data.Cost,
		data.CostType,
		data.Effect.Type,
		data.Effect.Value,
		data.Effect.Unit,
		data.Effect.Target,
		data.Effect.TargetID,
		data.IsActive,
		data.Sequence,
		now,
		now,
	).Scan(&lastInsertId)
	if database.IsDuplicateError(err) {
		return nil, apperror.ErrorAlreadyExists("upgrade", "slug", data.Slug)
	} else if err != nil {
		return nil, err
	}

	id = &lastInsertId

	return id, nil
}

// GetActiveUpgradesDB implements [UpgradeRepository].
func (r *upgradeRepository) GetActiveUpgradesDB(ctx context.Context) (res []entities.Upgrade, err error) {
	panic("unimplemented")
}

// GetUpgradeByIDDB implements [UpgradeRepository].
func (r *upgradeRepository) GetUpgradeByIDDB(ctx context.Context, id int64) (res *entities.Upgrade, err error) {
	rows, err := r.db.QueryContext(ctx, getUpgradeByIDQuery, id)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var data entities.Upgrade
	var effect entities.UpgradeEffect
	if rows.Next() {
		if err := rows.Scan(
			&data.ID,
			&data.Slug,
			&data.Name,
			&data.Description,
			&data.IsActive,
			&data.Sequence,
			&data.Cost,
			&data.CostType,
			&effect.Type,
			&effect.Value,
			&effect.Unit,
			&effect.Target,
			&effect.TargetID,
			&effect.TargetName,
		); err != nil {
			return nil, err
		}
	}

	data.Effect = effect

	return &data, nil
}

// GetUpgradeBySlugDB implements [UpgradeRepository].
func (r *upgradeRepository) GetUpgradeBySlugDB(ctx context.Context, slug string) (res *entities.Upgrade, err error) {
	panic("unimplemented")
}

// GetUpgradesDB implements [UpgradeRepository].
func (r *upgradeRepository) GetUpgradesDB(ctx context.Context, limit int, offset int) (res []entities.Upgrade, err error) {
	rows, err := r.db.QueryContext(ctx, getUpgradesQuery, limit, offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var upgrade entities.Upgrade
		if err := rows.Scan(
			&upgrade.ID,
			&upgrade.Slug,
			&upgrade.Name,
			&upgrade.Description,
			&upgrade.IsActive,
			&upgrade.Sequence,
		); err != nil {
			return nil, err
		}
		res = append(res, upgrade)
	}

	return res, nil
}

// UpdateUpgradeDB implements [UpgradeRepository].
func (r *upgradeRepository) UpdateUpgradeDB(ctx context.Context, id int64, data entities.Upgrade) (err error) {
	panic("unimplemented")
}

func (r *upgradeRepository) CountUpgradesDB(ctx context.Context) (totalRows int64, err error) {
	err = r.db.QueryRowContext(ctx, countUpgradesQuery).Scan(&totalRows)
	if err != nil {
		return 0, err
	}

	return totalRows, nil
}
