package usecase

import (
	"context"
	"database/sql"
	"github.com/winartodev/cat-cafe/pkg/helper"

	"github.com/winartodev/cat-cafe/internal/entities"
	"github.com/winartodev/cat-cafe/internal/repositories"
	"github.com/winartodev/cat-cafe/pkg/apperror"
)

type UserProgressionUseCase interface {
	InitializeUserProgression(ctx context.Context, userID int64, stageID int64, config *entities.GameStageConfig) (err error)
	LatestStageProgression(ctx context.Context) (res *entities.UserGameStageProgression, err error)
	DailyRewardProgression(ctx context.Context, userID int64) (res *entities.UserDailyReward, err error)
	GetActiveStageUpgrade(ctx context.Context, stageID int64) (res []entities.UserStageUpgrade, err error)
}

type userProgressionUseCase struct {
	userProgressionRepo repositories.UserProgressionRepository
	foodItemRepo        repositories.FoodItemRepository
}

func NewUserProgressionUseCase(userProgressionRepo repositories.UserProgressionRepository, foodItemRepo repositories.FoodItemRepository) UserProgressionUseCase {
	return &userProgressionUseCase{
		userProgressionRepo: userProgressionRepo,
		foodItemRepo:        foodItemRepo,
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

		err = userProgressionTx.MarkStageAsStartedDB(ctx, userID, stageID)
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
		return nil, apperror.ErrMissingKitchenConfig
	}

	kitchenProgression, err := repo.GetUserKitchenProgressDB(ctx, userID, stageID)
	if err != nil {
		return nil, err
	}

	if kitchenProgression != nil {
		return kitchenProgression, nil
	}

	stationLevels := make(map[string]entities.UserStationLevel)
	stationUpgrades := make(map[string]entities.UserStationUpgrade)
	var unlockedStations []string

	for _, station := range gameConfig.KitchenStations {
		if station.AutoUnlock {
			// Get Food Item for override check
			foodItem, err := u.foodItemRepo.GetFoodBySlugDB(ctx, station.FoodItemSlug)
			if err != nil {
				return nil, err
			}

			// Base stats for Level 1
			level1 := entities.UserStationLevel{
				Level:           1,
				Cost:            station.InitialCost,
				Profit:          station.InitialProfit,
				PreparationTime: station.CookingTime,
			}

			// Check override for level 1
			override, err := u.foodItemRepo.GetOverrideLevelByFoodItemIDAndLevelDB(ctx, foodItem.ID, 1)
			if err != nil {
				return nil, err
			}
			if override != nil {
				level1 = entities.UserStationLevel{
					Level:           override.Level,
					Cost:            override.Cost,
					Profit:          override.Profit,
					PreparationTime: override.PreparationTime,
				}
			}

			stationLevels[station.FoodItemSlug] = level1
			unlockedStations = append(unlockedStations, station.FoodItemSlug)
		} else {
			stationLevels[station.FoodItemSlug] = entities.UserStationLevel{
				Level:           0,
				Cost:            0,
				Profit:          0,
				PreparationTime: 0,
			}
		}

		stationUpgrades[station.FoodItemSlug] = entities.UserStationUpgrade{
			ProfitBonus:       1,
			ReduceCookingTime: 1,
			HelperCount:       0,
			CustomerCount:     0,
		}
	}

	newProgress := &entities.UserKitchenStageProgression{
		UserID:           userID,
		StageID:          stageID,
		StationLevels:    stationLevels,
		UnlockedStations: unlockedStations,
		StationUpgrades:  stationUpgrades,
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
		return nil, apperror.ErrMissingKitchenConfig
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

func (u *userProgressionUseCase) LatestStageProgression(ctx context.Context) (res *entities.UserGameStageProgression, err error) {
	userID, err := helper.GetUserIDFromContext(ctx)
	if err != nil {
		return nil, err
	}

	lastestProgression, err := u.userProgressionRepo.GetLatestGameStageProgressionDB(ctx, userID)
	if err != nil {
		return nil, err
	}

	if lastestProgression == nil {
		return nil, apperror.ErrStageNotFound
	}

	if lastestProgression.LastStartedAt == nil {
		return lastestProgression, apperror.ErrStageNotStarted
	}

	return lastestProgression, nil
}

func (u *userProgressionUseCase) GetActiveStageUpgrade(ctx context.Context, stageID int64) (res []entities.UserStageUpgrade, err error) {
	userID, err := helper.GetUserIDFromContext(ctx)
	if err != nil {
		return nil, err
	}

	upgrades, err := u.userProgressionRepo.GetCurrentStageUpgradeDB(ctx, userID, stageID)
	if err != nil {
		return nil, err
	}

	return upgrades, err
}
