package usecase

import (
	"context"
	"github.com/winartodev/cat-cafe/internal/entities"
	"github.com/winartodev/cat-cafe/internal/repositories"
	"github.com/winartodev/cat-cafe/pkg/apperror"
)

type UserUseCase interface {
	GetUserByID(ctx context.Context, id int64) (res *entities.User, err error)
	GetUserDailyRewardByID(ctx context.Context, id int64) (res *entities.UserDailyReward, err error)
	GetUserBalance(ctx context.Context, id int64) (res *entities.UserBalance, err error)
}

type userUseCase struct {
	userRepo            repositories.UserRepository
	userDailyRewardRepo repositories.UserDailyRewardRepository
}

func NewUserUseCase(userRepo repositories.UserRepository, userDailyRewardRepo repositories.UserDailyRewardRepository) UserUseCase {
	return &userUseCase{
		userRepo:            userRepo,
		userDailyRewardRepo: userDailyRewardRepo,
	}
}

func (c *userUseCase) GetUserByID(ctx context.Context, id int64) (res *entities.User, err error) {
	user, err := c.userRepo.GetUserByIDDB(ctx, id)
	if err != nil {
		return nil, err
	}

	if user == nil {
		return nil, apperror.ErrRecordNotFound
	}

	return user, nil
}

func (c *userUseCase) GetUserDailyRewardByID(ctx context.Context, id int64) (res *entities.UserDailyReward, err error) {
	user, err := c.GetUserByID(ctx, id)
	if err != nil {
		return nil, err
	}

	progress, err := c.userDailyRewardRepo.GetUserDailyRewardByIDDB(ctx, user.ID)
	if err != nil {
		return nil, err
	}

	if progress == nil {
		return nil, nil
	}

	return progress, nil
}

func (c *userUseCase) GetUserBalance(ctx context.Context, id int64) (res *entities.UserBalance, err error) {
	userBalance, err := c.userRepo.GetUserBalanceByIDDB(ctx, id)
	if err != nil {
		return nil, err
	}

	if userBalance == nil {
		return nil, apperror.ErrRecordNotFound
	}

	return userBalance, nil
}
