package handlers

import (
	"github.com/gofiber/fiber/v2"
	"github.com/winartodev/cat-cafe/internal/dto"
	"github.com/winartodev/cat-cafe/internal/middleware"
	"github.com/winartodev/cat-cafe/internal/usecase"
	"github.com/winartodev/cat-cafe/pkg/apperror"
	"github.com/winartodev/cat-cafe/pkg/helper"
	"github.com/winartodev/cat-cafe/pkg/response"
)

type GameHandler struct {
	GameUseCase usecase.GameUseCase
}

func NewGameHandler(gameUc usecase.GameUseCase) *GameHandler {
	return &GameHandler{
		GameUseCase: gameUc,
	}
}

func (g *GameHandler) SyncBalance(c *fiber.Ctx) error {
	var req dto.SyncBalanceRequest
	if err := c.BodyParser(&req); err != nil {
		return response.FailedResponse(c, fiber.StatusBadRequest, err)
	}
	ctx := c.Context()

	userID, _ := helper.GetUserIDFromContext(ctx)
	if userID <= 0 {
		return response.FailedResponse(c, fiber.StatusUnauthorized, apperror.ErrUnauthorized)
	}

	res, err := g.GameUseCase.UpdateUserBalance(ctx, req.CoinsEarned, userID)
	if err != nil {
		return response.FailedResponse(c, fiber.StatusInternalServerError, err)
	}

	return response.SuccessResponse(c, fiber.StatusOK, "Sync Balance Successfully", res, nil)
}

func (g *GameHandler) Route(r fiber.Router, m middleware.Middleware) error {
	game := r.Group("/game")
	game.Post("/sync-balance", m.WithUserAuth(g.SyncBalance))

	return nil
}
