package usecase

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"github.com/winartodev/cat-cafe/internal/entities"
	"github.com/winartodev/cat-cafe/pkg/apperror"
)

type stageUpgradeContext struct {
	userID              int64
	stageID             int64
	slug                string
	stageUpgrade        *entities.StageUpgrade
	userBalance         *entities.UserBalance
	balanceType         entities.UserBalanceType
	userProgress        *entities.UserKitchenStageProgression
	activeStageUpgrades []entities.UserStageUpgrade
}

func (g *gameUseCase) gatherStageUpgradeData(ctx context.Context, userID int64, slug string) (res *stageUpgradeContext, err error) {
	suctx := &stageUpgradeContext{
		userID: userID,
		slug:   slug,
	}

	latestStage, err := g.userProgressionUseCase.LatestStageProgression(ctx)
	if err != nil {
		return nil, err
	}

	suctx.stageID = latestStage.StageID

	suctx.activeStageUpgrades, err = g.userProgressionUseCase.GetActiveStageUpgrade(ctx, suctx.stageID)
	if err != nil {
		return nil, err
	}

	suctx.stageUpgrade, err = g.stageUpgradeRepo.GetUpgradeByStageIDAndSlugDB(ctx, suctx.stageID, slug)
	if err != nil {
		return nil, err
	}

	if suctx.stageUpgrade == nil {
		return nil, apperror.ErrRecordNotFound
	}

	suctx.userProgress, err = g.userProgressionRepo.GetUserKitchenProgressDB(ctx, userID, suctx.stageID)
	if err != nil {
		return nil, err
	}

	suctx.userBalance, err = g.userUseCase.GetUserBalance(ctx, userID)
	if err != nil {
		return nil, err
	}

	if suctx.userBalance == nil {
		return nil, apperror.ErrRecordNotFound
	}

	return suctx, nil
}

func (g *gameUseCase) validateStageUpgrade(ctx *stageUpgradeContext) error {
	if ctx.stageUpgrade == nil {
		return apperror.ErrorInvalidRequest("stage upgrade is empty")
	}

	for _, stageUpgrade := range ctx.activeStageUpgrades {
		if stageUpgrade.Upgrade.Slug == ctx.stageUpgrade.Upgrade.Slug && stageUpgrade.IsPurchased {
			return apperror.ErrorInvalidRequest(fmt.Sprintf("upgrade with name %s already purchased", ctx.stageUpgrade.Upgrade.Name))
		}
	}

	upgradeCost := ctx.stageUpgrade.Upgrade.Cost
	userBalance := ctx.userBalance

	switch ctx.stageUpgrade.Upgrade.CostType {
	case entities.UpgradeCostTypeCoin:
		ctx.balanceType = entities.BalanceTypeCoin
		if userBalance.Coin < upgradeCost {
			return apperror.ErrInsufficientCoins
		}

	case entities.UpgradeCostTypeGem:
		ctx.balanceType = entities.BalanceTypeGem
		if userBalance.Gem < upgradeCost {
			return apperror.ErrInsufficientGems
		}
	}

	return nil
}

func (g *gameUseCase) executeStageUpgradeTransaction(ctx context.Context, upgradeContext *stageUpgradeContext) (err error) {
	return g.userProgressionRepo.WithUserProgressionTx(ctx, func(tx *sql.Tx) error {
		upgrade := upgradeContext.stageUpgrade.Upgrade
		upgradeEffect := upgrade.Effect
		userID := upgradeContext.userID
		stageID := upgradeContext.stageID

		userProgressionTx := g.userProgressionRepo.WithTx(tx)
		userRepoTx := g.userRepo.WithTx(tx)

		err = userRepoTx.UpdateUserBalanceWithTx(ctx, userID, upgradeContext.balanceType, -upgrade.Cost)
		if err != nil {
			return err
		}

		err = userProgressionTx.CreateUpgradeStageProgression(ctx, entities.UserStageUpgrade{
			UserID:             userID,
			StageID:            stageID,
			GameStageUpgradeID: upgradeContext.stageUpgrade.ID,
		})
		if errors.Is(err, apperror.ErrConflict) {
			return apperror.ErrorInvalidRequest(fmt.Sprintf("upgrade with name %s already purchased", upgrade.Name))
		} else if err != nil {
			return err
		}

		if upgradeEffect.Target == upgradeTargetFood {
			progress := upgradeContext.userProgress
			if progress == nil {
				progress = &entities.UserKitchenStageProgression{
					UserID:          userID,
					StageID:         stageID,
					StationUpgrades: make(map[string]entities.UserStationUpgrade),
				}
			}

			stationUpgrade := progress.StationUpgrades
			if stationUpgrade == nil {
				stationUpgrade = make(map[string]entities.UserStationUpgrade)
			}

			foodSlug := upgradeEffect.TargetName
			currentBonus := stationUpgrade[foodSlug]

			if currentBonus.ProfitBonus == 0 {
				currentBonus.ProfitBonus = 1.0
			}

			if currentBonus.ReduceCookingTime == 0 && (upgradeEffect.Unit == entities.UpgradeEffectUnitMultiplier || upgradeEffect.Unit == entities.UpgradeEffectUnitPercentage) {
				currentBonus.ReduceCookingTime = 1.0
			}

			switch upgradeEffect.Type {
			case entities.UpgradeEffectTypeAddHelper:
				currentBonus.HelperCount += int64(upgradeEffect.Value)

			case entities.UpgradeEffectTypeAddCustomer:
				currentBonus.CustomerCount += int64(upgradeEffect.Value)

			case entities.UpgradeEffectTypeReduceCookingTime:
				currentBonus.ReduceCookingTime = upgradeEffect.CalculateNewValue(currentBonus.ReduceCookingTime)

			case entities.UpgradeEffectTypeProfit:
				currentBonus.ProfitBonus = upgradeEffect.CalculateNewValue(currentBonus.ProfitBonus)
			}

			stationUpgrade[foodSlug] = currentBonus
			progress.StationUpgrades = stationUpgrade

			return userProgressionTx.UpdateKitchenStationUpgradeDB(ctx, userID, stageID, stationUpgrade)
		}

		return nil
	})
}
