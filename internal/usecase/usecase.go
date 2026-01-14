package usecase

import "github.com/winartodev/cat-cafe/internal/repositories"

type UseCase struct {
	UserUseCase        UserUseCase
	DailyRewardUseCase DailyRewardUseCase
}

func SetUpUseCase(repo repositories.Repository) *UseCase {
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

	return &UseCase{
		UserUseCase:        userUC,
		DailyRewardUseCase: dailyRewardUC,
	}
}
