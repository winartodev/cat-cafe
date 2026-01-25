package usecase

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"github.com/winartodev/cat-cafe/internal/dto"
	"github.com/winartodev/cat-cafe/internal/entities"
	"github.com/winartodev/cat-cafe/internal/repositories"
	"github.com/winartodev/cat-cafe/pkg/apperror"
	"github.com/winartodev/cat-cafe/pkg/helper"
	"math"
)

type GameUseCase interface {
	UpdateUserBalance(ctx context.Context, coinEarned int64, userID int64) (res *dto.UserBalanceResponse, err error)
	GetUserGameData(ctx context.Context, userID int64) (res *entities.Game, err error)
	GetGameStages(ctx context.Context, userID int64) (res []entities.UserGameStage, nextStage *entities.UserNextGameStageInfo, err error)
	StartGameStage(ctx context.Context, userID int64, slug string) (stage *entities.GameStage, config *entities.GameStageConfig, nextStage *entities.UserNextGameStageInfo, err error)
	CompleteGameStage(ctx context.Context, userID int64, slug string) (err error)
	UpgradeKitchenStation(ctx context.Context, userID int64, slug string) (res *entities.UpgradeKitchenStation, err error)
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
	}
}

func (g *gameUseCase) UpdateUserBalance(ctx context.Context, coinEarned int64, userID int64) (res *dto.UserBalanceResponse, err error) {
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

func (g *gameUseCase) GetUserGameData(ctx context.Context, userID int64) (res *entities.Game, err error) {
	isDailyRewardAvailable, err := g.userUseCase.IsDailyRewardAvailable(ctx, userID)
	if err != nil {
		return nil, err
	}

	return &entities.Game{
		DailyRewardAvailable: isDailyRewardAvailable,
	}, nil
}

func (g *gameUseCase) GetGameStages(ctx context.Context, userID int64) (res []entities.UserGameStage, nextStage *entities.UserNextGameStageInfo, err error) {
	gameStages, err := g.gameStageRepo.GetActiveGameStagesDB(ctx)
	if err != nil || len(gameStages) == 0 {
		return nil, nil, err
	}

	latestProgress, err := g.userProgressionRepo.GetLatestGameStageProgressionDB(ctx, userID)
	if err != nil {
		return nil, nil, err
	}

	// Create new progression if player doesn't have any
	if latestProgress == nil {
		firstStage := gameStages[0]
		_, err = g.userProgressionRepo.CreateGameStageProgressionDB(ctx, userID, firstStage.ID)
		if err != nil {
			return nil, nil, err
		}

		return g.GetGameStages(ctx, userID)
	}

	stages, nextStage := g.mapToUserGameStage(gameStages, latestProgress)

	return stages, nextStage, nil
}

func (g *gameUseCase) StartGameStage(ctx context.Context, userID int64, slug string) (stage *entities.GameStage, config *entities.GameStageConfig, nextStage *entities.UserNextGameStageInfo, err error) {
	stage, err = g.gameStageRepo.GetGameStageBySlugDB(ctx, slug)
	if err != nil {
		return nil, nil, nil, err
	}

	// Get user's available stages
	userStages, nextStage, err := g.GetGameStages(ctx, userID)
	if err != nil {
		return nil, nil, nil, err
	}

	// Check if user can access this stage
	var canAccess bool
	for _, userStage := range userStages {
		if userStage.Slug == slug {
			// Only allow if status is Current or Complete (if you allow replay)
			if userStage.Status == entities.GSStatusCurrent {
				canAccess = true
			} else if userStage.Status == entities.GSStatusComplete {
				// Decide if you want to allow replaying completed stages
				return nil, nil, nil, apperror.ErrStageAlreadyCompleted
			}
			break
		}
	}

	if !canAccess {
		return nil, nil, nil, apperror.ErrStageNotUnlocked
	}

	config, err = g.gameStageRepo.GetGameConfigByIDDB(ctx, stage.ID)
	if err != nil {
		return nil, nil, nil, err
	}

	err = g.userProgressionUseCase.InitializeUserProgression(ctx, userID, stage.ID, config)
	if err != nil {
		return nil, nil, nil, err
	}

	return stage, config, nextStage, nil
}

func (g *gameUseCase) mapToUserGameStage(stages []entities.GameStage, lastProgress *entities.UserGameStageProgression) ([]entities.UserGameStage, *entities.UserNextGameStageInfo) {
	isFoundCurrent := false

	var res = make([]entities.UserGameStage, len(stages))
	var nextStage *entities.UserNextGameStageInfo

	for i, stage := range stages {
		currentStage := entities.UserGameStage{
			Slug:     stage.Slug,
			Name:     stage.Name,
			Sequence: stage.Sequence,
		}

		// User hasn't started any stage yet
		if lastProgress == nil {
			// First stage should be current/available
			if i == 0 {
				currentStage.Status = entities.GSStatusCurrent
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

					nextStage = g.setNextStageIfExists(stages, i+1, entities.GSStatusCurrent)
				} else {
					currentStage.Status = entities.GSStatusCurrent

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
					currentStage.Status = entities.GSStatusCurrent
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

func (g *gameUseCase) CompleteGameStage(ctx context.Context, userID int64, slug string) (err error) {
	stage, err := g.gameStageRepo.GetGameStageBySlugDB(ctx, slug)
	if err != nil {
		return err
	}

	if stage == nil {
		return apperror.ErrRecordNotFound
	}

	// TODO: FIX THIS EITHER USE LATEST GAME STAGE PROGRESSION OR BY STAGE GAME
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

func (g *gameUseCase) UpgradeKitchenStation(ctx context.Context, userID int64, slug string) (res *entities.UpgradeKitchenStation, err error) {
	// Get latest progression
	latestProgression, err := g.userProgressionRepo.GetLatestGameStageProgressionDB(ctx, userID)
	if err != nil {
		return nil, err
	}
	if latestProgression == nil {
		return nil, apperror.ErrRecordNotFound
	}

	// Get food item by slug
	foodItem, err := g.foodItemRepo.GetFoodBySlugDB(ctx, slug)
	if err != nil {
		return nil, err
	}
	if foodItem == nil {
		return nil, apperror.ErrRecordNotFound
	}

	stageID := latestProgression.StageID

	// Get kitchen station
	kitchenStation, err := g.kitchenStationRepo.GetKitchenStationByFoodIDDB(ctx, stageID, foodItem.ID)
	if err != nil {
		return nil, err
	}
	if kitchenStation == nil {
		return nil, apperror.ErrRecordNotFound
	}

	// Get kitchen config
	kitchenConfig, err := g.kitchenConfigRepo.GetKitchenConfigByStageIDDB(ctx, stageID)
	if err != nil {
		return nil, err
	}

	// Get user kitchen progress
	userProgress, err := g.userProgressionRepo.GetUserKitchenProgressDB(ctx, userID, stageID)
	if err != nil {
		return nil, err
	}

	// Get user phase progress
	phaseProgress, err := g.userProgressionRepo.GetUserKitchenPhaseProgressionDB(ctx, userID, kitchenConfig.ID)
	if err != nil {
		return nil, err
	}

	// Get user balance
	userBalance, err := g.userRepo.GetUserBalanceByIDDB(ctx, userID)
	if err != nil {
		return nil, err
	}
	if userBalance == nil {
		return nil, apperror.ErrRecordNotFound
	}

	stationLevels := userProgress.StationLevels
	currentLevel, exists := stationLevels[slug]
	if !exists {
		currentLevel = 0
	}

	if currentLevel >= kitchenConfig.MaxLevel {
		return nil, apperror.ErrMaxLevelReached
	}

	nextLevel := currentLevel + 1

	oldPhaseInfo := g.calculatePhaseInfo(currentLevel, kitchenConfig)
	newPhaseInfo := g.calculatePhaseInfo(nextLevel, kitchenConfig)

	upgradeCost := g.calculateUpgradeCost(
		foodItem.StartingPrice,
		nextLevel,
		kitchenConfig,
		newPhaseInfo.CurrentPhase,
	)

	currentProfit := g.calculateProfit(
		foodItem.StartingPrice,
		nextLevel,
		kitchenConfig,
		newPhaseInfo.CurrentPhase,
		0,
	)

	preparationTIme := g.calculateCurrentProcessTime(
		foodItem.StartingPreparation,
		1, 1,
	)

	if userBalance.Coin < upgradeCost {
		return nil, apperror.ErrInsufficientCoins
	}

	var currentCoin = userBalance.Coin - upgradeCost

	var phaseTransitioned bool
	var phaseRewards []entities.PhaseRewardInfo
	var allPhaseRewardsCollected []entities.PhaseRewardInfo

	var newTableCount int64

	err = g.userProgressionRepo.WithUserProgressionTx(ctx, func(tx *sql.Tx) error {
		userRepo := g.userRepo.WithTx(tx)
		userProgressionRepo := g.userProgressionRepo.WithTx(tx)

		// Deduct coins
		err := userRepo.UpdateUserBalanceWithTx(ctx, userID, entities.BalanceTypeCoin, -upgradeCost)
		if err != nil {
			return err
		}

		stationLevels[slug] = nextLevel
		userProgress.StationLevels = stationLevels
		err = userProgressionRepo.UpdateUserKitchenProgressDB(ctx, userID, stageID, userProgress)

		if newPhaseInfo.CurrentPhase > oldPhaseInfo.CurrentPhase {
			phaseTransitioned = true

			err = g.updateKitchenPhaseProgression(
				ctx,
				tx,
				userID,
				kitchenConfig.ID,
				phaseProgress,
				newPhaseInfo.CurrentPhase,
			)
			if err != nil {
				return err
			}

			// Get new table count
			if int(newPhaseInfo.CurrentPhase) <= len(kitchenConfig.TableCountPerPhases) {
				newTableCount = kitchenConfig.TableCountPerPhases[newPhaseInfo.CurrentPhase-1]
			}

			phaseRewards, err = g.collectPhaseCompletionRewards(
				ctx,
				tx,
				userID,
				stageID,
				kitchenConfig,
				oldPhaseInfo.CurrentPhase,
				newPhaseInfo.CurrentPhase,
			)
			if err != nil {
				// Log error but don't fail upgrade
				fmt.Printf("Error collecting phase rewards: %v\n", err)
				phaseRewards = []entities.PhaseRewardInfo{}
				return err
			}
		}

		if nextLevel >= kitchenConfig.MaxLevel {
			// ✅ Collect all remaining phase rewards
			allPhaseRewardsCollected, err = g.collectAllRemainingPhaseRewards(
				ctx,
				tx,
				userID,
				stageID,
				kitchenConfig,
				newPhaseInfo.CurrentPhase,
			)
			if err != nil {
				fmt.Printf("Error collecting all phase rewards: %v\n", err)
				allPhaseRewardsCollected = []entities.PhaseRewardInfo{}
			}
		}

		return nil
	})
	if err != nil {
		return nil, err
	}

	fmt.Println("==================================================")
	fmt.Printf("UPGRADE SUCCESS: %s\n", slug)
	fmt.Printf("Level     : %d -> %d\n", currentLevel, nextLevel)
	fmt.Printf("Cost      : %d Coins\n", upgradeCost)
	fmt.Printf("Balance   : %d -> %d\n", userBalance.Coin, currentCoin)
	fmt.Println("--------------------------------------------------")
	fmt.Printf("Phase     : %d -> %d (Transition: %t)\n", oldPhaseInfo.CurrentPhase, newPhaseInfo.CurrentPhase, phaseTransitioned)
	fmt.Printf("Profit    : %d\n", currentProfit)
	fmt.Printf("Prep Time : %v\n", preparationTIme)

	if phaseTransitioned {
		fmt.Printf("Phase Rewards  : %+v\n", phaseRewards)
		fmt.Printf("New Table Count: %d\n", newTableCount)
	}

	if nextLevel >= kitchenConfig.MaxLevel {
		fmt.Println("MAX LEVEL REACHED!")
		fmt.Printf("Final Rewards: %+v\n", allPhaseRewardsCollected)
	}

	fmt.Println("==================================================")

	response := &entities.UpgradeKitchenStation{
		Success:        true,
		NewLevel:       nextLevel,
		IsMaxLevel:     nextLevel >= kitchenConfig.MaxLevel,
		NewCoinBalance: currentCoin,
		CoinsSpent:     upgradeCost,

		// Current values
		CurrentProfit:   currentProfit,
		CurrentPrepTime: preparationTIme,
		ProfitPerSecond: float64(currentProfit) / preparationTIme,

		// Phase info
		PhaseTransitioned:      phaseTransitioned,
		CurrentPhase:           newPhaseInfo.CurrentPhase,
		CurrentPhaseStartLevel: newPhaseInfo.CurrentPhaseStart,
		CurrentPhaseLastLevel:  newPhaseInfo.CurrentPhaseEnd,
		PhaseProfitMultiplier:  newPhaseInfo.ProfitMultiplier,

		// Table count
		NewTableCount: newTableCount,

		// Rewards
		PhaseRewards:    phaseRewards,
		AllPhaseRewards: allPhaseRewardsCollected,
	}

	return response, nil
}

func (g *gameUseCase) validateLastProgression(lastProgression *entities.UserGameStageProgression, stage *entities.GameStage) error {
	if lastProgression == nil || stage == nil {
		return apperror.ErrStageNotUnlocked
	}

	if lastProgression.StageID != stage.ID {
		return apperror.ErrStageNotUnlocked
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

func (g *gameUseCase) calculatePhaseInfo(level int64, config *entities.StageKitchenConfig) entities.PhaseInfo {
	info := entities.PhaseInfo{
		CurrentPhase:      1,
		CurrentPhaseStart: 1,
		CurrentPhaseEnd:   config.MaxLevel,
	}

	if len(config.TransitionPhaseLevels) == 0 {
		return info
	}

	for i := 0; i < len(config.TransitionPhaseLevels); i++ {
		if config.TransitionPhaseLevels[i] < level {
			info.CurrentPhase = int64(i + 1)
			info.CurrentPhaseStart = config.TransitionPhaseLevels[i]

			// Calculate end phase
			if i+1 < len(config.TransitionPhaseLevels) {
				info.CurrentPhaseEnd = config.TransitionPhaseLevels[i+1] - 1
			} else {
				info.CurrentPhaseEnd = config.MaxLevel
			}

			// Get profit multiplier for this phase
			if i < len(config.PhaseProfitMultipliers) {
				info.ProfitMultiplier = config.PhaseProfitMultipliers[i]
			}

			// Get upgrade cost multiplier for this phase
			if i < len(config.PhaseUpgradeCostMultipliers) {
				info.UpgradeCostMultiplier = config.PhaseUpgradeCostMultipliers[i]
			}
		} else {
			break
		}
	}

	return info
}

// calculateUpgradeCost: basePrice * (upgradeCostMultiply/100)^(level-1) * phaseMultiplier
func (g *gameUseCase) calculateUpgradeCost(
	basePrice int64,
	level int64,
	config *entities.StageKitchenConfig,
	currentPhase int64,
) int64 {
	phaseMultiplier := 1.0
	if int(currentPhase) <= len(config.PhaseUpgradeCostMultipliers) {
		phaseMultiplier = config.PhaseUpgradeCostMultipliers[currentPhase-1]
	}

	cost := float64(basePrice)
	multiplier := float64(config.UpgradeCostMultiply) / 100.0

	cost = cost * math.Pow(multiplier, float64(level-1)) * phaseMultiplier

	return int64(math.Ceil(cost))
}

// calculateProfit = basePrice × (costMultiplier/100)^(level-1) × phaseCostMultiplier
func (g *gameUseCase) calculateProfit(
	baseProfit int64,
	level int64,
	config *entities.StageKitchenConfig,
	currentPhase int64,
	bonusMultiply float64,
) int64 {
	phaseMultiplier := 1.0
	if int(currentPhase) <= len(config.PhaseProfitMultipliers) {
		phaseMultiplier = config.PhaseProfitMultipliers[currentPhase-1]
	}

	// Calculate base profit for level
	profit := float64(baseProfit)
	multiplier := float64(config.UpgradeProfitMultiply) / 100.0

	profit = profit * math.Pow(multiplier, float64(level-1)) * phaseMultiplier

	// Apply bonus (kitchen profit bonus from upgrades/buffs)
	if bonusMultiply > 0 {
		profit = profit * (1 + bonusMultiply)
	}

	return int64(math.Ceil(profit))
}

// calculateCurrentProcessTime = baseTime - (baseTime × permanentReduction) + (baseTime × temporaryModifier)
func (g *gameUseCase) calculateCurrentProcessTime(
	baseProcessTime float64,
	reduceTimeMultiply float64,
	bonusReduceTime float64,
) float64 {
	processTime := baseProcessTime - (baseProcessTime * reduceTimeMultiply) + (baseProcessTime * bonusReduceTime)

	if processTime < 0.1 {
		processTime = 0.1 // Minimum process time
	}

	return processTime
}

func (g *gameUseCase) updateKitchenPhaseProgression(
	ctx context.Context,
	tx *sql.Tx,
	userID int64,
	kitchenConfigID int64,
	currentProgress *entities.UserKitchenPhaseProgression,
	newPhase int64,
) error {
	userProgressionRepo := g.userProgressionRepo.WithTx(tx)

	// Don't update if same phase
	if int64(currentProgress.CurrentPhase) >= int64(newPhase) {
		return nil
	}

	// Add completed phases
	completedPhases := currentProgress.CompletedPhases

	// Mark all phases between current and new as completed
	for phase := currentProgress.CurrentPhase; phase < newPhase; phase++ {
		alreadyCompleted := false
		for _, cp := range completedPhases {
			if cp == phase {
				alreadyCompleted = true
				break
			}
		}

		if !alreadyCompleted {
			completedPhases = append(completedPhases, phase)
		}
	}

	// Update phase progression
	return userProgressionRepo.UpdateUserKitchenPhaseProgressionDB(ctx, userID, kitchenConfigID, &entities.UserKitchenPhaseProgression{
		CurrentPhase:    newPhase,
		CompletedPhases: completedPhases,
	})
}

func (g *gameUseCase) collectPhaseCompletionRewards(
	ctx context.Context,
	tx *sql.Tx,
	userID int64,
	stageID int64,
	kitchenConfig *entities.StageKitchenConfig,
	fromPhase int64,
	toPhase int64,
) ([]entities.PhaseRewardInfo, error) {
	kitchenConfigRepo := g.kitchenConfigRepo.WithTx(tx)
	userProgression := g.userProgressionRepo.WithTx(tx)

	var rewardInfos []entities.PhaseRewardInfo
	phaseRewards, err := kitchenConfigRepo.GetKitchenCompletionRewardsDB(ctx, stageID)
	if err != nil {
		return nil, err
	}

	// Collect rewards for completed phases
	// Example: fromPhase=1, toPhase=2 → collect reward for phase 1
	for phase := fromPhase; phase < toPhase; phase++ {
		for _, phaseReward := range phaseRewards {
			if phaseReward.PhaseNumber == phase {
				fmt.Println("Get Rewards")
				// Check if already claimed
				claimed, err := userProgression.IsPhaseRewardAlreadyClaimedDB(
					ctx,
					userID,
					kitchenConfig.ID,
					phase,
					phaseReward.RewardID,
				)
				if err != nil {
					return nil, err
				}

				if claimed {
					continue
				}

				// Grant the reward
				rewardInfo, err := g.grantPhaseReward(
					ctx,
					tx,
					userID,
					kitchenConfig.ID,
					&phaseReward,
				)
				if err != nil {
					// Log error but continue with other rewards
					return nil, err
				}

				rewardInfos = append(rewardInfos, rewardInfo)

			}
		}
	}

	return rewardInfos, nil
}

func (g *gameUseCase) grantPhaseReward(
	ctx context.Context,
	tx *sql.Tx,
	userID int64,
	kitchenConfigID int64,
	phaseReward *entities.KitchenPhaseCompletionRewards,
) (entities.PhaseRewardInfo, error) {
	fmt.Println("grantPhaseReward")
	rewardRepo := g.rewardRepo.WithTx(tx)
	userRepo := g.userRepo.WithTx(tx)
	userProgression := g.userProgressionRepo.WithTx(tx)

	fmt.Println("phaseReward: ", phaseReward.RewardID)

	// Get reward details
	reward, err := rewardRepo.GetRewardByIDDB(ctx, phaseReward.RewardID)
	if err != nil {
		return entities.PhaseRewardInfo{}, err
	}

	if reward == nil {
		return entities.PhaseRewardInfo{}, apperror.ErrRecordNotFound
	}

	rewardType := reward.RewardType
	rewardTypeEnum, err := entities.ToRewardType(rewardType.Slug)
	if err != nil {
		return entities.PhaseRewardInfo{}, err
	}

	if rewardTypeEnum.RequiresBalanceUpdate() {
		balanceType := rewardTypeEnum.ToUserBalance()
		err = userRepo.UpdateUserBalanceWithTx(ctx, userID, balanceType, reward.Amount)
		if err != nil {
			return entities.PhaseRewardInfo{}, err
		}
	} else if rewardTypeEnum.IsSentExternally() {
		// TODO: Call External API to give GoPay Coin to player
	}

	// Record that reward was claimed
	err = userProgression.CreateUserKitchenClaimRewardDB(
		ctx,
		entities.UserKitchenPhaseRewardClaim{
			UserID:          userID,
			KitchenConfigID: kitchenConfigID,
			RewardID:        phaseReward.RewardID,
			CurrentPhase:    phaseReward.PhaseNumber,
		},
	)
	if err != nil {
		return entities.PhaseRewardInfo{}, err
	}

	return entities.PhaseRewardInfo{
		PhaseNumber: int(phaseReward.PhaseNumber),
		RewardType:  rewardTypeEnum.String(),
		RewardSlug:  reward.Slug,
		RewardName:  reward.Name,
		Amount:      reward.Amount,
	}, nil
}

func (g *gameUseCase) collectAllRemainingPhaseRewards(
	ctx context.Context,
	tx *sql.Tx,
	userID int64,
	stageID int64,
	kitchenConfig *entities.StageKitchenConfig,
	currentPhase int64,
) ([]entities.PhaseRewardInfo, error) {
	fmt.Println("collectAllRemainingPhaseRewards")
	phaseReward := g.kitchenConfigRepo.WithTx(tx)
	userProgression := g.userProgressionRepo.WithTx(tx)

	var rewardInfos []entities.PhaseRewardInfo

	// Get all phase rewards
	phaseRewards, err := phaseReward.GetKitchenCompletionRewardsDB(ctx, stageID)
	if err != nil {
		return nil, err
	}

	// Collect rewards from current phase onwards
	for _, phaseReward := range phaseRewards {
		if phaseReward.PhaseNumber >= int64(currentPhase) {
			// Check if already claimed
			claimed, err := userProgression.IsPhaseRewardAlreadyClaimedDB(
				ctx,
				userID,
				kitchenConfig.ID,
				phaseReward.PhaseNumber,
				phaseReward.RewardID,
			)
			if err != nil {
				continue
			}

			if claimed {
				continue
			}

			// Grant the reward
			rewardInfo, err := g.grantPhaseReward(
				ctx,
				tx,
				userID,
				kitchenConfig.ID,
				&phaseReward,
			)
			if err != nil {
				fmt.Printf("Error granting reward: %v\n", err)
				continue
			}

			rewardInfos = append(rewardInfos, rewardInfo)
		}
	}

	return rewardInfos, nil
}
