package handlers

import (
	"github.com/gofiber/fiber/v2"
	"github.com/winartodev/cat-cafe/internal/dto"
	"github.com/winartodev/cat-cafe/internal/usecase"
	"github.com/winartodev/cat-cafe/pkg/helper"
	"github.com/winartodev/cat-cafe/pkg/response"
)

type GameStageHandler struct {
	GameStageUseCase usecase.GameStageUseCase
}

func NewGameStageHandler(gameStageUseCase usecase.GameStageUseCase) *GameStageHandler {
	return &GameStageHandler{
		GameStageUseCase: gameStageUseCase,
	}
}

func (h *GameStageHandler) CreateGameStage(c *fiber.Ctx) error {
	var req dto.CreateGameStageRequest
	if err := c.BodyParser(&req); err != nil {
		return response.FailedResponse(c, fiber.StatusBadRequest, err)
	}

	ctx := c.Context()
	gameStage, stageConfig, err := req.ToEntities()
	if err != nil {
		return response.FailedResponse(c, fiber.StatusBadRequest, err)
	}

	data, err := h.GameStageUseCase.CreateGameStage(ctx, gameStage, stageConfig)
	if err != nil {
		return response.FailedResponse(c, fiber.StatusInternalServerError, err)
	}

	helper.PrettyPrint(gameStage)
	helper.PrettyPrint(stageConfig)

	return response.SuccessResponse(c, fiber.StatusCreated, "Game Stage Successfully Created", data, nil)
}

func (h *GameStageHandler) UpdateGameStage(c *fiber.Ctx) error {
	id, err := helper.GetParam[int64](c, "id")
	if err != nil {
		return response.FailedResponse(c, fiber.StatusBadRequest, err)
	}

	var req dto.UpdateGameStageRequest
	if err := c.BodyParser(&req); err != nil {
		return response.FailedResponse(c, fiber.StatusBadRequest, err)
	}

	ctx := c.Context()
	gameStage, stageConfig, err := req.ToEntities(id)
	if err != nil {
		return response.FailedResponse(c, fiber.StatusBadRequest, err)
	}

	data, err := h.GameStageUseCase.UpdateGameStage(ctx, gameStage, stageConfig)
	if err != nil {
		return response.FailedResponse(c, fiber.StatusInternalServerError, err)
	}

	return response.SuccessResponse(c, fiber.StatusOK, "Game Stage Successfully Updated", data, nil)
}

func (h *GameStageHandler) GetGameStages(c *fiber.Ctx) error {
	params := helper.GetPaginationParams(c)

	ctx := c.Context()

	res, totalRows, err := h.GameStageUseCase.GetGameStages(ctx, params.Limit, params.Offset)
	if err != nil {
		return response.FailedResponse(c, fiber.StatusInternalServerError, err)
	}

	data := dto.ToGameStageResponses(res)
	meta := helper.CreatePaginationMeta(params.Page, params.Limit, totalRows)

	return response.SuccessResponse(c, fiber.StatusOK, "Game Stage Successfully Retrieved", data, meta)
}

func (h *GameStageHandler) GetGameStage(c *fiber.Ctx) error {
	id, err := helper.GetParam[int64](c, "id")
	if err != nil {
		return response.FailedResponse(c, fiber.StatusBadRequest, err)
	}

	ctx := c.Context()

	gameStage, gameConfig, err := h.GameStageUseCase.GetGameStageByID(ctx, id)
	if err != nil {
		return response.FailedResponse(c, fiber.StatusInternalServerError, err)
	}

	return response.SuccessResponse(c, fiber.StatusOK, "Game Stage Successfully Retrieved", dto.ToGameStageDetailResponse(gameStage, gameConfig), nil)
}

func (h *GameStageHandler) Route(open fiber.Router, userAuth fiber.Router, internalAuth fiber.Router) error {
	gameStages := internalAuth.Group("/game-stages")

	gameStages.Post("/", h.CreateGameStage)
	gameStages.Put("/:id", h.UpdateGameStage)
	gameStages.Get("/", h.GetGameStages)
	gameStages.Get("/:id", h.GetGameStage)

	return nil
}
