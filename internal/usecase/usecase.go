package usecase

import (
	"github.com/winartodev/cat-cafe/internal/repositories"
	"github.com/winartodev/cat-cafe/pkg/jwt"
)

type UseCase struct {
	UserUseCase            UserUseCase
	UserProgressionUseCase UserProgressionUseCase
	RewardUseCase          RewardUseCase
	DailyRewardUseCase     DailyRewardUseCase
	AuthUseCase            AuthUseCase
	GameUseCase            GameUseCase
	GameStageUseCase       GameStageUseCase
	FoodItemUseCase        FoodItemUseCase
	UpgradeUseCase         UpgradeUseCase
}

func SetUpUseCase(repo repositories.Repository, jwt_ *jwt.JWT) *UseCase {
	userProgressionUC := NewUserProgressionUseCase(
		repo.UserProgressionRepository,
		repo.FoodItemRepository,
	)

	userUC := NewUserUseCase(
		repo.UserRepository,
		userProgressionUC,
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

	upgradeUC := NewUpgradeUseCase(
		repo.UpgradeRepository,
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
		repo.UpgradeRepository,
		repo.StageUpgradeRepository,
	)

	gameUC := NewGameUseCase(
		userUC,
		userProgressionUC,
		repo.UserRepository,
		repo.UserProgressionRepository,
		repo.GameStageRepository,
		repo.FoodItemRepository,
		repo.KitchenStationRepository,
		repo.StageKitchenConfigRepository,
		repo.RewardRepository,
		repo.StageUpgradeRepository,
	)

	authUC := NewAuthUseCase(
		userUC,
		gameUC,
		repo.UserRepository,
		jwt_,
	)

	return &UseCase{
		UserUseCase:            userUC,
		UserProgressionUseCase: userProgressionUC,
		RewardUseCase:          rewardUC,
		DailyRewardUseCase:     dailyRewardUC,
		AuthUseCase:            authUC,
		GameUseCase:            gameUC,
		GameStageUseCase:       gameStageUC,
		FoodItemUseCase:        foodItemUC,
		UpgradeUseCase:         upgradeUC,
	}
}
