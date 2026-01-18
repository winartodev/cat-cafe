package usecase

import (
	"github.com/winartodev/cat-cafe/internal/repositories"
	"github.com/winartodev/cat-cafe/pkg/jwt"
)

type UseCase struct {
	UserUseCase        UserUseCase
	DailyRewardUseCase DailyRewardUseCase
	AuthUseCase        AuthUseCase
	GameUseCase        GameUseCase
}

func SetUpUseCase(repo repositories.Repository, jwt_ *jwt.JWT) *UseCase {
	userUC := NewUserUseCase(
		repo.UserRepository,
		repo.UserDailyRewardRepository,
	)

	dailyRewardUC := NewDailyRewardUseCase(
		repo.DailyRewardRepository,
		repo.UserDailyRewardRepository,
		repo.UserRepository,
		userUC,
	)

	authUC := NewAuthUseCase(
		userUC,
		repo.UserRepository,
		jwt_,
	)

	gameUC := NewGameUseCase(
		userUC,
		repo.UserRepository,
	)

	return &UseCase{
		UserUseCase:        userUC,
		DailyRewardUseCase: dailyRewardUC,
		AuthUseCase:        authUC,
		GameUseCase:        gameUC,
	}
}
