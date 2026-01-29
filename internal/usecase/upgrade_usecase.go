package usecase

import (
	"context"
	"fmt"

	"github.com/winartodev/cat-cafe/internal/entities"
	"github.com/winartodev/cat-cafe/internal/repositories"
	"github.com/winartodev/cat-cafe/pkg/apperror"
)

type UpgradeUseCase interface {
	CreateUpgrade(ctx context.Context, upgrade entities.Upgrade) (err error)
}

type upgradeUseCase struct {
	upgradeRepo  repositories.UpgradeRepository
	foodItemRepo repositories.FoodItemRepository
}

func NewUpgradeUseCase(
	upgradeRepo repositories.UpgradeRepository,
	foodItemRepo repositories.FoodItemRepository,
) UpgradeUseCase {
	return &upgradeUseCase{
		upgradeRepo:  upgradeRepo,
		foodItemRepo: foodItemRepo,
	}
}

func (u *upgradeUseCase) CreateUpgrade(ctx context.Context, upgrade entities.Upgrade) (err error) {
	var targetID int64
	if u.isEffectTargetFood(upgrade.Effect.Target) {
		foodItem, err := u.foodItemRepo.GetFoodBySlugDB(ctx, upgrade.Effect.TargetName)
		if err != nil {
			return err
		}

		if foodItem == nil {
			return apperror.ErrorNotFound("food item:", upgrade.Effect.TargetName)
		}

		targetID = foodItem.ID
	}

	fmt.Println(targetID)
	id, err := u.upgradeRepo.CreateUpgradeDB(ctx, upgrade)
	if err != nil {
		return err
	}

	upgrade.ID = *id

	return nil
}

func (u *upgradeUseCase) isEffectTargetFood(target entities.UpgradeEffectTarget) bool {
	return target == entities.UpgradeEffectTargetFood
}
