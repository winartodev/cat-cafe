package usecase

import (
	"context"
	"github.com/winartodev/cat-cafe/internal/dto"
	"github.com/winartodev/cat-cafe/internal/entities"
	"github.com/winartodev/cat-cafe/internal/repositories"
	"github.com/winartodev/cat-cafe/pkg/helper"
)

type GameUseCase interface {
	UpdateUserBalance(ctx context.Context, coinEarned int64, userID int64) (res *dto.UserBalanceResponse, err error)
	GetUserGameData(ctx context.Context, userID int64) (res *entities.Game, err error)
}

type gameUseCase struct {
	userUseCase UserUseCase
	userRepo    repositories.UserRepository
}

func NewGameUseCase(userUc UserUseCase, userRepo repositories.UserRepository) GameUseCase {
	return &gameUseCase{
		userUseCase: userUc,
		userRepo:    userRepo,
	}
}

func (g *gameUseCase) UpdateUserBalance(ctx context.Context, coinEarned int64, userID int64) (res *dto.UserBalanceResponse, err error) {
	// TODO: should we validate earning rate ???

	err = g.userRepo.BalanceWithTx(ctx, func(txRepo repositories.UserRepository) error {
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
