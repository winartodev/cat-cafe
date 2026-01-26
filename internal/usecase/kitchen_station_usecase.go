package usecase

import (
	"context"
	"database/sql"
	"fmt"
	"github.com/winartodev/cat-cafe/internal/entities"
	"github.com/winartodev/cat-cafe/internal/repositories"
	"github.com/winartodev/cat-cafe/pkg/apperror"
	"math"
)

type unlockContext struct {
	userID        int64
	stageID       int64
	slug          string
	foodItem      *entities.FoodItem
	kitchenConfig *entities.StageKitchenConfig
	userProgress  *entities.UserKitchenStageProgression
	userBalance   *entities.UserBalance
}

type upgradeContext struct {
	userID         int64
	stageID        int64
	slug           string
	foodItem       *entities.FoodItem
	kitchenStation *entities.KitchenStation
	kitchenConfig  *entities.StageKitchenConfig
	userProgress   *entities.UserKitchenStageProgression
	phaseProgress  *entities.UserKitchenPhaseProgression
	userBalance    *entities.UserBalance
	currentLevel   int64
	nextLevel      int64
}

type unlockResult struct {
	unlockCost     int64
	newCoinBalance int64
	unlockedSlug   string
}

// upgradeResult holds the results of the upgrade operation
type upgradeResult struct {
	upgradeCost              int64
	newCoinBalance           int64
	currentProfit            int64
	preparationTime          float64
	oldPhaseInfo             phaseInfo
	newPhaseInfo             phaseInfo
	phaseTransitioned        bool
	newTableCount            int64
	phaseRewards             []entities.PhaseRewardInfo
	allPhaseRewardsCollected []entities.PhaseRewardInfo
}

type phaseInfo struct {
	CurrentPhase          int64
	CurrentPhaseStart     int64
	CurrentPhaseEnd       int64
	ProfitMultiplier      float64
	UpgradeCostMultiplier float64
	TableCount            int64
}

func (g *gameUseCase) UnlockKitchenStation(ctx context.Context, userID int64, slug string) (*entities.UnlockKitchenStation, error) {
	// Step 1: Gather data for unlock
	unlockCtx, err := g.gatherUnlockData(ctx, userID, slug)
	if err != nil {
		return nil, err
	}

	// Step 2: Validate unlock requirements
	if err := g.validateUnlockRequirements(unlockCtx); err != nil {
		return nil, err
	}

	// Step 3: Calculate unlock cost
	result := g.calculateUnlockCost(unlockCtx)

	// Step 4: Check sufficient funds
	if unlockCtx.userBalance.Coin < result.unlockCost {
		return nil, apperror.ErrInsufficientCoins
	}

	// Step 5: Execute unlock transaction
	if err := g.executeUnlockTransaction(ctx, unlockCtx, result); err != nil {
		return nil, err
	}

	// Step 6: Log unlock details
	g.logUnlockDetails(unlockCtx, result)

	// Step 7: Build and return response
	return g.buildUnlockResponse(unlockCtx, result), nil
}

func (g *gameUseCase) UpgradeKitchenStation(ctx context.Context, userID int64, slug string) (*entities.UpgradeKitchenStation, error) {
	// Gather all required data
	upgradeCtx, err := g.gatherUpgradeData(ctx, userID, slug)
	if err != nil {
		return nil, err
	}

	// Validate upgrade requirements
	if err := g.validateUpgradeRequirements(upgradeCtx); err != nil {
		return nil, err
	}

	// Calculate upgrade metrics
	result := g.calculateUpgradeMetrics(upgradeCtx)

	// Check sufficient funds
	if upgradeCtx.userBalance.Coin < result.upgradeCost {
		return nil, apperror.ErrInsufficientCoins
	}

	// Execute upgrade transaction
	if err := g.executeUpgradeTransaction(ctx, upgradeCtx, result); err != nil {
		return nil, err
	}

	// Log upgrade details
	g.logUpgradeDetails(upgradeCtx, result)

	// Build and return response
	return g.buildUpgradeResponse(upgradeCtx, result), nil
}

func (g *gameUseCase) gatherUnlockData(ctx context.Context, userID int64, slug string) (*unlockContext, error) {
	uctx := &unlockContext{
		userID: userID,
		slug:   slug,
	}

	// Get latest progression
	latestProgression, err := g.userProgressionRepo.GetLatestGameStageProgressionDB(ctx, userID)
	if err != nil {
		return nil, err
	}
	if latestProgression == nil {
		return nil, apperror.ErrRecordNotFound
	}
	uctx.stageID = latestProgression.StageID

	// Get food item
	uctx.foodItem, err = g.foodItemRepo.GetFoodBySlugDB(ctx, slug)
	if err != nil {
		return nil, err
	}
	if uctx.foodItem == nil {
		return nil, apperror.ErrRecordNotFound
	}

	// Get kitchen config
	uctx.kitchenConfig, err = g.kitchenConfigRepo.GetKitchenConfigByStageIDDB(ctx, uctx.stageID)
	if err != nil {
		return nil, err
	}

	// Get user kitchen progress
	uctx.userProgress, err = g.userProgressionRepo.GetUserKitchenProgressDB(ctx, userID, uctx.stageID)
	if err != nil {
		return nil, err
	}

	// Get user balance
	uctx.userBalance, err = g.userRepo.GetUserBalanceByIDDB(ctx, userID)
	if err != nil {
		return nil, err
	}
	if uctx.userBalance == nil {
		return nil, apperror.ErrRecordNotFound
	}

	return uctx, nil
}

func (g *gameUseCase) gatherUpgradeData(ctx context.Context, userID int64, slug string) (*upgradeContext, error) {
	upgradeContext := &upgradeContext{
		userID: userID,
		slug:   slug,
	}

	// Get latest progression
	latestProgression, err := g.userProgressionRepo.GetLatestGameStageProgressionDB(ctx, userID)
	if err != nil {
		return nil, err
	}
	if latestProgression == nil {
		return nil, apperror.ErrRecordNotFound
	}
	upgradeContext.stageID = latestProgression.StageID

	// Get food item
	upgradeContext.foodItem, err = g.foodItemRepo.GetFoodBySlugDB(ctx, slug)
	if err != nil {
		return nil, err
	}
	if upgradeContext.foodItem == nil {
		return nil, apperror.ErrRecordNotFound
	}

	// Get kitchen station
	upgradeContext.kitchenStation, err = g.kitchenStationRepo.GetKitchenStationByFoodIDDB(ctx, upgradeContext.stageID, upgradeContext.foodItem.ID)
	if err != nil {
		return nil, err
	}
	if upgradeContext.kitchenStation == nil {
		return nil, apperror.ErrRecordNotFound
	}

	// Get kitchen config
	upgradeContext.kitchenConfig, err = g.kitchenConfigRepo.GetKitchenConfigByStageIDDB(ctx, upgradeContext.stageID)
	if err != nil {
		return nil, err
	}

	// Get user kitchen progress
	upgradeContext.userProgress, err = g.userProgressionRepo.GetUserKitchenProgressDB(ctx, userID, upgradeContext.stageID)
	if err != nil {
		return nil, err
	}

	// Get user phase progress
	upgradeContext.phaseProgress, err = g.userProgressionRepo.GetUserKitchenPhaseProgressionDB(ctx, userID, upgradeContext.kitchenConfig.ID)
	if err != nil {
		return nil, err
	}

	// Get user balance
	upgradeContext.userBalance, err = g.userRepo.GetUserBalanceByIDDB(ctx, userID)
	if err != nil {
		return nil, err
	}
	if upgradeContext.userBalance == nil {
		return nil, apperror.ErrRecordNotFound
	}

	return upgradeContext, nil
}

func (g *gameUseCase) validateUnlockRequirements(unlockContext *unlockContext) error {
	// Check if station is already unlocked
	if g.isStationUnlocked(unlockContext.userProgress.UnlockedStations, unlockContext.slug) {
		return apperror.ErrStationAlreadyUnlocked
	}

	return nil
}

func (g *gameUseCase) validateUpgradeRequirements(upgradeContext *upgradeContext) error {
	// Check if station is unlocked
	if !g.isStationUnlocked(upgradeContext.userProgress.UnlockedStations, upgradeContext.slug) {
		return apperror.ErrStationNotUnlocked
	}

	// Get current level
	currentLevel, exists := upgradeContext.userProgress.StationLevels[upgradeContext.slug]
	if !exists {
		currentLevel = 0
	}
	upgradeContext.currentLevel = currentLevel
	upgradeContext.nextLevel = currentLevel + 1

	// Check max level
	if upgradeContext.currentLevel >= upgradeContext.kitchenConfig.MaxLevel {
		return apperror.ErrMaxLevelReached
	}

	return nil
}

func (g *gameUseCase) isStationUnlocked(unlockedStations []string, slug string) bool {
	if len(unlockedStations) == 0 {
		return false
	}

	for _, station := range unlockedStations {
		if station == slug {
			return true
		}
	}
	return false
}

func (g *gameUseCase) calculateUnlockCost(unlockContext *unlockContext) *unlockResult {
	return &unlockResult{
		unlockCost:     unlockContext.foodItem.StartingPrice,
		newCoinBalance: unlockContext.userBalance.Coin - unlockContext.foodItem.StartingPrice,
		unlockedSlug:   unlockContext.slug,
	}
}

func (g *gameUseCase) calculateUpgradeMetrics(upgradeContext *upgradeContext) *upgradeResult {
	result := &upgradeResult{}

	result.oldPhaseInfo = g.calculatePhaseInfo(upgradeContext.currentLevel, upgradeContext.kitchenConfig)
	result.newPhaseInfo = g.calculatePhaseInfo(upgradeContext.nextLevel, upgradeContext.kitchenConfig)

	result.upgradeCost = g.calculateUpgradeCost(
		upgradeContext.foodItem.StartingPrice,
		upgradeContext.nextLevel,
		upgradeContext.kitchenConfig,
		result.newPhaseInfo.CurrentPhase,
	)

	result.currentProfit = g.calculateProfit(
		upgradeContext.foodItem.StartingPrice,
		upgradeContext.nextLevel,
		upgradeContext.kitchenConfig,
		result.newPhaseInfo.CurrentPhase,
		0,
	)

	result.preparationTime = g.calculateCurrentProcessTime(
		upgradeContext.foodItem.StartingPreparation,
		1, 1,
	)

	result.newCoinBalance = upgradeContext.userBalance.Coin - result.upgradeCost
	result.phaseTransitioned = result.newPhaseInfo.CurrentPhase > result.oldPhaseInfo.CurrentPhase

	return result
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
func (g *gameUseCase) calculateCurrentProcessTime(baseProcessTime float64, reduceTimeMultiply float64, bonusReduceTime float64) float64 {
	processTime := baseProcessTime - (baseProcessTime * reduceTimeMultiply) + (baseProcessTime * bonusReduceTime)

	if processTime < 0.1 {
		processTime = 0.1 // Minimum process time
	}

	return processTime
}

func (g *gameUseCase) executeUnlockTransaction(ctx context.Context, unlockContext *unlockContext, result *unlockResult) error {
	return g.userProgressionRepo.WithUserProgressionTx(ctx, func(tx *sql.Tx) error {
		userRepo := g.userRepo.WithTx(tx)
		userProgressionRepo := g.userProgressionRepo.WithTx(tx)

		// Deduct coins
		if err := userRepo.UpdateUserBalanceWithTx(ctx, unlockContext.userID, entities.BalanceTypeCoin, -result.unlockCost); err != nil {
			return err
		}

		// Unlock the station
		if err := g.unlockStation(ctx, userProgressionRepo, unlockContext); err != nil {
			return err
		}

		return nil
	})
}

func (g *gameUseCase) executeUpgradeTransaction(ctx context.Context, upgradeContext *upgradeContext, result *upgradeResult) error {
	return g.userProgressionRepo.WithUserProgressionTx(ctx, func(tx *sql.Tx) error {
		userRepo := g.userRepo.WithTx(tx)
		userProgressionRepo := g.userProgressionRepo.WithTx(tx)

		// Deduct coins
		if err := userRepo.UpdateUserBalanceWithTx(ctx, upgradeContext.userID, entities.BalanceTypeCoin, -result.upgradeCost); err != nil {
			return err
		}

		// Update station level
		if err := g.updateStationLevel(ctx, userProgressionRepo, upgradeContext, result); err != nil {
			return err
		}

		// Handle phase transition
		if result.phaseTransitioned {
			if err := g.handlePhaseTransition(ctx, tx, upgradeContext, result); err != nil {
				return err
			}
		}

		// Handle max level rewards
		if upgradeContext.nextLevel >= upgradeContext.kitchenConfig.MaxLevel {
			if err := g.handleMaxLevelRewards(ctx, tx, upgradeContext, result); err != nil {
				// Log but don't fail the transaction
				fmt.Printf("Error collecting all phase rewards: %v\n", err)
			}
		}

		return nil
	})
}

func (g *gameUseCase) unlockStation(ctx context.Context, repo repositories.UserProgressionRepository, uctx *unlockContext) error {
	// Add slug to unlocked stations
	unlockedStations := uctx.userProgress.UnlockedStations
	unlockedStations = append(unlockedStations, uctx.slug)
	uctx.userProgress.UnlockedStations = unlockedStations

	// Initialize station level to 0
	if uctx.userProgress.StationLevels == nil {
		uctx.userProgress.StationLevels = make(map[string]int64)
	}
	uctx.userProgress.StationLevels[uctx.slug] = 1

	// Update progress in database
	return repo.UpdateUserKitchenProgressDB(ctx, uctx.userID, uctx.stageID, uctx.userProgress)
}

func (g *gameUseCase) updateStationLevel(ctx context.Context, repo repositories.UserProgressionRepository, upgradeContext *upgradeContext, result *upgradeResult) error {
	stationLevels := upgradeContext.userProgress.StationLevels
	stationLevels[upgradeContext.slug] = upgradeContext.nextLevel
	upgradeContext.userProgress.StationLevels = stationLevels
	return repo.UpdateUserKitchenProgressDB(ctx, upgradeContext.userID, upgradeContext.stageID, upgradeContext.userProgress)
}

func (g *gameUseCase) handlePhaseTransition(ctx context.Context, tx *sql.Tx, upgradeContext *upgradeContext, result *upgradeResult) error {
	// Update phase progression
	if err := g.updateKitchenPhaseProgression(
		ctx,
		tx,
		upgradeContext.userID,
		upgradeContext.kitchenConfig.ID,
		upgradeContext.phaseProgress,
		result.newPhaseInfo.CurrentPhase,
	); err != nil {
		return err
	}

	// Update table count
	if int(result.newPhaseInfo.CurrentPhase) <= len(upgradeContext.kitchenConfig.TableCountPerPhases) {
		result.newTableCount = upgradeContext.kitchenConfig.TableCountPerPhases[result.newPhaseInfo.CurrentPhase-1]
	}

	// Collect phase rewards
	phaseRewards, err := g.collectPhaseCompletionRewards(
		ctx,
		tx,
		upgradeContext.userID,
		upgradeContext.stageID,
		upgradeContext.kitchenConfig,
		result.oldPhaseInfo.CurrentPhase,
		result.newPhaseInfo.CurrentPhase,
	)
	if err != nil {
		fmt.Printf("Error collecting phase rewards: %v\n", err)
		result.phaseRewards = []entities.PhaseRewardInfo{}
		return err
	}
	result.phaseRewards = phaseRewards

	return nil
}

func (g *gameUseCase) calculatePhaseInfo(level int64, config *entities.StageKitchenConfig) phaseInfo {
	info := phaseInfo{
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

func (g *gameUseCase) handleMaxLevelRewards(ctx context.Context, tx *sql.Tx, upgradeContext *upgradeContext, result *upgradeResult) error {
	allRewards, err := g.collectAllRemainingPhaseRewards(
		ctx,
		tx,
		upgradeContext.userID,
		upgradeContext.stageID,
		upgradeContext.kitchenConfig,
		result.newPhaseInfo.CurrentPhase,
	)
	if err != nil {
		return err
	}
	result.allPhaseRewardsCollected = allRewards
	return nil
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

func (g *gameUseCase) logUnlockDetails(unlockContext *unlockContext, result *unlockResult) {
	fmt.Println("==================================================")
	fmt.Printf("UNLOCK SUCCESS: %s\n", unlockContext.slug)
	fmt.Printf("Cost      : %d Coins\n", result.unlockCost)
	fmt.Printf("Balance   : %d -> %d\n", unlockContext.userBalance.Coin, result.newCoinBalance)
	fmt.Printf("Total Unlocked Stations: %d\n", len(unlockContext.userProgress.UnlockedStations)+1)
	fmt.Println("==================================================")
}

func (g *gameUseCase) logUpgradeDetails(upgradeContext *upgradeContext, result *upgradeResult) {
	fmt.Println("==================================================")
	fmt.Printf("UPGRADE SUCCESS: %s\n", upgradeContext.slug)
	fmt.Printf("Level     : %d -> %d\n", upgradeContext.currentLevel, upgradeContext.nextLevel)
	fmt.Printf("Cost      : %d Coins\n", result.upgradeCost)
	fmt.Printf("Balance   : %d -> %d\n", upgradeContext.userBalance.Coin, result.newCoinBalance)
	fmt.Println("--------------------------------------------------")
	fmt.Printf("Phase     : %d -> %d (Transition: %t)\n",
		result.oldPhaseInfo.CurrentPhase,
		result.newPhaseInfo.CurrentPhase,
		result.phaseTransitioned)
	fmt.Printf("Profit    : %d\n", result.currentProfit)
	fmt.Printf("Prep Time : %v\n", result.preparationTime)

	if result.phaseTransitioned {
		fmt.Printf("Phase Rewards  : %+v\n", result.phaseRewards)
		fmt.Printf("New Table Count: %d\n", result.newTableCount)
	}

	if upgradeContext.nextLevel >= upgradeContext.kitchenConfig.MaxLevel {
		fmt.Println("MAX LEVEL REACHED!")
		fmt.Printf("Final Rewards: %+v\n", result.allPhaseRewardsCollected)
	}

	fmt.Println("==================================================")
}

func (g *gameUseCase) buildUnlockResponse(unlockContext *unlockContext, result *unlockResult) *entities.UnlockKitchenStation {
	return &entities.UnlockKitchenStation{
		UnlockedSlug:   result.unlockedSlug,
		NewCoinBalance: result.newCoinBalance,
		CoinsSpent:     result.unlockCost,
		StationName:    unlockContext.foodItem.Name,
		StationLevel:   1,
	}
}

func (g *gameUseCase) buildUpgradeResponse(upgradeContext *upgradeContext, result *upgradeResult) *entities.UpgradeKitchenStation {
	return &entities.UpgradeKitchenStation{
		NewLevel:       upgradeContext.nextLevel,
		IsMaxLevel:     upgradeContext.nextLevel >= upgradeContext.kitchenConfig.MaxLevel,
		NewCoinBalance: result.newCoinBalance,
		CoinsSpent:     result.upgradeCost,

		// Current values
		CurrentProfit:   result.currentProfit,
		CurrentPrepTime: result.preparationTime,
		ProfitPerSecond: float64(result.currentProfit) / result.preparationTime,

		// Phase info
		PhaseTransitioned:      result.phaseTransitioned,
		CurrentPhase:           result.newPhaseInfo.CurrentPhase,
		CurrentPhaseStartLevel: result.newPhaseInfo.CurrentPhaseStart,
		CurrentPhaseLastLevel:  result.newPhaseInfo.CurrentPhaseEnd,
		PhaseProfitMultiplier:  result.newPhaseInfo.ProfitMultiplier,

		// Table count
		NewTableCount: result.newTableCount,

		// Rewards
		PhaseRewards:    result.phaseRewards,
		AllPhaseRewards: result.allPhaseRewardsCollected,
	}
}
