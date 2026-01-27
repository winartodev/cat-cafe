package repositories

import (
	"context"
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"time"

	"github.com/redis/go-redis/v9"
	"github.com/winartodev/cat-cafe/internal/entities"
	"github.com/winartodev/cat-cafe/pkg/apperror"
	"github.com/winartodev/cat-cafe/pkg/database"
	"github.com/winartodev/cat-cafe/pkg/helper"
)

const (
	userDailyRewardRedisKey = "user_daily_reward:%d"
)

type UserProgressionRepository interface {
	WithTx(tx *sql.Tx) UserProgressionRepository
	WithUserProgressionTx(ctx context.Context, fn func(tx *sql.Tx) error) error

	GetUserDailyRewardByIDDB(ctx context.Context, id int64) (res *entities.UserDailyReward, err error)
	UpsertDailyRewardProgressionDB(ctx context.Context, userID int64, longestStreak int64, currentStreak int64, lastClaim time.Time) (err error)

	GetGameStageProgressionDB(ctx context.Context, userID int64, stageID int64) (res *entities.UserGameStageProgression, err error)
	GetLatestGameStageProgressionDB(ctx context.Context, userID int64) (res *entities.UserGameStageProgression, err error)
	CreateGameStageProgressionDB(ctx context.Context, userID int64, stageID int64) (*int64, error)
	CheckStageProgressionExistsDB(ctx context.Context, userID int64, stageID int64) (bool, error)
	MarkStageAsCompleteDB(ctx context.Context, userID int64, stageID int64) error

	CreateUserKitchenProgressionDB(ctx context.Context, data *entities.UserKitchenStageProgression) (err error)
	UpdateUserKitchenProgressDB(ctx context.Context, userID int64, stageID int64, progress *entities.UserKitchenStageProgression) (err error)
	GetUserKitchenProgressDB(ctx context.Context, userID int64, stageID int64) (res *entities.UserKitchenStageProgression, err error)

	GetUserKitchenPhaseProgressionDB(ctx context.Context, userID int64, kitchenConfigID int64) (res *entities.UserKitchenPhaseProgression, err error)
	CreateUserKitchenPhaseProgressionDB(ctx context.Context, data *entities.UserKitchenPhaseProgression) (err error)
	UpdateUserKitchenPhaseProgressionDB(ctx context.Context, userID int64, kitchenConfigID int64, data *entities.UserKitchenPhaseProgression) (err error)

	IsPhaseRewardAlreadyClaimedDB(ctx context.Context, userID int64, kitchenConfigID int64, phaseNumber int64, rewardID int64) (claimed bool, err error)
	CreateUserKitchenClaimRewardDB(ctx context.Context, data entities.UserKitchenPhaseRewardClaim) (err error)

	//GetUserDailyRewardRedis(ctx context.Context, userID int64) (res *entities.UserDailyReward, err error)
	//SetUserDailyRewardRedis(ctx context.Context, userID int64, progress *entities.UserDailyReward, ttl time.Duration) (err error)
	//DeleteUserDailyRewardRedis(ctx context.Context, userID int64) error
}

type userProgressionRepository struct {
	BaseRepository
}

func NewUserProgressionRepository(db *sql.DB, redis *redis.Client) UserProgressionRepository {
	return &userProgressionRepository{
		BaseRepository{
			db:    db,
			pool:  db,
			redis: redis,
		},
	}
}

func (r *userProgressionRepository) WithTx(tx *sql.Tx) UserProgressionRepository {
	if tx == nil {
		return r
	}

	return &userProgressionRepository{
		BaseRepository{
			db:    tx,
			pool:  r.pool,
			redis: r.redis,
		},
	}
}

func (r *userProgressionRepository) WithUserProgressionTx(ctx context.Context, fn func(tx *sql.Tx) error) error {
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

func (r *userProgressionRepository) GetUserDailyRewardByIDDB(ctx context.Context, id int64) (res *entities.UserDailyReward, err error) {
	var data entities.UserDailyReward
	err = r.db.QueryRowContext(ctx, getUserDailyRewardProgressQuery, id).Scan(
		&data.ID,
		&data.UserID,
		&data.LongestStreak,
		&data.CurrentStreak,
		&data.LastClaimDate,
	)

	if errors.Is(err, sql.ErrNoRows) {
		return nil, nil
	} else if err != nil {
		return nil, err
	}

	return &data, err
}

func (r *userProgressionRepository) UpsertDailyRewardProgressionDB(ctx context.Context, userID int64, longestStreak int64, currentStreak int64, lastClaim time.Time) (err error) {
	now := time.Now()

	_, err = r.db.ExecContext(ctx, upsertUserDailyRewardProgressQuery,
		userID,
		longestStreak,
		currentStreak,
		lastClaim,
		now,
	)

	return err
}

func (r *userProgressionRepository) GetUserDailyRewardRedis(ctx context.Context, userID int64) (res *entities.UserDailyReward, err error) {
	key := r.userDailyRewardKey(userID)
	val, err := r.redis.Get(ctx, key).Result()

	if errors.Is(err, redis.Nil) {
		return nil, nil
	} else if err != nil {
		return nil, err
	}

	var progress entities.UserDailyReward
	err = json.Unmarshal([]byte(val), &progress)
	if err != nil {
		return nil, err
	}

	return &progress, nil
}

func (r *userProgressionRepository) SetUserDailyRewardRedis(ctx context.Context, userID int64, progress *entities.UserDailyReward, ttl time.Duration) (err error) {
	key := r.userDailyRewardKey(userID)
	data, err := json.Marshal(progress)
	if err != nil {
		return err
	}

	return r.redis.Set(ctx, key, data, ttl).Err()
}

func (r *userProgressionRepository) DeleteUserDailyRewardRedis(ctx context.Context, userID int64) error {
	key := r.userDailyRewardKey(userID)
	return r.redis.Del(ctx, key).Err()
}

func (r *userProgressionRepository) userDailyRewardKey(userID int64) string {
	return fmt.Sprintf(userDailyRewardRedisKey, userID)
}

func (r *userProgressionRepository) GetActiveGameStagesDB(ctx context.Context) (res []entities.GameStage, err error) {
	rows, err := r.db.QueryContext(ctx, getActiveGameStagesQuery)
	if err != nil {
		return nil, err
	}

	for rows.Next() {
		var data entities.GameStage

		err = rows.Scan(
			&data.Slug,
			&data.Name,
			&data.Sequence,
		)
		res = append(res, data)
	}

	return res, nil
}

func (r *userProgressionRepository) GetGameStageProgressionDB(ctx context.Context, userID int64, stageID int64) (res *entities.UserGameStageProgression, err error) {
	row := r.db.QueryRowContext(ctx, getGameStageProgressionQuery, userID, stageID)
	var data entities.UserGameStageProgression
	err = row.Scan(
		&data.ID,
		&data.UserID,
		&data.StageID,
		&data.IsComplete,
		&data.CompletedAt,
	)
	if errors.Is(err, sql.ErrNoRows) {
		return nil, nil
	} else if err != nil {
		return nil, err
	}

	return &data, nil
}

func (r *userProgressionRepository) CreateGameStageProgressionDB(ctx context.Context, userID int64, stageID int64) (*int64, error) {
	var id int64
	err := r.db.QueryRowContext(ctx, insertGameStageProgressionQuery,
		userID,
		stageID,
	).Scan(&id)
	if database.IsDuplicateError(err) {
		return nil, apperror.ErrConflict
	} else if err != nil {
		return nil, err
	}

	return &id, nil
}

func (r *userProgressionRepository) CheckStageProgressionExistsDB(ctx context.Context, userID int64, stageID int64) (bool, error) {
	var exists bool

	err := r.db.QueryRowContext(ctx, checkStageProgressionExistsQuery, userID, stageID).Scan(&exists)
	if err != nil {
		return false, err
	}
	return exists, nil
}

func (r *userProgressionRepository) MarkStageAsCompleteDB(ctx context.Context, userID int64, stageID int64) error {
	result, err := r.db.ExecContext(ctx, markStageAsComplete, helper.NowUTC(), userID, stageID)
	if err != nil {
		return err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}

	if rowsAffected == 0 {
		return apperror.ErrNoUpdateRecord
	}

	return nil
}

func (r *userProgressionRepository) GetLatestGameStageProgressionDB(ctx context.Context, userID int64) (res *entities.UserGameStageProgression, err error) {
	row := r.db.QueryRowContext(ctx, getLatestGameStageProgressionQuery, userID)
	var data entities.UserGameStageProgression
	err = row.Scan(
		&data.ID,
		&data.UserID,
		&data.StageID,
		&data.IsComplete,
		&data.CompletedAt,
	)
	if errors.Is(err, sql.ErrNoRows) {
		return nil, nil
	} else if err != nil {
		return nil, err
	}

	return &data, nil
}

func (r *userProgressionRepository) GetUserKitchenProgressDB(ctx context.Context, userID int64, stageID int64) (res *entities.UserKitchenStageProgression, err error) {
	var progression entities.UserKitchenStageProgression
	var stationLevelsJSON []byte
	var unlockedStationsJSON []byte

	err = r.db.QueryRowContext(ctx,
		getUserKitchenProgressQuery,
		userID,
		stageID,
	).Scan(
		&progression.ID,
		&progression.UserID,
		&progression.StageID,
		&stationLevelsJSON,
		&unlockedStationsJSON,
	)
	if errors.Is(err, sql.ErrNoRows) {
		return nil, nil
	} else if err != nil {
		return nil, err
	}

	if err := json.Unmarshal(stationLevelsJSON, &progression.StationLevels); err != nil {
		return nil, err
	}

	if err := json.Unmarshal(unlockedStationsJSON, &progression.UnlockedStations); err != nil {
		return nil, err
	}

	return &progression, nil
}

func (r *userProgressionRepository) CreateUserKitchenProgressionDB(ctx context.Context, data *entities.UserKitchenStageProgression) (err error) {
	stationLevelsJSON, err := json.Marshal(data.StationLevels)
	if err != nil {
		return err
	}

	unlockedStationsJSON, err := json.Marshal(data.UnlockedStations)
	if err != nil {
		return err
	}

	_, err = r.db.ExecContext(ctx, insertUserKitchenProgressQuery,
		data.UserID,
		data.StageID,
		stationLevelsJSON,
		unlockedStationsJSON,
	)
	if err != nil {
		return err
	}

	return nil
}

func (r *userProgressionRepository) UpdateUserKitchenProgressDB(ctx context.Context, userID int64, stageID int64, progress *entities.UserKitchenStageProgression) (err error) {
	stationLevelsJSON, err := json.Marshal(progress.StationLevels)
	if err != nil {
		return err
	}

	unlockedStationsJSON, err := json.Marshal(progress.UnlockedStations)
	if err != nil {
		return err
	}

	result, err := r.db.ExecContext(
		ctx,
		updateUserKitchenProgressQuery,
		userID,
		stageID,
		stationLevelsJSON,
		unlockedStationsJSON,
		helper.NowUTC(),
	)
	if err != nil {
		return err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}

	if rowsAffected == 0 {
		return apperror.ErrNoUpdateRecord
	}

	return nil

}

func (r *userProgressionRepository) GetUserKitchenPhaseProgressionDB(ctx context.Context, userID int64, kitchenConfigID int64) (res *entities.UserKitchenPhaseProgression, err error) {
	var data entities.UserKitchenPhaseProgression
	var completedPhaseJSON []byte
	err = r.db.QueryRowContext(ctx, getUserKitchenPhaseProgressQuery,
		userID,
		kitchenConfigID,
	).Scan(
		&data.UserID,
		&data.KitchenConfigID,
		&data.CurrentPhase,
		&completedPhaseJSON,
	)
	if errors.Is(err, sql.ErrNoRows) {
		return nil, nil
	} else if err != nil {
		return nil, err
	}

	return &data, nil
}

func (r *userProgressionRepository) CreateUserKitchenPhaseProgressionDB(ctx context.Context, data *entities.UserKitchenPhaseProgression) (err error) {
	completedPhasesJSON, err := json.Marshal(data.CompletedPhases)
	if err != nil {
		return err
	}

	_, err = r.db.ExecContext(ctx, insertUserKitchenPhaseProgressQuery,
		data.UserID,
		data.KitchenConfigID,
		data.CurrentPhase,
		completedPhasesJSON,
	)
	if err != nil {
		return err
	}

	return nil
}

func (r *userProgressionRepository) UpdateUserKitchenPhaseProgressionDB(ctx context.Context, userID int64, kitchenConfigID int64, data *entities.UserKitchenPhaseProgression) (err error) {
	completedPhasesJSON, err := json.Marshal(data.CompletedPhases)
	if err != nil {
		return fmt.Errorf("failed to marshal completed_phases: %w", err)
	}

	now := helper.NowUTC()
	result, err := r.db.ExecContext(
		ctx,
		updateUserKitchenPhaseProgressQuery,
		userID,
		kitchenConfigID,
		data.CurrentPhase,
		completedPhasesJSON,
		now,
	)
	if err != nil {
		return nil
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}

	if rowsAffected == 0 {
		return err
	}

	return nil

}

func (r *userProgressionRepository) IsPhaseRewardAlreadyClaimedDB(ctx context.Context, userID int64, kitchenConfigID int64, phaseNumber int64, rewardID int64) (res bool, err error) {
	var claimed int64
	err = r.db.QueryRowContext(ctx, getUserPhaseRewardClaimedQuery,
		userID,
		kitchenConfigID,
		phaseNumber,
		rewardID,
	).Scan(&claimed)
	if errors.Is(err, sql.ErrNoRows) {
		return false, nil
	} else if err != nil {
		return false, err
	}

	fmt.Println("claimed ", claimed)

	return claimed == 1, nil
}

func (r *userProgressionRepository) CreateUserKitchenClaimRewardDB(ctx context.Context, data entities.UserKitchenPhaseRewardClaim) (err error) {
	now := helper.NowUTC()
	_, err = r.db.ExecContext(ctx, createUserKitchenPhaseClaimReward,
		data.UserID,
		data.KitchenConfigID,
		data.RewardID,
		data.CurrentPhase,
		now,
	)
	if database.IsDuplicateError(err) {
		return apperror.ErrConflict
	} else if err != nil {
		return err
	}

	return nil
}
