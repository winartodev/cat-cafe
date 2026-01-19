package usecase

import (
	"github.com/winartodev/cat-cafe/internal/repositories"
	"github.com/winartodev/cat-cafe/pkg/jwt"
)

type UseCase struct {
	UserUseCase        UserUseCase
	RewardUseCase      RewardUseCase
	DailyRewardUseCase DailyRewardUseCase
	AuthUseCase        AuthUseCase
	GameUseCase        GameUseCase
}

func SetUpUseCase(repo repositories.Repository, jwt_ *jwt.JWT) *UseCase {
	userUC := NewUserUseCase(
		repo.UserRepository,
		repo.UserDailyRewardRepository,
	)

	rewardUC := NewRewardUseCase(
		repo.RewardRepository,
	)

	dailyRewardUC := NewDailyRewardUseCase(
		repo.DailyRewardRepository,
		repo.UserDailyRewardRepository,
		repo.UserRepository,
		userUC,
		rewardUC,
	)
	
	gameUC := NewGameUseCase(
		userUC,
		repo.UserRepository,
	)

	authUC := NewAuthUseCase(
		userUC,
		gameUC,
		repo.UserRepository,
		jwt_,
	)

	return &UseCase{
		UserUseCase:        userUC,
		RewardUseCase:      rewardUC,
		DailyRewardUseCase: dailyRewardUC,
		AuthUseCase:        authUC,
		GameUseCase:        gameUC,
	}
}
