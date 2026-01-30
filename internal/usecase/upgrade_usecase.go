package usecase

import (
	"context"

	"github.com/winartodev/cat-cafe/internal/entities"
	"github.com/winartodev/cat-cafe/internal/repositories"
	"github.com/winartodev/cat-cafe/pkg/apperror"
)

type UpgradeUseCase interface {
	CreateUpgrade(ctx context.Context, upgrade *entities.Upgrade) (res *entities.Upgrade, err error)
	GetUpgrades(ctx context.Context, limit, offset int) (res []entities.Upgrade, totalRows int64, err error)
	GetUpgradeByID(ctx context.Context, id int64) (res *entities.Upgrade, err error)
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

func (u *upgradeUseCase) CreateUpgrade(ctx context.Context, upgrade *entities.Upgrade) (res *entities.Upgrade, err error) {
	if u.isEffectTargetFood(upgrade.Effect.Target) {
		foodItem, err := u.foodItemRepo.GetFoodBySlugDB(ctx, upgrade.Effect.TargetName)
		if err != nil {
			return nil, err
		}

		if foodItem == nil {
			return nil, apperror.ErrorNotFound("food item:", upgrade.Effect.TargetName)
		}

		upgrade.Effect.TargetID = foodItem.ID
	}

	id, err := u.upgradeRepo.CreateUpgradeDB(ctx, *upgrade)
	if err != nil {
		return nil, err
	}

	upgrade.ID = *id

	return upgrade, nil
}

func (u *upgradeUseCase) GetUpgrades(ctx context.Context, limit, offset int) (res []entities.Upgrade, totalRows int64, err error) {
	res, err = u.upgradeRepo.GetUpgradesDB(ctx, limit, offset)
	if err != nil {
		return nil, 0, err
	}

	totalRows, err = u.upgradeRepo.CountUpgradesDB(ctx)
	if err != nil {
		return nil, 0, err
	}

	return res, totalRows, nil
}

// GetUpgradeByID implements UpgradeUseCase.
func (u *upgradeUseCase) GetUpgradeByID(ctx context.Context, id int64) (res *entities.Upgrade, err error) {
	return u.upgradeRepo.GetUpgradeByIDDB(ctx, id)
}

func (u *upgradeUseCase) isEffectTargetFood(target entities.UpgradeEffectTarget) bool {
	return target == entities.UpgradeEffectTargetFood
}
