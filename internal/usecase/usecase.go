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
	GameStageUseCase   GameStageUseCase
	FoodItemUseCase    FoodItemUseCase
}

func SetUpUseCase(repo repositories.Repository, jwt_ *jwt.JWT) *UseCase {
	userUC := NewUserUseCase(
		repo.UserRepository,
		repo.UserProgressionRepository,
	)

	rewardUC := NewRewardUseCase(
		repo.RewardRepository,
	)

	dailyRewardUC := NewDailyRewardUseCase(
		repo.DailyRewardRepository,
		repo.UserProgressionRepository,
		repo.UserRepository,
		userUC,
		rewardUC,
	)

	foodItemUC := NewFoodItemUseCase(
		repo.FoodItemRepository,
	)

	gameStageUC := NewGameStageUseCase(
		repo.GameStageRepository,
		repo.StageCustomerConfigRepository,
		repo.StageStaffConfigRepository,
		repo.StageKitchenConfigRepository,
		repo.StageCameraConfigRepository,
		repo.RewardRepository,
		repo.KitchenStationRepository,
		repo.FoodItemRepository,
	)

	gameUC := NewGameUseCase(
		userUC,
		repo.UserRepository,
		repo.UserProgressionRepository,
		repo.GameStageRepository,
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
		GameStageUseCase:   gameStageUC,
		FoodItemUseCase:    foodItemUC,
	}
}
