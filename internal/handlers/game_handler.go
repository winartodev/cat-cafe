package handlers

import (
	"context"
	"errors"
	"github.com/gofiber/fiber/v2"
	"github.com/winartodev/cat-cafe/internal/dto"
	"github.com/winartodev/cat-cafe/internal/usecase"
	"github.com/winartodev/cat-cafe/pkg/apperror"
	"github.com/winartodev/cat-cafe/pkg/helper"
	"github.com/winartodev/cat-cafe/pkg/response"
)

// GameHandler is used for interaction with public players
type GameHandler struct {
	GameUseCase        usecase.GameUseCase
	DailyRewardUseCase usecase.DailyRewardUseCase
}

func NewGameHandler(gameUc usecase.GameUseCase, dailyRewardUc usecase.DailyRewardUseCase) *GameHandler {
	return &GameHandler{
		GameUseCase:        gameUc,
		DailyRewardUseCase: dailyRewardUc,
	}
}

func (h *GameHandler) SyncBalance(c *fiber.Ctx) error {
	var req dto.SyncBalanceRequest
	if err := c.BodyParser(&req); err != nil {
		return response.FailedResponse(c, fiber.StatusBadRequest, err)
	}
	ctx := c.Context()

	userID, _ := helper.GetUserIDFromContext(ctx)
	if userID <= 0 {
		return response.FailedResponse(c, fiber.StatusUnauthorized, apperror.ErrUnauthorized)
	}

	res, err := h.GameUseCase.UpdateUserBalance(ctx, req.CoinsEarned, userID)
	if err != nil {
		return response.FailedResponse(c, fiber.StatusInternalServerError, err)
	}

	return response.SuccessResponse(c, fiber.StatusOK, "Sync Balance Successfully", res, nil)
}

func (h *GameHandler) GetDailyRewardStatus(c *fiber.Ctx) error {
	userID := helper.GetUserID(c)
	ctx := context.WithValue(c.Context(), helper.ContextUserIDKey, userID)

	rewards, dailyRewardIdx, isNewDay, err := h.DailyRewardUseCase.GetDailyRewardStatus(ctx)
	if err != nil {
		if errors.Is(err, apperror.ErrRecordNotFound) {
			return response.FailedResponse(c, fiber.StatusNotFound, err)
		}

		return response.FailedResponse(c, fiber.StatusInternalServerError, err)
	}

	return response.SuccessResponse(c, fiber.StatusOK, "Daily Reward Status Successfully Retrieved", dto.ToDailyRewardStatus(rewards, dailyRewardIdx, isNewDay), nil)
}

func (h *GameHandler) ClaimReward(c *fiber.Ctx) error {
	userID := helper.GetUserID(c)
	ctx := context.WithValue(c.Context(), helper.ContextUserIDKey, userID)

	reward, newBalance, err := h.DailyRewardUseCase.ClaimDailyReward(ctx)
	if err != nil {
		return response.FailedResponse(c, fiber.StatusInternalServerError, err)
	}

	return response.SuccessResponse(c, fiber.StatusOK, "Daily Reward Claimed Successfully", dto.ToClaimDailyRewardResponse(reward, newBalance), nil)
}

func (h *GameHandler) GetCurrentStage(c *fiber.Ctx) error {
	userID := helper.GetUserID(c)
	ctx := context.WithValue(c.Context(), helper.ContextUserIDKey, userID)

	stages, err := h.GameUseCase.GetGameStages(ctx, userID)
	if err != nil {
		return response.FailedResponse(c, fiber.StatusInternalServerError, err)
	}

	return response.SuccessResponse(c, fiber.StatusOK, "Current Stage Successfully Retrieved", dto.ToUserGameStageResponses(stages), nil)
}

func (h *GameHandler) Route(open fiber.Router, userAuth fiber.Router, internalAuth fiber.Router) error {
	game := userAuth.Group("/game")

	// Player game interactions
	game.Post("/sync-balance", h.SyncBalance)

	// Player game stages
	game.Get("/stages", h.GetCurrentStage)

	// Player daily reward interactions
	game.Get("/daily-reward/status", h.GetDailyRewardStatus)
	game.Post("/daily-reward/claim", h.ClaimReward)

	return nil
}
