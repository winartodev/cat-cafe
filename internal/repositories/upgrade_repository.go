package repositories

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"github.com/lib/pq"

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
	GetActiveUpgradesDB(ctx context.Context, stageID int64) (res []entities.Upgrade, err error)
	CountUpgradesDB(ctx context.Context) (totalRows int64, err error)
	GetUpgradesBySlugsDB(ctx context.Context, slugs []string) ([]entities.Upgrade, error)
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

func (r *upgradeRepository) GetActiveUpgradesDB(ctx context.Context, stageID int64) (res []entities.Upgrade, err error) {
	rows, err := r.db.QueryContext(ctx, getActiveUpgradesQuery)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var data entities.Upgrade
		var effect entities.UpgradeEffect
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
		data.Effect = effect
		res = append(res, data)
	}

	return res, nil
}

func (r *upgradeRepository) GetUpgradeByIDDB(ctx context.Context, id int64) (res *entities.Upgrade, err error) {
	row := r.db.QueryRowContext(ctx, getUpgradeByIDQuery, id)
	res, err = r.scanUpgrade(row)
	if errors.Is(err, sql.ErrNoRows) {
		return nil, apperror.ErrorNotFound("upgrade", "id", fmt.Sprint(id))
	} else if err != nil {
		return nil, err
	}

	return res, nil
}

func (r *upgradeRepository) GetUpgradeBySlugDB(ctx context.Context, slug string) (res *entities.Upgrade, err error) {
	row := r.db.QueryRowContext(ctx, getUpgradeBySlugQuery, slug)
	res, err = r.scanUpgrade(row)
	if errors.Is(err, sql.ErrNoRows) {
		return nil, apperror.ErrorNotFound("upgrade", "slug", slug)
	} else if err != nil {
		return nil, err
	}

	return res, nil
}

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

func (r *upgradeRepository) UpdateUpgradeDB(ctx context.Context, id int64, data entities.Upgrade) (err error) {
	now := helper.NowUTC()

	_, err = r.db.ExecContext(ctx, updateUpgradeQuery,
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
		id,
	)
	if database.IsDuplicateError(err) {
		return apperror.ErrorAlreadyExists("upgrade", "slug", data.Slug)
	} else if err != nil {
		return err
	}

	return nil
}

func (r *upgradeRepository) CountUpgradesDB(ctx context.Context) (totalRows int64, err error) {
	err = r.db.QueryRowContext(ctx, countUpgradesQuery).Scan(&totalRows)
	if err != nil {
		return 0, err
	}

	return totalRows, nil
}

func (r *upgradeRepository) scanUpgrade(rows *sql.Row) (res *entities.Upgrade, err error) {
	var data entities.Upgrade
	var effect entities.UpgradeEffect
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

	data.Effect = effect

	return &data, nil
}

func (r *upgradeRepository) GetUpgradesBySlugsDB(ctx context.Context, slugs []string) ([]entities.Upgrade, error) {
	var upgrades []entities.Upgrade
	rows, err := r.db.QueryContext(ctx, "SELECT id, slug FROM upgrades WHERE slug = ANY($1)", pq.Array(slugs))
	if err != nil {
		return nil, err
	}

	defer rows.Close()

	for rows.Next() {
		var upgrade entities.Upgrade

		if err := rows.Scan(
			&upgrade.ID,
			&upgrade.Slug,
		); err != nil {
			return nil, err
		}

		upgrades = append(upgrades, upgrade)
	}

	return upgrades, nil
}
