package handlers

import (
	"fmt"
	"github.com/winartodev/cat-cafe/internal/middleware"

	"github.com/gofiber/fiber/v2"
	"github.com/winartodev/cat-cafe/internal/usecase"
)

type Registerer interface {
	Route(r fiber.Router, m middleware.Middleware) error
}

func register(api fiber.Router, m middleware.Middleware, items ...Registerer) error {
	for _, item := range items {
		if err := item.Route(api, m); err != nil {
			return fmt.Errorf("route registration failed: %w", err)
		}
	}
	return nil
}

func SetupHandler(app *fiber.App, uc usecase.UseCase, m middleware.Middleware) {
	dailyRewardHandler := NewDailyRewardHandler(
		uc.DailyRewardUseCase,
	)

	authHandler := NewAuthHandler(
		uc.AuthUseCase,
	)

	gameHandler := NewGameHandler(
		uc.GameUseCase,
	)

	api := app.Group("/api")

	if err := register(api, m,
		dailyRewardHandler,
		authHandler,
		gameHandler,
	); err != nil {
		panic(err)
	}
}
