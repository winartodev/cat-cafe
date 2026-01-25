package repositories

import (
	"context"
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/lib/pq"
	_ "github.com/lib/pq"
	"github.com/winartodev/cat-cafe/internal/entities"
	"github.com/winartodev/cat-cafe/pkg/apperror"
	"github.com/winartodev/cat-cafe/pkg/database"
	"github.com/winartodev/cat-cafe/pkg/helper"
)

type GameStageRepository interface {
	WithTx(tx *sql.Tx) GameStageRepository
	GameStageWithTx(ctx context.Context, fn func(tx *sql.Tx) error) error
	CreateGameStageWithTxDB(ctx context.Context, data *entities.GameStage) (*int64, error)
	UpdateGameStageWithTxDB(ctx context.Context, data *entities.GameStage) error
	GetGameStagesDB(ctx context.Context, limit, offset int) (res []entities.GameStage, totalRows int64, err error)
	GetGameStageByIDDB(ctx context.Context, id int64) (*entities.GameStage, error)
	GetGameStageBySlugDB(ctx context.Context, slug string) (*entities.GameStage, error)

	GetGameConfigByIDDB(ctx context.Context, stageID int64) (*entities.GameStageConfig, error)

	GetActiveGameStagesDB(ctx context.Context) ([]entities.GameStage, error)
}

type gameStageRepository struct {
	BaseRepository
}

func NewGameStageRepository(db *sql.DB) GameStageRepository {
	return &gameStageRepository{
		BaseRepository: BaseRepository{
			db:   db,
			pool: db,
		},
	}
}

func (r *gameStageRepository) WithTx(tx *sql.Tx) GameStageRepository {
	if tx == nil {
		return r
	}

	return &gameStageRepository{
		BaseRepository{
			db:   tx,
			pool: r.pool,
		},
	}
}

func (r *gameStageRepository) GameStageWithTx(ctx context.Context, fn func(tx *sql.Tx) error) error {
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

func (r *gameStageRepository) CreateGameStageWithTxDB(ctx context.Context, data *entities.GameStage) (*int64, error) {
	now := helper.NowUTC()
	var id int64

	err := r.db.QueryRowContext(ctx, insertIntoGameStageQuery,
		data.Slug,
		data.Name,
		data.StartingCoin,
		data.StagePrize,
		data.IsActive,
		data.Sequence,
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

func (r *gameStageRepository) UpdateGameStageWithTxDB(ctx context.Context, data *entities.GameStage) error {
	now := helper.NowUTC()

	_, err := r.db.ExecContext(ctx, updateGameStageQuery,
		data.ID,
		data.Name,
		data.StartingCoin,
		data.StagePrize,
		now,
		data.IsActive,
		data.Sequence,
	)
	if err != nil {
		return err
	}

	return nil
}

func (r *gameStageRepository) GetGameStageByIDDB(ctx context.Context, id int64) (*entities.GameStage, error) {
	row := r.db.QueryRowContext(ctx, getGameStageByIDQuery, id)
	return r.scanGameStageRow(row)
}

func (r *gameStageRepository) GetGameStageBySlugDB(ctx context.Context, slug string) (*entities.GameStage, error) {
	row := r.db.QueryRowContext(ctx, getGameStageBySlugQuery, slug)
	return r.scanGameStageRow(row)
}

func (r *gameStageRepository) GetGameStagesDB(ctx context.Context, limit, offset int) (res []entities.GameStage, totalRows int64, err error) {
	rows, err := r.db.QueryContext(ctx, getGameStageQuery, limit, offset)
	if err != nil {
		return nil, 0, err
	}

	for rows.Next() {
		var gameStage entities.GameStage
		err = rows.Scan(
			&gameStage.ID,
			&gameStage.Slug,
			&gameStage.Name,
			&gameStage.StartingCoin,
			&gameStage.StagePrize,
			&gameStage.IsActive,
			&gameStage.Sequence,
		)
		if err != nil {
			return nil, 0, err
		}

		res = append(res, gameStage)
	}

	err = r.db.QueryRowContext(ctx, countGameStagesQuery).Scan(&totalRows)
	if err != nil {
		return nil, 0, err
	}

	return res, totalRows, nil
}

func (r *gameStageRepository) GetGameConfigByIDDB(ctx context.Context, stageID int64) (*entities.GameStageConfig, error) {
	row := r.db.QueryRowContext(ctx, getGameStageConfig, stageID)
	return r.scanGameStageConfigRow(row)
}

func (r *gameStageRepository) scanGameStageRow(row *sql.Row) (*entities.GameStage, error) {
	var gameStage entities.GameStage
	err := row.Scan(
		&gameStage.ID,
		&gameStage.Slug,
		&gameStage.Name,
		&gameStage.StartingCoin,
		&gameStage.StagePrize,
		&gameStage.IsActive,
		&gameStage.Sequence,
	)
	if errors.Is(err, sql.ErrNoRows) {
		return nil, nil
	} else if err != nil {
		return nil, err
	}

	return &gameStage, nil
}

func (r *gameStageRepository) scanGameStageConfigRow(row *sql.Row) (*entities.GameStageConfig, error) {
	var gameStageConfig entities.GameStageConfig
	var scc entities.StageCustomerConfig
	var ssc entities.StageStaffConfig
	var skc entities.StageKitchenConfig
	var scc2 entities.StageCameraConfig
	var rewardsJSON []byte
	var kitchenStationJSON []byte

	err := row.Scan(
		&scc.CustomerSpawnTime,
		&scc.MaxCustomerOrderCount,
		&scc.MaxCustomerOrderVariant,
		&scc.StartingOrderTableCount,
		&ssc.StartingStaffManager,
		&ssc.StartingStaffHelper,
		&skc.ID,
		&skc.MaxLevel,
		&skc.UpgradeProfitMultiply,
		&skc.UpgradeCostMultiply,
		pq.Array(&skc.TransitionPhaseLevels),
		pq.Array(&skc.PhaseProfitMultipliers),
		pq.Array(&skc.PhaseUpgradeCostMultipliers),
		pq.Array(&skc.TableCountPerPhases),
		&rewardsJSON,
		&kitchenStationJSON,
		&scc2.ZoomSize,
		&scc2.MaxBoundX,
		&scc2.MinBoundX,
		&scc2.MinBoundY,
		&scc2.MaxBoundY,
	)
	if errors.Is(err, sql.ErrNoRows) {
		return nil, nil
	} else if err != nil {
		return nil, err
	}

	if len(rewardsJSON) > 0 {
		if err := json.Unmarshal(rewardsJSON, &gameStageConfig.KitchenPhaseReward); err != nil {
			return nil, fmt.Errorf("parse rewards into map: %w", err)
		}
	}

	if len(kitchenStationJSON) > 0 {
		if err := json.Unmarshal(kitchenStationJSON, &gameStageConfig.KitchenStations); err != nil {
			return nil, fmt.Errorf("parse kithcen stations into map: %w", err)
		}
	}

	gameStageConfig.CustomerConfig = &scc
	gameStageConfig.StaffConfig = &ssc
	gameStageConfig.KitchenConfig = &skc
	gameStageConfig.CameraConfig = &scc2

	return &gameStageConfig, nil
}

func (r *gameStageRepository) GetActiveGameStagesDB(ctx context.Context) (res []entities.GameStage, err error) {
	rows, err := r.db.QueryContext(ctx, getActiveGameStagesQuery)
	if err != nil {
		return nil, err
	}

	for rows.Next() {
		var gameStage entities.GameStage
		err = rows.Scan(
			&gameStage.ID,
			&gameStage.Slug,
			&gameStage.Name,
			&gameStage.Sequence,
		)

		res = append(res, gameStage)
	}

	return res, nil
}
