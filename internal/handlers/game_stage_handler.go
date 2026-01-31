package handlers

import (
	"github.com/gofiber/fiber/v2"
	"github.com/winartodev/cat-cafe/internal/dto"
	"github.com/winartodev/cat-cafe/internal/usecase"
	"github.com/winartodev/cat-cafe/pkg/apperror"
	"github.com/winartodev/cat-cafe/pkg/helper"
	"github.com/winartodev/cat-cafe/pkg/response"
)

type GameStageHandler struct {
	GameStageUseCase usecase.GameStageUseCase
	errorHandler     *apperror.ErrorHandler
}

func NewGameStageHandler(gameStageUseCase usecase.GameStageUseCase) *GameStageHandler {
	return &GameStageHandler{
		GameStageUseCase: gameStageUseCase,
		errorHandler:     apperror.NewErrorHandler(),
	}
}

func (h *GameStageHandler) CreateGameStage(c *fiber.Ctx) error {
	var req dto.CreateGameStageRequest
	if err := c.BodyParser(&req); err != nil {
		return response.FailedResponse(c, h.errorHandler, err)
	}

	ctx := c.Context()
	gameStage, stageConfig, err := req.ToEntities()
	if err != nil {
		return response.FailedResponse(c, h.errorHandler, err)
	}

	data, err := h.GameStageUseCase.CreateGameStage(ctx, gameStage, stageConfig)
	if err != nil {
		return response.FailedResponse(c, h.errorHandler, err)
	}

	return response.SuccessResponse(c, fiber.StatusCreated, "Game Stage Successfully Created", data, nil)
}

func (h *GameStageHandler) UpdateGameStage(c *fiber.Ctx) error {
	id, err := helper.GetParam[int64](c, "id")
	if err != nil {
		return response.FailedResponse(c, h.errorHandler, err)
	}

	var req dto.UpdateGameStageRequest
	if err := c.BodyParser(&req); err != nil {
		return response.FailedResponse(c, h.errorHandler, err)
	}

	ctx := c.Context()
	gameStage, stageConfig, err := req.ToEntities(id)
	if err != nil {
		return response.FailedResponse(c, h.errorHandler, err)
	}

	data, err := h.GameStageUseCase.UpdateGameStage(ctx, gameStage, stageConfig)
	if err != nil {
		return response.FailedResponse(c, h.errorHandler, err)
	}

	return response.SuccessResponse(c, fiber.StatusOK, "Game Stage Successfully Updated", data, nil)
}

func (h *GameStageHandler) GetGameStages(c *fiber.Ctx) error {
	params := helper.GetPaginationParams(c)

	ctx := c.Context()

	res, totalRows, err := h.GameStageUseCase.GetGameStages(ctx, params.Limit, params.Offset)
	if err != nil {
		return response.FailedResponse(c, h.errorHandler, err)
	}

	data := dto.ToGameStageResponses(res)
	meta := helper.CreatePaginationMeta(params.Page, params.Limit, totalRows)

	return response.SuccessResponse(c, fiber.StatusOK, "Game Stage Successfully Retrieved", data, meta)
}

func (h *GameStageHandler) GetGameStage(c *fiber.Ctx) error {
	id, err := helper.GetParam[int64](c, "id")
	if err != nil {
		return response.FailedResponse(c, h.errorHandler, err)
	}

	ctx := c.Context()

	gameStage, gameConfig, err := h.GameStageUseCase.GetGameStageByID(ctx, id)
	if err != nil {
		return response.FailedResponse(c, h.errorHandler, err)
	}

	return response.SuccessResponse(c, fiber.StatusOK, "Game Stage Successfully Retrieved", dto.ToGameStageDetailResponse(gameStage, gameConfig), nil)
}

func (h *GameStageHandler) CreateStageUpgrade(c *fiber.Ctx) error {
	var req dto.CreateStageUpgradeRequest
	if err := c.BodyParser(&req); err != nil {
		return response.FailedResponse(c, h.errorHandler, err)
	}

	ctx := c.Context()

	err := h.GameStageUseCase.CreateStageUpgrade(ctx, req.Stage, req.Upgrades)
	if err != nil {
		return response.FailedResponse(c, h.errorHandler, err)
	}

	return response.SuccessResponse(c, fiber.StatusOK, "Stage Upgrade Successfully Created", dto.ToUpgradeStageResponse(req.Stage, req.Upgrades), nil)
}

func (h *GameStageHandler) GetGameStageUpgrades(c *fiber.Ctx) error {
	slug, err := helper.GetParam[string](c, "slug")
	if err != nil {
		return response.FailedResponse(c, h.errorHandler, err)
	}

	paginateParam := helper.GetPaginationParams(c)

	ctx := c.Context()
	stageUpgrades, total, err := h.GameStageUseCase.GetStageUpgrades(ctx, slug, paginateParam.Limit, paginateParam.Offset)
	if err != nil {
		return response.FailedResponse(c, h.errorHandler, err)
	}

	meta := helper.CreatePaginationMeta(paginateParam.Page, paginateParam.Limit, total)

	return response.SuccessResponse(c, fiber.StatusOK, "", dto.ToStageUpgradesResponseDTO(slug, stageUpgrades), meta)
}

func (h *GameStageHandler) UpdateStageUpgrade(c *fiber.Ctx) error {
	slug, err := helper.GetParam[string](c, "slug")
	if err != nil {
		return response.FailedResponse(c, h.errorHandler, err)
	}

	var req dto.UpdateStageUpgradeRequest
	if err := c.BodyParser(&req); err != nil {
		return response.FailedResponse(c, h.errorHandler, err)
	}

	ctx := c.Context()

	err = h.GameStageUseCase.UpdateStageUpgrades(ctx, slug, req.Upgrades)
	if err != nil {
		return response.FailedResponse(c, h.errorHandler, err)
	}

	return response.SuccessResponse(c, fiber.StatusOK, "Stage Upgrade Successfully Updated", dto.ToUpgradeStageResponse(slug, req.Upgrades), nil)
}

func (h *GameStageHandler) Route(open fiber.Router, userAuth fiber.Router, internalAuth fiber.Router) error {
	gameStages := internalAuth.Group("/game-stages")

	gameStages.Post("/", h.CreateGameStage)
	gameStages.Put("/:id", h.UpdateGameStage)
	gameStages.Get("/", h.GetGameStages)
	gameStages.Get("/:id", h.GetGameStage)

	stageUpgrade := internalAuth.Group("/stage-upgrades")

	stageUpgrade.Post("/", h.CreateStageUpgrade)
	stageUpgrade.Get("/:slug", h.GetGameStageUpgrades)
	stageUpgrade.Put("/:slug", h.UpdateStageUpgrade)

	return nil
}
