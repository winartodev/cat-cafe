package handlers

import (
	"fmt"

	"github.com/winartodev/cat-cafe/internal/middleware"

	"github.com/gofiber/fiber/v2"
	"github.com/winartodev/cat-cafe/internal/usecase"
)

type Registerer interface {
	Route(open fiber.Router, userAuth fiber.Router, internalAuth fiber.Router) error
}

func register(open fiber.Router, userAuth fiber.Router, internalAuth fiber.Router, items ...Registerer) error {
	for _, item := range items {
		if err := item.Route(open, userAuth, internalAuth); err != nil {
			return fmt.Errorf("route registration failed: %w", err)
		}
	}
	return nil
}

func SetupHandler(app *fiber.App, uc usecase.UseCase, middleware middleware.Middleware) {
	rewardHandler := NewRewardHandler(
		uc.RewardUseCase,
		uc.DailyRewardUseCase,
	)

	foodItemHandler := NewFoodItemHandler(
		uc.FoodItemUseCase,
	)

	gameStageHandler := NewGameStageHandler(
		uc.GameStageUseCase,
	)

	authHandler := NewAuthHandler(
		uc.AuthUseCase,
	)

	gameHandler := NewGameHandler(
		uc.GameUseCase,
		uc.DailyRewardUseCase,
	)

	upgradeHandler := NewUpgradeHandler(
		uc.UpgradeUseCase,
	)

	tutorialHandler := NewTutorialHandler(
		uc.TutorialUseCase,
	)

	api := app.Group("/api")
	userAuth := api.Group("/v1", middleware.WithUserAuth())
	internalAuth := api.Group("/internal")

	if err := register(api, userAuth, internalAuth,
		rewardHandler,
		foodItemHandler,
		authHandler,
		gameHandler,
		gameStageHandler,
		upgradeHandler,
		tutorialHandler,
	); err != nil {
		panic(err)
	}
}
