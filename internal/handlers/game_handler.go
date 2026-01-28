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
	errorHandler       *apperror.ErrorHandler
}

func NewGameHandler(gameUc usecase.GameUseCase, dailyRewardUc usecase.DailyRewardUseCase) *GameHandler {
	return &GameHandler{
		GameUseCase:        gameUc,
		DailyRewardUseCase: dailyRewardUc,
		errorHandler:       apperror.NewErrorHandler(),
	}
}

func (h *GameHandler) SyncBalance(c *fiber.Ctx) error {
	var req dto.SyncBalanceRequest
	if err := c.BodyParser(&req); err != nil {
		return response.FailedResponse(c, h.errorHandler, err)
	}
	ctx := c.Context()

	userID, _ := helper.GetUserIDFromContext(ctx)
	if userID <= 0 {
		return response.FailedResponse(c, h.errorHandler, apperror.ErrUnauthorized)
	}

	res, err := h.GameUseCase.UpdateUserBalance(ctx, req.CoinsEarned, userID)
	if err != nil {
		return response.FailedResponse(c, h.errorHandler, err)
	}

	return response.SuccessResponse(c, fiber.StatusOK, "Sync Balance Successfully", res, nil)
}

func (h *GameHandler) GetDailyRewardStatus(c *fiber.Ctx) error {
	userID := helper.GetUserID(c)
	ctx := context.WithValue(c.Context(), helper.ContextUserIDKey, userID)

	rewards, dailyRewardIdx, isNewDay, err := h.DailyRewardUseCase.GetDailyRewardStatus(ctx)
	if err != nil {
		if errors.Is(err, apperror.ErrRecordNotFound) {
			return response.FailedResponse(c, h.errorHandler, err)
		}

		return response.FailedResponse(c, h.errorHandler, err)
	}

	return response.SuccessResponse(c, fiber.StatusOK, "Daily Reward Status Successfully Retrieved", dto.ToDailyRewardStatus(rewards, dailyRewardIdx, isNewDay), nil)
}

func (h *GameHandler) ClaimReward(c *fiber.Ctx) error {
	userID := helper.GetUserID(c)
	ctx := context.WithValue(c.Context(), helper.ContextUserIDKey, userID)

	reward, newBalance, err := h.DailyRewardUseCase.ClaimDailyReward(ctx)
	if err != nil {
		return response.FailedResponse(c, h.errorHandler, err)
	}

	return response.SuccessResponse(c, fiber.StatusOK, "Daily Reward Claimed Successfully", dto.ToClaimDailyRewardResponse(reward, newBalance), nil)
}

func (h *GameHandler) GetCurrentStage(c *fiber.Ctx) error {
	userID := helper.GetUserID(c)
	ctx := context.WithValue(c.Context(), helper.ContextUserIDKey, userID)

	stages, _, err := h.GameUseCase.GetGameStages(ctx, userID)
	if err != nil {
		return response.FailedResponse(c, h.errorHandler, err)
	}

	return response.SuccessResponse(c, fiber.StatusOK, "Current Stage Successfully Retrieved", dto.ToUserGameStageResponses(stages), nil)
}

func (h *GameHandler) StartGameStage(c *fiber.Ctx) error {
	slug, err := helper.GetParam[string](c, "slug")
	if err != nil {
		return response.FailedResponse(c, h.errorHandler, apperror.ErrInvalidParam)
	}

	userID := helper.GetUserID(c)

	gameStage, config, nextStage, err := h.GameUseCase.StartGameStage(c.Context(), userID, slug)
	if err != nil {
		return response.FailedResponse(c, h.errorHandler, err)
	}

	return response.SuccessResponse(c, fiber.StatusOK, "Game Stage Successfully Started", dto.ToUserDetailGameStageResponse(gameStage, config, nextStage), nil)
}

func (h *GameHandler) CompleteGameStage(c *fiber.Ctx) error {
	slug, err := helper.GetParam[string](c, "slug")
	if err != nil {
		return response.FailedResponse(c, h.errorHandler, apperror.ErrInvalidParam)
	}

	userID := helper.GetUserID(c)

	err = h.GameUseCase.CompleteGameStage(c.Context(), userID, slug)
	if err != nil {
		return response.FailedResponse(c, h.errorHandler, err)
	}

	return response.SuccessResponse(c, fiber.StatusOK, "Game Stage Successfully Started", nil, nil)
}

func (h *GameHandler) UpgradeKitchenStation(c *fiber.Ctx) error {
	slug, err := helper.GetParam[string](c, "slug")
	if err != nil {
		return response.FailedResponse(c, h.errorHandler, apperror.ErrInvalidParam)
	}

	userID := helper.GetUserID(c)

	res, err := h.GameUseCase.UpgradeKitchenStation(c.Context(), userID, slug)
	if err != nil {
		return response.FailedResponse(c, h.errorHandler, err)
	}

	return response.SuccessResponse(c, fiber.StatusOK, "Kitchen Station Successfully Upgraded", res, nil)
}

func (h *GameHandler) PurchaseKitchenStation(c *fiber.Ctx) error {
	slug, err := helper.GetParam[string](c, "slug")
	if err != nil {
		return response.FailedResponse(c, h.errorHandler, apperror.ErrInvalidParam)
	}

	userID := helper.GetUserID(c)

	res, err := h.GameUseCase.UnlockKitchenStation(c.Context(), userID, slug)
	if err != nil {
		return response.FailedResponse(c, h.errorHandler, err)
	}

	return response.SuccessResponse(c, fiber.StatusOK, "Kitchen Station Successfully Purchased", res, nil)
}

func (h *GameHandler) Route(open fiber.Router, userAuth fiber.Router, internalAuth fiber.Router) error {
	game := userAuth.Group("/game")

	// Player game interactions
	game.Post("/sync-balance", h.SyncBalance)

	// Player game stages
	game.Get("/stages", h.GetCurrentStage)
	game.Post("/stages/:slug/start", h.StartGameStage)
	game.Post("/stages/:slug/complete", h.CompleteGameStage)

	// Player Kitchen Station
	game.Post("/stations/:slug/purchase", h.PurchaseKitchenStation)
	game.Post("/stations/:slug/upgrade", h.UpgradeKitchenStation)

	// Player daily reward interactions
	game.Get("/daily-reward/status", h.GetDailyRewardStatus)
	game.Post("/daily-reward/claim", h.ClaimReward)

	return nil
}
