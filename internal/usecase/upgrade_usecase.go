package usecase

import (
	"context"

	"github.com/winartodev/cat-cafe/internal/entities"
	"github.com/winartodev/cat-cafe/internal/repositories"
	"github.com/winartodev/cat-cafe/pkg/apperror"
)

type UpgradeUseCase interface {
	CreateUpgrade(ctx context.Context, data entities.Upgrade) (res *entities.Upgrade, err error)
	GetUpgrades(ctx context.Context, limit, offset int) (res []entities.Upgrade, totalRows int64, err error)
	GetUpgradeByID(ctx context.Context, id int64) (res *entities.Upgrade, err error)
	UpdateUpgrade(ctx context.Context, id int64, data entities.Upgrade) (res *entities.Upgrade, err error)
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

func (u *upgradeUseCase) CreateUpgrade(ctx context.Context, data entities.Upgrade) (res *entities.Upgrade, err error) {
	if err := u.resolveEffectTargetID(ctx, &data); err != nil {
		return nil, err
	}

	id, err := u.upgradeRepo.CreateUpgradeDB(ctx, data)
	if err != nil {
		return nil, err
	}

	data.ID = *id

	return &data, nil
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

func (u *upgradeUseCase) UpdateUpgrade(ctx context.Context, id int64, data entities.Upgrade) (res *entities.Upgrade, err error) {
	if err := u.resolveEffectTargetID(ctx, &data); err != nil {
		return nil, err
	}

	err = u.upgradeRepo.UpdateUpgradeDB(ctx, id, data)
	if err != nil {
		return nil, err
	}

	upgrade, err := u.upgradeRepo.GetUpgradeByIDDB(ctx, id)
	if err != nil {
		return nil, err
	}

	return upgrade, nil
}

func (u *upgradeUseCase) resolveEffectTargetID(ctx context.Context, data *entities.Upgrade) (err error) {
	target := data.Effect.Target
	switch target {
	case entities.UpgradeEffectTargetFood:
		foodItem, err := u.foodItemRepo.GetFoodBySlugDB(ctx, data.Effect.TargetName)
		if err != nil {
			return err
		}

		if foodItem == nil {
			return apperror.ErrorNotFound("food item:", data.Effect.TargetName)
		}

		data.Effect.TargetID = foodItem.ID

		return nil
	}

	return nil
}
