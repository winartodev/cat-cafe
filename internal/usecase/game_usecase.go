package usecase

import (
	"context"
	"database/sql"
	"errors"
	"github.com/winartodev/cat-cafe/internal/dto"
	"github.com/winartodev/cat-cafe/internal/entities"
	"github.com/winartodev/cat-cafe/internal/repositories"
	"github.com/winartodev/cat-cafe/pkg/apperror"
	"github.com/winartodev/cat-cafe/pkg/helper"
)

const (
	upgradeTargetFood = "food"
)

// GameUseCase is used for interaction with player
type GameUseCase interface {
	UpdateUserBalance(ctx context.Context, coinEarned int64) (res *dto.UserBalanceResponse, err error)
	GetUserGameData(ctx context.Context) (res *entities.Game, err error)

	GetGameStages(ctx context.Context) (res []entities.UserGameStage, nextStage *entities.UserNextGameStageInfo, err error)
	GetCurrentGameStage(ctx context.Context) (stage *entities.GameStage, config *entities.GameStageConfig, nextStage *entities.UserNextGameStageInfo, err error)
	StartGameStage(ctx context.Context, slug string) (stage *entities.GameStage, config *entities.GameStageConfig, nextStage *entities.UserNextGameStageInfo, err error)
	CompleteGameStage(ctx context.Context, slug string) (err error)

	UnlockKitchenStation(ctx context.Context, slug string) (res *entities.UnlockKitchenStation, err error)
	UpgradeKitchenStation(ctx context.Context, slug string) (res *entities.UpgradeKitchenStation, err error)

	GetStageUpgrades(ctx context.Context) (res []entities.UserStageUpgrade, err error)
	PurchaseStageUpgrade(ctx context.Context, slug string) (res *entities.Upgrade, err error)
}

type gameUseCase struct {
	userUseCase            UserUseCase
	userProgressionUseCase UserProgressionUseCase

	userProgressionRepo repositories.UserProgressionRepository
	userRepo            repositories.UserRepository
	gameStageRepo       repositories.GameStageRepository
	foodItemRepo        repositories.FoodItemRepository
	kitchenStationRepo  repositories.KitchenStationRepository
	kitchenConfigRepo   repositories.StageKitchenConfigRepository
	rewardRepo          repositories.RewardRepository
	stageUpgradeRepo    repositories.StageUpgradeRepository
}

func NewGameUseCase(
	userUc UserUseCase,
	userProgressionUC UserProgressionUseCase,
	userRepo repositories.UserRepository,
	userProgressionRepo repositories.UserProgressionRepository,
	gameStageRepo repositories.GameStageRepository,
	foodItemRepo repositories.FoodItemRepository,
	kitchenStationRepo repositories.KitchenStationRepository,
	kitchenConfigRepo repositories.StageKitchenConfigRepository,
	rewardRepo repositories.RewardRepository,
	stageUpgradeRepo repositories.StageUpgradeRepository,
) GameUseCase {
	return &gameUseCase{
		userUseCase:            userUc,
		userProgressionUseCase: userProgressionUC,
		userRepo:               userRepo,
		userProgressionRepo:    userProgressionRepo,
		gameStageRepo:          gameStageRepo,
		foodItemRepo:           foodItemRepo,
		kitchenStationRepo:     kitchenStationRepo,
		kitchenConfigRepo:      kitchenConfigRepo,
		rewardRepo:             rewardRepo,
		stageUpgradeRepo:       stageUpgradeRepo,
	}
}

func (g *gameUseCase) UpdateUserBalance(ctx context.Context, coinEarned int64) (res *dto.UserBalanceResponse, err error) {
	userID, err := helper.GetUserIDFromContext(ctx)
	if err != nil {
		return nil, err
	}

	// TODO: should we validate earning rate ???

	err = g.userRepo.BalanceWithTx(ctx, func(tx *sql.Tx) error {
		txRepo := g.userRepo.WithTx(tx)

		if err := txRepo.UpdateUserBalanceWithTx(ctx, userID, entities.BalanceTypeCoin, coinEarned); err != nil {
			return err
		}

		if err := txRepo.UpdateLastSyncBalanceWithTx(ctx, userID, helper.NowUTC()); err != nil {
			return err
		}

		return nil
	})
	if err != nil {
		return nil, err
	}

	_ = g.userRepo.DeleteUserRedis(ctx, userID)

	user, err := g.userUseCase.GetUserByID(ctx, userID)
	if err != nil {
		return nil, err
	}

	return &dto.UserBalanceResponse{
		Coin: user.UserBalance.Coin,
		Gem:  user.UserBalance.Gem,
	}, nil
}

func (g *gameUseCase) GetUserGameData(ctx context.Context) (res *entities.Game, err error) {
	userID, err := helper.GetUserIDFromContext(ctx)
	if err != nil {
		return nil, err
	}

	isDailyRewardAvailable, err := g.userUseCase.IsDailyRewardAvailable(ctx, userID)
	if err != nil {
		return nil, err
	}

	return &entities.Game{
		DailyRewardAvailable: isDailyRewardAvailable,
	}, nil
}

func (g *gameUseCase) GetGameStages(ctx context.Context) (res []entities.UserGameStage, nextStage *entities.UserNextGameStageInfo, err error) {
	userID, err := helper.GetUserIDFromContext(ctx)
	if err != nil {
		return nil, nil, err
	}

	gameStages, err := g.gameStageRepo.GetActiveGameStagesDB(ctx)
	if err != nil || len(gameStages) == 0 {
		return nil, nil, err
	}

	latestProgress, err := g.userProgressionUseCase.LatestStageProgression(ctx)
	if err != nil && !errors.Is(err, apperror.ErrStageNotFound) && !errors.Is(err, apperror.ErrStageNotStarted) {
		return nil, nil, err
	}

	// Create new progression if player doesn't have any
	if latestProgress == nil {
		firstStage := gameStages[0]
		_, err = g.userProgressionRepo.CreateGameStageProgressionDB(ctx, userID, firstStage.ID)
		if err != nil {
			return nil, nil, err
		}

		return g.GetGameStages(ctx)
	}

	stages, nextStage := g.mapToUserGameStage(gameStages, latestProgress)

	return stages, nextStage, nil
}

func (g *gameUseCase) GetCurrentGameStage(ctx context.Context) (stage *entities.GameStage, config *entities.GameStageConfig, nextStage *entities.UserNextGameStageInfo, err error) {
	userID, err := helper.GetUserIDFromContext(ctx)
	if err != nil {
		return nil, nil, nil, err
	}

	lastestStage, err := g.userProgressionUseCase.LatestStageProgression(ctx)
	if err != nil {
		return nil, nil, nil, err
	}

	stageID := lastestStage.StageID

	stage, err = g.gameStageRepo.GetGameStageByIDDB(ctx, stageID)
	if err != nil {
		return nil, nil, nil, err
	}

	config, err = g.gameStageRepo.GetGameConfigByIDDB(ctx, stageID)
	if err != nil {
		return nil, nil, nil, err
	}

	if err = g.gatherUserProgressionData(ctx, userID, stageID, config); err != nil {
		return nil, nil, nil, err
	}

	return stage, config, nextStage, nil
}

func (g *gameUseCase) StartGameStage(ctx context.Context, slug string) (stage *entities.GameStage, config *entities.GameStageConfig, nextStage *entities.UserNextGameStageInfo, err error) {
	userID, err := helper.GetUserIDFromContext(ctx)
	if err != nil {
		return nil, nil, nil, err
	}

	stage, err = g.gameStageRepo.GetGameStageBySlugDB(ctx, slug)
	if err != nil {
		return nil, nil, nil, err
	}

	// Get user's available stages
	userStages, nextStage, err := g.GetGameStages(ctx)
	if err != nil {
		return nil, nil, nil, err
	}

	// Check if user can access this stage
	var canAccess bool
	for _, userStage := range userStages {
		if userStage.Slug == slug {
			// Only allow if status is Current or Complete (if you allow replay)
			if userStage.Status == entities.GSStatusAvailable || userStage.Status == entities.GSStatusComplete {
				canAccess = true
			}
			break
		}
	}

	if !canAccess {
		return nil, nil, nil, apperror.ErrStageLocked
	}

	config, err = g.gameStageRepo.GetGameConfigByIDDB(ctx, stage.ID)
	if err != nil {
		return nil, nil, nil, err
	}

	err = g.userProgressionUseCase.InitializeUserProgression(ctx, userID, stage.ID, config)
	if err != nil {
		return nil, nil, nil, err
	}

	err = g.gatherUserProgressionData(ctx, userID, stage.ID, config)
	if err != nil {
		return nil, nil, nil, err
	}

	return stage, config, nextStage, nil
}

func (g *gameUseCase) gatherUserProgressionData(ctx context.Context, userID int64, stageID int64, config *entities.GameStageConfig) (err error) {
	if config == nil {
		config, err = g.gameStageRepo.GetGameConfigByIDDB(ctx, stageID)
		if err != nil {
			return err
		}
	}

	// Get user kitchen progression to include in response
	userKitchenProgress, err := g.userProgressionRepo.GetUserKitchenProgressDB(ctx, userID, stageID)
	if err != nil {
		return err
	}

	helper.PrettyPrint(userKitchenProgress == nil)

	// Calculate next level stats for all unlocked stations
	if userKitchenProgress != nil && len(userKitchenProgress.UnlockedStations) > 0 {
		userKitchenProgress.NextLevelStats = make(map[string]entities.UserStationLevel)

		for _, slug := range userKitchenProgress.UnlockedStations {
			currentLevel, exists := userKitchenProgress.StationLevels[slug]
			if !exists || currentLevel.Level == 0 {
				continue
			}

			// Skip if at max level
			if currentLevel.Level >= config.KitchenConfig.MaxLevel {
				continue
			}

			nextLevel := currentLevel.Level + 1
			phaseInfo := g.calculatePhaseInfo(nextLevel, config.KitchenConfig)

			// Calculate next cost and profit
			nextCost := g.calculateUpgradeCost(
				currentLevel.Cost,
				currentLevel.Level,
				config.KitchenConfig,
				phaseInfo.CurrentPhase,
			)
			nextProfit := g.calculateProfit(
				currentLevel.Profit,
				currentLevel.Level,
				config.KitchenConfig,
				phaseInfo.CurrentPhase,
				0,
			)

			// Check for override
			foodItem, err := g.foodItemRepo.GetFoodBySlugDB(ctx, slug)
			if err == nil && foodItem != nil {
				override, err := g.foodItemRepo.GetOverrideLevelByFoodItemIDAndLevelDB(ctx, foodItem.ID, int(nextLevel))
				if err == nil && override != nil {
					nextCost = override.Cost
					nextProfit = override.Profit
				}
			}

			phaseRewards, err := g.kitchenConfigRepo.GetKitchenCompletionRewardByPhaseNumberDB(ctx, config.KitchenConfig.ID, phaseInfo.CurrentPhase)
			if err != nil {
				return err
			}

			currentLevel.Reward = phaseRewards.Reward
			userKitchenProgress.StationLevels[slug] = currentLevel

			userKitchenProgress.NextLevelStats[slug] = entities.UserStationLevel{
				Level:  nextLevel,
				Cost:   nextCost,
				Profit: nextProfit,
			}
		}
	}

	// Attach user progression to config for DTO mapping
	config.UserProgress = userKitchenProgress

	return nil
}

func (g *gameUseCase) mapToUserGameStage(stages []entities.GameStage, lastProgress *entities.UserGameStageProgression) ([]entities.UserGameStage, *entities.UserNextGameStageInfo) {
	isFoundCurrent := false

	var res = make([]entities.UserGameStage, len(stages))
	var nextStage *entities.UserNextGameStageInfo

	for i, stage := range stages {
		currentStage := entities.UserGameStage{
			Slug:        stage.Slug,
			Name:        stage.Name,
			Description: stage.Description,
			Sequence:    stage.Sequence,
		}

		// User hasn't started any stage yet
		if lastProgress == nil {
			// First stage should be current/available
			if i == 0 {
				currentStage.Status = entities.GSStatusAvailable
				isFoundCurrent = true

				nextStage = g.setNextStageIfExists(stages, i+1, entities.GSStatusLocked)
			} else {
				// All other stages are locked
				currentStage.Status = entities.GSStatusLocked
			}
		} else {
			// User has some progression

			// This is the stage user is currently on or has completed
			if lastProgress.StageID == stage.ID {
				if lastProgress.IsComplete {
					currentStage.Status = entities.GSStatusComplete

					nextStage = g.setNextStageIfExists(stages, i+1, entities.GSStatusAvailable)
				} else {
					currentStage.Status = entities.GSStatusAvailable

					nextStage = g.setNextStageIfExists(stages, i+1, entities.GSStatusLocked)
				}
			} else if currentStage.Sequence < g.getSequenceByID(stages, lastProgress.StageID) {
				// This stage comes before the user's last progress, so it's completed
				currentStage.Status = entities.GSStatusComplete
			} else {
				// This stage comes after the user's last progress

				// If last progress is complete, and we haven't set a current stage yet,
				// this is the next available stage
				if lastProgress.IsComplete && !isFoundCurrent {
					currentStage.Status = entities.GSStatusAvailable
					isFoundCurrent = true

					nextStage = g.setNextStageIfExists(stages, i+1, entities.GSStatusLocked)
				} else {
					// Stage is still locked (either already found current or last progress incomplete)
					currentStage.Status = entities.GSStatusLocked
				}
			}
		}

		res[i] = currentStage
	}

	return res, nextStage
}

func (g *gameUseCase) CompleteGameStage(ctx context.Context, slug string) (err error) {
	userID, err := helper.GetUserIDFromContext(ctx)
	if err != nil {
		return err
	}

	stage, err := g.gameStageRepo.GetGameStageBySlugDB(ctx, slug)
	if err != nil {
		return err
	}

	if stage == nil {
		return apperror.ErrRecordNotFound
	}

	userStageProgress, err := g.userProgressionRepo.GetGameStageProgressionDB(ctx, userID, stage.ID)
	if err != nil {
		return err
	}

	if err = g.validateLastProgression(userStageProgress, stage); err != nil {
		return err
	}

	err = g.userProgressionRepo.MarkStageAsCompleteDB(ctx, userID, stage.ID)
	if err != nil {
		return err
	}

	gameStages, err := g.gameStageRepo.GetActiveGameStagesDB(ctx)
	if err != nil {
		return err
	}

	for i, s := range gameStages {
		if s.ID == stage.ID && i+1 < len(gameStages) {
			nextStageID := gameStages[i+1].ID
			_, err = g.userProgressionRepo.CreateGameStageProgressionDB(ctx, userID, nextStageID)
			if err != nil && !errors.Is(err, apperror.ErrConflict) {
				return err
			}
			break
		}
	}

	return nil
}

func (g *gameUseCase) UnlockKitchenStation(ctx context.Context, slug string) (*entities.UnlockKitchenStation, error) {
	userID, err := helper.GetUserIDFromContext(ctx)
	if err != nil {
		return nil, err
	}

	// Gather data for unlock
	unlockCtx, err := g.gatherUnlockData(ctx, userID, slug)
	if err != nil {
		return nil, err
	}

	// Validate unlock requirements
	if err := g.validateUnlockRequirements(unlockCtx); err != nil {
		return nil, err
	}

	// Calculate unlock cost
	result := g.calculateUnlockCost(unlockCtx)

	// Check sufficient funds
	if unlockCtx.userBalance.Coin < result.unlockCost {
		return nil, apperror.ErrInsufficientCoins
	}

	// Execute unlock transaction
	if err := g.executeUnlockTransaction(ctx, unlockCtx, result); err != nil {
		return nil, err
	}

	// Build and return response
	return g.buildUnlockResponse(unlockCtx, result), nil
}

func (g *gameUseCase) UpgradeKitchenStation(ctx context.Context, slug string) (*entities.UpgradeKitchenStation, error) {
	userID, err := helper.GetUserIDFromContext(ctx)
	if err != nil {
		return nil, err
	}

	// Gather all required data
	upgradeCtx, err := g.gatherUpgradeData(ctx, userID, slug)
	if err != nil {
		return nil, err
	}

	// Validate upgrade requirements
	if err := g.validateUpgradeRequirements(upgradeCtx); err != nil {
		return nil, err
	}

	// process override level if any
	overrideCurrentLevel, overrideNextLevel, err := g.proceedOverrideLevel(ctx, upgradeCtx)
	if err != nil {
		return nil, err
	}

	// Calculate upgrade metrics
	result := g.calculateUpgradeMetrics(upgradeCtx, overrideCurrentLevel, overrideNextLevel)

	// Check sufficient funds
	if upgradeCtx.userBalance.Coin < result.upgradeCost {
		return nil, apperror.ErrInsufficientCoins
	}

	// Execute upgrade transaction
	if err := g.executeUpgradeTransaction(ctx, upgradeCtx, result); err != nil {
		return nil, err
	}

	// Build and return response
	return g.buildUpgradeResponse(upgradeCtx, result), nil
}

func (g *gameUseCase) validateLastProgression(lastProgression *entities.UserGameStageProgression, stage *entities.GameStage) error {
	if lastProgression == nil || stage == nil {
		return apperror.ErrStageLocked
	}

	if lastProgression.StageID != stage.ID {
		return apperror.ErrStageLocked
	}

	if lastProgression.IsComplete {
		return apperror.ErrStageAlreadyCompleted
	}

	return nil
}

func (g *gameUseCase) getSequenceByID(stages []entities.GameStage, id int64) int64 {
	for _, s := range stages {
		if s.ID == id {
			return s.Sequence
		}
	}
	return 0
}

func (g *gameUseCase) setNextStageIfExists(stages []entities.GameStage, index int, status entities.GameStageStatus) *entities.UserNextGameStageInfo {
	if index >= len(stages) {
		return nil
	}

	return &entities.UserNextGameStageInfo{
		Slug:   stages[index].Slug,
		Name:   stages[index].Name,
		Status: status,
	}
}

func (g *gameUseCase) GetStageUpgrades(ctx context.Context) (res []entities.UserStageUpgrade, err error) {
	_, err = helper.GetUserIDFromContext(ctx)
	if err != nil {
		return nil, err
	}

	latestProgression, err := g.userProgressionUseCase.LatestStageProgression(ctx)
	if err != nil {
		return nil, err
	}

	stageID := latestProgression.StageID
	stageUpgrades, err := g.userProgressionUseCase.GetActiveStageUpgrade(ctx, stageID)
	if err != nil {
		return nil, err
	}

	return stageUpgrades, nil
}

func (g *gameUseCase) PurchaseStageUpgrade(ctx context.Context, slug string) (res *entities.Upgrade, err error) {
	userID, err := helper.GetUserIDFromContext(ctx)
	if err != nil {
		return nil, err
	}

	upgradeData, err := g.gatherStageUpgradeData(ctx, userID, slug)
	if err != nil {
		return nil, err
	}

	if err = g.validateStageUpgrade(upgradeData); err != nil {
		return nil, err
	}

	err = g.executeStageUpgradeTransaction(ctx, upgradeData)
	if err != nil {
		return nil, err
	}

	return &upgradeData.stageUpgrade.Upgrade, nil
}
