package usecase

import (
	"context"
	"database/sql"

	"github.com/winartodev/cat-cafe/internal/entities"
	"github.com/winartodev/cat-cafe/internal/repositories"
	"github.com/winartodev/cat-cafe/pkg/apperror"
)

type FoodItemUseCase interface {
	CreateFood(ctx context.Context, data entities.FoodItem, overrideLevels []entities.FoodItemOverrideLevel) (res *entities.FoodItem, err error)
	UpdateFood(ctx context.Context, id int64, data entities.FoodItem, overrideLevels []entities.FoodItemOverrideLevel) (res *entities.FoodItem, err error)
	GetFoodBySlug(ctx context.Context, slug string) (res *entities.FoodItem, overrideLevels []entities.FoodItemOverrideLevel, err error)
	GetFoods(ctx context.Context, limit, offset int) (res []entities.FoodItem, totalRow int64, err error)
	GetFoodByID(ctx context.Context, id int64) (res *entities.FoodItem, overrideLevels []entities.FoodItemOverrideLevel, err error)
}

type foodItemUseCase struct {
	foodItemRepo repositories.FoodItemRepository
}

func NewFoodItemUseCase(foodItemRepo repositories.FoodItemRepository) FoodItemUseCase {
	return &foodItemUseCase{
		foodItemRepo: foodItemRepo,
	}
}

func (u *foodItemUseCase) CreateFood(ctx context.Context, data entities.FoodItem, overrideLevels []entities.FoodItemOverrideLevel) (res *entities.FoodItem, err error) {
	err = u.foodItemRepo.FoodItemWithTx(ctx, func(tx *sql.Tx) error {
		foodItemTx := u.foodItemRepo.WithTx(tx)
		id, err := foodItemTx.CreateFoodDB(ctx, data)
		if err != nil {
			return err
		}

		data.ID = *id

		if overrideLevels != nil {
			err = foodItemTx.CreateOverrideLevelDB(ctx, data.ID, overrideLevels)
			if err != nil {
				return err
			}
		}
		return nil
	})

	if err != nil {
		return nil, err
	}

	return &data, nil
}

func (u *foodItemUseCase) UpdateFood(ctx context.Context, id int64, data entities.FoodItem, overrideLevels []entities.FoodItemOverrideLevel) (res *entities.FoodItem, err error) {
	err = u.foodItemRepo.FoodItemWithTx(ctx, func(tx *sql.Tx) error {
		foodItemTx := u.foodItemRepo.WithTx(tx)
		err = foodItemTx.UpdateFoodDB(ctx, id, data)
		if err != nil {
			return err
		}

		if overrideLevels != nil {
			err = foodItemTx.DeleteOverrideLevelDB(ctx, id)
			if err != nil {
				return err
			}

			err = foodItemTx.CreateOverrideLevelDB(ctx, id, overrideLevels)
			if err != nil {
				return err
			}
		}
		return nil
	})

	if err != nil {
		return nil, err
	}

	return &data, nil
}

func (u *foodItemUseCase) GetFoodBySlug(ctx context.Context, slug string) (res *entities.FoodItem, overrideLevels []entities.FoodItemOverrideLevel, err error) {
	foodItem, err := u.foodItemRepo.GetFoodBySlugDB(ctx, slug)
	if err != nil {
		return nil, nil, err
	}

	if foodItem == nil {
		return nil, nil, apperror.ErrRecordNotFound
	}

	foodOverrideLevels, err := u.foodItemRepo.GetOverrideLevelDB(ctx, foodItem.ID)
	if err != nil {
		return nil, nil, err
	}

	if foodOverrideLevels == nil {
		foodOverrideLevels = []entities.FoodItemOverrideLevel{}
	}

	return foodItem, foodOverrideLevels, nil
}

func (u *foodItemUseCase) GetFoods(ctx context.Context, limit, offset int) (res []entities.FoodItem, totalRow int64, err error) {
	res, err = u.foodItemRepo.GetFoodsDB(ctx, limit, offset)
	if err != nil {
		return nil, 0, err
	}

	totalRow, err = u.foodItemRepo.CountFoodItemDB(ctx)
	if err != nil {
		return nil, 0, err
	}

	return res, totalRow, err
}

func (u *foodItemUseCase) GetFoodByID(ctx context.Context, id int64) (res *entities.FoodItem, overrideLevels []entities.FoodItemOverrideLevel, err error) {
	foodItem, err := u.foodItemRepo.GetFoodByIDDB(ctx, id)
	if err != nil {
		return nil, nil, err
	}

	if foodItem == nil {
		return nil, nil, apperror.ErrRecordNotFound
	}

	foodOverrideLevels, err := u.foodItemRepo.GetOverrideLevelDB(ctx, foodItem.ID)
	if err != nil {
		return nil, nil, err
	}

	if foodOverrideLevels == nil {
		foodOverrideLevels = []entities.FoodItemOverrideLevel{}
	}

	return foodItem, foodOverrideLevels, nil
}
