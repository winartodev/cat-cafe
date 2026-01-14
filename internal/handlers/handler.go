package handlers

import (
	"fmt"

	"github.com/gofiber/fiber/v2"
	"github.com/winartodev/cat-cafe/internal/controllers"
)

type Registerer interface {
	Route(route fiber.Router) error
}

func register(api fiber.Router, items ...Registerer) error {
	for _, item := range items {
		if err := item.Route(api); err != nil {
			return fmt.Errorf("route registration failed: %w", err)
		}
	}
	return nil
}

func SetupHandler(app *fiber.App, ctrl controllers.Controller) {
	catHandler := NewCatHandler(
		ctrl.CatController,
	)

	api := app.Group("/api")

	if err := register(api, catHandler); err != nil {
		panic(err)
	}
}
