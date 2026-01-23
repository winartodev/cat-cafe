package usecase

import (
	"context"
	"github.com/winartodev/cat-cafe/internal/entities"
	"github.com/winartodev/cat-cafe/internal/repositories"
	"github.com/winartodev/cat-cafe/pkg/apperror"
)

type FoodItemUseCase interface {
	CreateFood(ctx context.Context, data entities.FoodItem) (res *entities.FoodItem, err error)
	UpdateFood(ctx context.Context, id int64, data entities.FoodItem) (res *entities.FoodItem, err error)
	GetFoodBySlug(ctx context.Context, slug string) (res *entities.FoodItem, err error)
	GetFoods(ctx context.Context, limit, offset int) (res []entities.FoodItem, totalRow int64, err error)
	GetFoodByID(ctx context.Context, id int64) (res *entities.FoodItem, err error)
}

type foodItemUseCase struct {
	foodItemRepo repositories.FoodItemRepository
}

func NewFoodItemUseCase(foodItemRepo repositories.FoodItemRepository) FoodItemUseCase {
	return &foodItemUseCase{
		foodItemRepo: foodItemRepo,
	}
}

func (u *foodItemUseCase) CreateFood(ctx context.Context, data entities.FoodItem) (res *entities.FoodItem, err error) {
	id, err := u.foodItemRepo.CreateFoodDB(ctx, data)
	if err != nil {
		return nil, err
	}

	if id == nil {
		return nil, apperror.ErrFailedRetrieveID
	}

	data.ID = *id

	return &data, err
}

func (u *foodItemUseCase) UpdateFood(ctx context.Context, id int64, data entities.FoodItem) (res *entities.FoodItem, err error) {
	err = u.foodItemRepo.UpdateFoodDB(ctx, id, data)
	if err != nil {
		return nil, err
	}

	res, err = u.GetFoodByID(ctx, id)
	if err != nil {
		return nil, err
	}

	return res, err
}

func (u *foodItemUseCase) GetFoodBySlug(ctx context.Context, slug string) (res *entities.FoodItem, err error) {
	reward, err := u.foodItemRepo.GetFoodBySlugDB(ctx, slug)
	if err != nil {
		return nil, err
	}

	if reward == nil {
		return nil, apperror.ErrRecordNotFound
	}

	return reward, nil
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

func (u *foodItemUseCase) GetFoodByID(ctx context.Context, id int64) (res *entities.FoodItem, err error) {
	reward, err := u.foodItemRepo.GetFoodByIDDB(ctx, id)
	if err != nil {
		return nil, err
	}

	if reward == nil {
		return nil, apperror.ErrRecordNotFound
	}

	return reward, nil
}
