package usecase

import (
	"context"
	"database/sql"
	"fmt"
	"github.com/winartodev/cat-cafe/internal/entities"
	"github.com/winartodev/cat-cafe/internal/repositories"
)

type UserProgressionUseCase interface {
	InitializeUserProgression(ctx context.Context, userID int64, stageID int64, config *entities.GameStageConfig) (err error)

	DailyRewardProgression(ctx context.Context, userID int64) (res *entities.UserDailyReward, err error)
}

type userProgressionUseCase struct {
	userProgressionRepo repositories.UserProgressionRepository
}

func NewUserProgressionUseCase(userProgressionRepo repositories.UserProgressionRepository) UserProgressionUseCase {
	return &userProgressionUseCase{
		userProgressionRepo: userProgressionRepo,
	}
}

func (u *userProgressionUseCase) InitializeUserProgression(ctx context.Context, userID int64, stageID int64, config *entities.GameStageConfig) (err error) {
	err = u.userProgressionRepo.WithUserProgressionTx(ctx, func(tx *sql.Tx) error {
		userProgressionTx := u.userProgressionRepo.WithTx(tx)

		_, err := u.getOrCreateKitchenProgress(ctx, userProgressionTx, userID, stageID, config)
		if err != nil {
			return err
		}

		_, err = u.getOrCreateKitchenPhaseProgression(ctx, userProgressionTx, userID, config.KitchenConfig.ID)
		if err != nil {
			return err
		}

		return nil
	})
	if err != nil {
		return err
	}

	return nil
}

func (u *userProgressionUseCase) getOrCreateKitchenProgress(
	ctx context.Context,
	repo repositories.UserProgressionRepository,
	userID int64,
	stageID int64,
	gameConfig *entities.GameStageConfig,
) (res *entities.UserKitchenStageProgression, err error) {
	if gameConfig == nil {
		return nil, fmt.Errorf("kitchen config is nil")
	}

	kitchenProgression, err := repo.GetUserKitchenProgressDB(ctx, userID, stageID)
	if err != nil {
		return nil, err
	}

	if kitchenProgression != nil {
		return kitchenProgression, nil
	}

	stationLevels := make(map[string]int64)
	var unlockedStations []string

	for _, station := range gameConfig.KitchenStations {
		if station.AutoUnlock {
			stationLevels[station.FoodItemSlug] = 1
			unlockedStations = append(unlockedStations, station.FoodItemSlug)
		} else {
			stationLevels[station.FoodItemSlug] = 0
		}
	}

	newProgress := &entities.UserKitchenStageProgression{
		UserID:           userID,
		StageID:          stageID,
		StationLevels:    stationLevels,
		UnlockedStations: unlockedStations,
	}

	err = repo.CreateUserKitchenProgressionDB(ctx, newProgress)
	if err != nil {
		return nil, err
	}

	return newProgress, nil
}

func (u *userProgressionUseCase) getOrCreateKitchenPhaseProgression(
	ctx context.Context,
	repo repositories.UserProgressionRepository,
	userID int64,
	kitchenConfigID int64,
) (res *entities.UserKitchenPhaseProgression, err error) {
	if kitchenConfigID == 0 {
		return nil, fmt.Errorf("kitchen config is nil")
	}

	kitchenPhaseProgression, err := repo.GetUserKitchenPhaseProgressionDB(ctx, userID, kitchenConfigID)
	if err != nil {
		return nil, err
	}

	if kitchenPhaseProgression != nil {
		return kitchenPhaseProgression, nil
	}

	newKitchenPhaseProgression := &entities.UserKitchenPhaseProgression{
		UserID:          userID,
		KitchenConfigID: kitchenConfigID,
		CurrentPhase:    1,
		CompletedPhases: make([]int64, 0),
	}

	err = repo.CreateUserKitchenPhaseProgressionDB(ctx, newKitchenPhaseProgression)
	if err != nil {
		return nil, err
	}

	return newKitchenPhaseProgression, nil
}

func (u *userProgressionUseCase) DailyRewardProgression(ctx context.Context, userID int64) (res *entities.UserDailyReward, err error) {
	return u.userProgressionRepo.GetUserDailyRewardByIDDB(ctx, userID)
}
