package usecase

import (
	"context"
	"github.com/winartodev/cat-cafe/internal/entities"
	"github.com/winartodev/cat-cafe/internal/repositories"
	"github.com/winartodev/cat-cafe/pkg/apperror"
	"github.com/winartodev/cat-cafe/pkg/helper"
	"time"
)

type UserUseCase interface {
	CreateUser(ctx context.Context, data entities.User) (res *entities.User, err error)
	GetUserByID(ctx context.Context, userID int64) (res *entities.User, err error)
	GetUserDailyRewardByID(ctx context.Context, userID int64) (res *entities.UserDailyReward, err error)
	GetUserBalance(ctx context.Context, userID int64) (res *entities.UserBalance, err error)
	GetUserByEmail(ctx context.Context, email string) (res *entities.User, err error)
	IsDailyRewardAvailable(ctx context.Context, userID int64) (isAvailable bool, err error)
}

type userUseCase struct {
	userRepo        repositories.UserRepository
	userProgression repositories.UserProgressionRepository
}

func NewUserUseCase(
	userRepo repositories.UserRepository,
	userProgression repositories.UserProgressionRepository,
) UserUseCase {
	return &userUseCase{
		userRepo:        userRepo,
		userProgression: userProgression,
	}
}

func (c *userUseCase) CreateUser(ctx context.Context, data entities.User) (res *entities.User, err error) {
	id, err := c.userRepo.CreateUserDB(ctx, &data)
	if err != nil {
		return nil, err
	}

	if id == nil {
		return nil, apperror.ErrFailedRetrieveID
	}

	data.ID = *id

	return &data, err
}

func (c *userUseCase) GetUserByID(ctx context.Context, userID int64) (res *entities.User, err error) {
	//userCache, err := c.userRepo.GetUserRedis(ctx, userID)
	//if err == nil && userCache != nil {
	//	return entities.UserFromCache(userCache), nil
	//}

	user, err := c.userRepo.GetUserByIDDB(ctx, userID)
	if err != nil {
		return nil, err
	}

	if user == nil {
		return nil, apperror.ErrRecordNotFound
	}

	//go func(u *entities.User) {
	//	_ = c.userRepo.SetUserRedis(context.Background(), userID, u.ToCache(), 24*time.Hour)
	//}(user)

	return user, nil
}

func (c *userUseCase) GetUserDailyRewardByID(ctx context.Context, userID int64) (res *entities.UserDailyReward, err error) {
	//progress, err := c.userProgressionRepo.GetUserDailyRewardRedis(ctx, userID)
	//if err == nil && progress != nil {
	//	return progress, nil
	//}

	user, err := c.GetUserByID(ctx, userID)
	if err != nil {
		return nil, err
	}

	progress, err := c.userProgression.GetUserDailyRewardByIDDB(ctx, user.ID)
	if err != nil {
		return nil, err
	}

	if progress == nil {
		return nil, nil
	}

	//go func(p *entities.UserDailyReward) {
	//	_ = c.userProgressionRepo.SetUserDailyRewardRedis(context.Background(), userID, p, 24*time.Hour)
	//}(progress)

	return progress, nil
}

func (c *userUseCase) GetUserBalance(ctx context.Context, userID int64) (res *entities.UserBalance, err error) {
	userBalance, err := c.userRepo.GetUserBalanceByIDDB(ctx, userID)
	if err != nil {
		return nil, err
	}

	if userBalance == nil {
		return nil, apperror.ErrRecordNotFound
	}

	return userBalance, nil
}

func (c *userUseCase) GetUserByEmail(ctx context.Context, email string) (res *entities.User, err error) {
	if !helper.IsEmailValid(email) {
		return nil, apperror.ErrInvalidEmail
	}

	user, err := c.userRepo.GetUserByEmailDB(ctx, email)
	if err != nil {
		return nil, err
	}

	if user == nil {
		return nil, apperror.ErrRecordNotFound
	}

	return user, nil
}

func (c *userUseCase) IsDailyRewardAvailable(ctx context.Context, userID int64) (isAvailable bool, err error) {
	progression, err := c.GetUserDailyRewardByID(ctx, userID)
	if err != nil {
		return false, err
	}

	now := helper.NowUTC()
	today := now.Truncate(24 * time.Hour)

	if progression == nil {
		progression = &entities.UserDailyReward{LongestStreak: 0, CurrentStreak: 0, LastClaimDate: nil}
	} else if progression.LastClaimDate != nil {
		// Check if user already claimed today
		lastClaim := progression.LastClaimDate.UTC().Truncate(24 * time.Hour)
		if today.Equal(lastClaim) {
			return false, nil
		}
	}

	return true, nil
}
