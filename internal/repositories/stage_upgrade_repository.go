package repositories

import (
	"context"
	"database/sql"

	"github.com/winartodev/cat-cafe/internal/entities"
	"github.com/winartodev/cat-cafe/pkg/apperror"
	"github.com/winartodev/cat-cafe/pkg/database"
	"github.com/winartodev/cat-cafe/pkg/helper"
)

type StageUpgradeRepository interface {
	WithTx(tx *sql.Tx) StageUpgradeRepository
	StageUpgradeWithTx(ctx context.Context, fn func(tx *sql.Tx) error) (err error)

	BulkCreateStageUpgradesDB(ctx context.Context, data []entities.StageUpgrade) (err error)
	GetStageUpgrades(ctx context.Context, stageID int64, limit, offset int) ([]entities.Upgrade, error)
	CountStageUpgrades(ctx context.Context, stageID int64) (int64, error)
}

type stageUpgradeRepository struct {
	BaseRepository
}

func NewStageUpgradeRepository(db *sql.DB) StageUpgradeRepository {
	return &stageUpgradeRepository{
		BaseRepository: BaseRepository{
			db:   db,
			pool: db,
		},
	}
}

func (r *stageUpgradeRepository) WithTx(tx *sql.Tx) StageUpgradeRepository {
	if tx == nil {
		return r
	}

	return &stageUpgradeRepository{
		BaseRepository: BaseRepository{
			db:   tx,
			pool: r.pool,
		},
	}
}

func (r *stageUpgradeRepository) StageUpgradeWithTx(ctx context.Context, fn func(tx *sql.Tx) error) (err error) {
	tx, err := r.pool.BeginTx(ctx, nil)
	if err != nil {
		return err
	}

	defer tx.Rollback()

	if err := fn(tx); err != nil {
		return err
	}

	return tx.Commit()
}

func (r *stageUpgradeRepository) BulkCreateStageUpgradesDB(ctx context.Context, data []entities.StageUpgrade) (err error) {
	if len(data) == 0 {
		return nil
	}

	numFields := 4
	onConflict := " ON CONFLICT (game_stage_id, upgrade_id) DO NOTHING"
	queryString := r.BuildBulkInsertQuery(insertStageUpgradeQuery, len(data), numFields, onConflict)

	args := make([]interface{}, 0, len(data)*numFields)
	now := helper.NowUTC()

	for _, item := range data {
		args = append(args,
			item.StageID,
			item.UpgradeID,
			now,
			now,
		)
	}

	_, err = r.db.ExecContext(ctx, queryString, args...)
	if database.IsDuplicateError(err) {
		return apperror.ErrConflict
	} else if err != nil {
		return err
	}

	return nil
}

func (r *stageUpgradeRepository) GetStageUpgrades(ctx context.Context, stageID int64, limit, offset int) ([]entities.Upgrade, error) {
	query := getStageUpgradeQuery + ` ORDER BY u.sequence ASC LIMIT $2 OFFSET $3`

	rows, err := r.db.QueryContext(ctx, query, stageID, limit, offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var upgrades []entities.Upgrade
	for rows.Next() {
		var u entities.Upgrade
		var e entities.UpgradeEffect
		var stageSlugIgnored string

		err := rows.Scan(
			&stageSlugIgnored,
			&u.ID,
			&u.Slug,
			&u.Name,
			&u.Description,
			&u.Cost,
			&u.CostType,
			&e.Type,
			&e.Value,
			&e.Unit,
			&e.Target,
			&e.TargetID,
			&u.IsActive,
			&u.Sequence,
			&u.CreatedAt,
			&u.UpdatedAt,
		)
		if err != nil {
			return nil, err
		}
		u.Effect = e
		upgrades = append(upgrades, u)
	}

	return upgrades, nil
}

func (r *stageUpgradeRepository) CountStageUpgrades(ctx context.Context, stageID int64) (int64, error) {
	var count int64

	err := r.db.QueryRowContext(ctx, getStageUpgradeCountQuery, stageID).Scan(&count)
	if err != nil {
		return 0, err
	}

	return count, nil
}
