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

	userID := helper.GetUserID(c)
	ctx := context.WithValue(c.Context(), helper.ContextUserIDKey, userID)

	res, err := h.GameUseCase.UpdateUserBalance(ctx, req.CoinsEarned)
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

func (h *GameHandler) GetAllStages(c *fiber.Ctx) error {
	userID := helper.GetUserID(c)
	ctx := context.WithValue(c.Context(), helper.ContextUserIDKey, userID)

	stages, _, err := h.GameUseCase.GetGameStages(ctx)
	if err != nil {
		return response.FailedResponse(c, h.errorHandler, err)
	}

	return response.SuccessResponse(c, fiber.StatusOK, "Current Stage Successfully Retrieved", dto.ToUserGameStageResponses(stages), nil)
}

func (h *GameHandler) GetCurrentStage(c *fiber.Ctx) error {
	userID := helper.GetUserID(c)
	ctx := context.WithValue(c.Context(), helper.ContextUserIDKey, userID)

	gameStage, config, nextStage, err := h.GameUseCase.GetCurrentGameStage(ctx)
	if err != nil {
		return response.FailedResponse(c, h.errorHandler, err)
	}

	return response.SuccessResponse(c, fiber.StatusOK, "Current Stage Successfully Retrieved", dto.ToUserDetailGameStageResponse(gameStage, config, nextStage), nil)
}

func (h *GameHandler) StartGameStage(c *fiber.Ctx) error {
	slug, err := helper.GetParam[string](c, "slug")
	if err != nil {
		return response.FailedResponse(c, h.errorHandler, apperror.ErrInvalidParam)
	}

	userID := helper.GetUserID(c)
	ctx := context.WithValue(c.Context(), helper.ContextUserIDKey, userID)

	gameStage, config, nextStage, err := h.GameUseCase.StartGameStage(ctx, slug)
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
	ctx := context.WithValue(c.Context(), helper.ContextUserIDKey, userID)

	err = h.GameUseCase.CompleteGameStage(ctx, slug)
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
	ctx := context.WithValue(c.Context(), helper.ContextUserIDKey, userID)

	res, err := h.GameUseCase.UpgradeKitchenStation(ctx, slug)
	if err != nil {
		return response.FailedResponse(c, h.errorHandler, err)
	}

	return response.SuccessResponse(c, fiber.StatusOK, "Kitchen Station Successfully Upgraded", dto.ToUserUpgradeKitchenResponse(res), nil)
}

func (h *GameHandler) PurchaseKitchenStation(c *fiber.Ctx) error {
	slug, err := helper.GetParam[string](c, "slug")
	if err != nil {
		return response.FailedResponse(c, h.errorHandler, apperror.ErrInvalidParam)
	}

	userID := helper.GetUserID(c)
	ctx := context.WithValue(c.Context(), helper.ContextUserIDKey, userID)

	res, err := h.GameUseCase.UnlockKitchenStation(ctx, slug)
	if err != nil {
		return response.FailedResponse(c, h.errorHandler, err)
	}

	return response.SuccessResponse(c, fiber.StatusOK, "Kitchen Station Successfully Purchased", dto.ToUserUnlockKitchenResponse(res), nil)
}

func (h *GameHandler) GetStageUpgrades(c *fiber.Ctx) error {
	userID := helper.GetUserID(c)
	ctx := context.WithValue(c.Context(), helper.ContextUserIDKey, userID)

	res, err := h.GameUseCase.GetStageUpgrades(ctx)
	if err != nil {
		return response.FailedResponse(c, h.errorHandler, err)
	}

	return response.SuccessResponse(c, fiber.StatusOK, "Upgrades Successfully Retrieved", dto.ToUserStageUpgradesResponse(res), nil)
}

func (h *GameHandler) PurchaseStageUpgrade(c *fiber.Ctx) error {
	slug, err := helper.GetParam[string](c, "slug")
	if err != nil {
		return response.FailedResponse(c, h.errorHandler, apperror.ErrInvalidParam)
	}

	userID := helper.GetUserID(c)
	ctx := context.WithValue(c.Context(), helper.ContextUserIDKey, userID)
	res, err := h.GameUseCase.PurchaseStageUpgrade(ctx, slug)
	if err != nil {
		return response.FailedResponse(c, h.errorHandler, err)
	}

	return response.SuccessResponse(c, fiber.StatusOK, "Upgrade Successfully Purchased", dto.ToUserPurchasedStageUpgradeResponse(res), nil)
}

func (h *GameHandler) Route(open fiber.Router, userAuth fiber.Router, internalAuth fiber.Router) error {
	game := userAuth.Group("/game")

	// Player game stages
	stages := game.Group("/stages")
	stages.Get("/", h.GetAllStages)
	stages.Get("/current", h.GetCurrentStage)
	stages.Post("/:slug/start", h.StartGameStage)
	stages.Post("/:slug/complete", h.CompleteGameStage)

	// Player Kitchen Station
	stations := game.Group("/stations")
	stations.Post("/:slug/purchase", h.PurchaseKitchenStation)
	stations.Post("/:slug/upgrade", h.UpgradeKitchenStation)

	// Player Upgrade
	upgrades := game.Group("/upgrades")
	upgrades.Get("/", h.GetStageUpgrades)
	upgrades.Post("/:slug/purchase", h.PurchaseStageUpgrade)

	// Player Economy & Rewards
	game.Post("/sync-balance", h.SyncBalance)
	game.Get("/daily-reward/status", h.GetDailyRewardStatus)
	game.Post("/daily-reward/claim", h.ClaimReward)

	return nil
}
