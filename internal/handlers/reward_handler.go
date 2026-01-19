package handlers

import (
	"errors"
	"github.com/gofiber/fiber/v2"
	"github.com/winartodev/cat-cafe/internal/dto"
	"github.com/winartodev/cat-cafe/internal/middleware"
	"github.com/winartodev/cat-cafe/internal/usecase"
	"github.com/winartodev/cat-cafe/pkg/apperror"
	"github.com/winartodev/cat-cafe/pkg/helper"
	"github.com/winartodev/cat-cafe/pkg/response"
)

type RewardHandler struct {
	RewardUseCase      usecase.RewardUseCase
	DailyRewardUseCase usecase.DailyRewardUseCase
}

func NewRewardHandler(rewardUseCase usecase.RewardUseCase, dailyRewardUseCase usecase.DailyRewardUseCase) *RewardHandler {
	return &RewardHandler{
		RewardUseCase:      rewardUseCase,
		DailyRewardUseCase: dailyRewardUseCase,
	}
}

func (h *RewardHandler) CreateReward(c *fiber.Ctx) error {
	var request dto.CreateRewardRequest
	if err := c.BodyParser(&request); err != nil {
		return response.FailedResponse(c, fiber.StatusBadRequest, err)
	}

	res, err := h.RewardUseCase.CreateReward(c.Context(), request.ToEntity())
	if err != nil {
		if errors.Is(err, apperror.ErrNoUpdateRecord) {
			return response.FailedResponse(c, fiber.StatusNotFound, err)
		}

		return response.FailedResponse(c, fiber.StatusInternalServerError, err)
	}

	return response.SuccessResponse(c, fiber.StatusOK, "Reward Successfully Created", dto.ToRewardResponse(res), nil)
}

func (h *RewardHandler) GetRewards(c *fiber.Ctx) error {
	params := helper.GetPaginationParams(c)

	res, totalRows, err := h.RewardUseCase.GetRewards(c.Context(), params.Limit, params.Offset)
	if err != nil {
		if errors.Is(err, apperror.ErrNoUpdateRecord) {
			return response.FailedResponse(c, fiber.StatusNotFound, err)
		}

		return response.FailedResponse(c, fiber.StatusInternalServerError, err)
	}

	data := dto.ToRewardsResponse(res)
	meta := helper.CreatePaginationMeta(params.Page, params.Limit, totalRows)

	return response.SuccessResponse(c, fiber.StatusOK, "Reward Successfully Retrieved", data, meta)
}

func (h *RewardHandler) CreateRewardType(c *fiber.Ctx) error {
	var request dto.CreateRewardTypeRequest
	if err := c.BodyParser(&request); err != nil {
		return response.FailedResponse(c, fiber.StatusBadRequest, err)
	}

	res, err := h.RewardUseCase.CreateRewardType(c.Context(), *request.ToEntity())
	if err != nil {
		if errors.Is(err, apperror.ErrConflict) {
			return response.FailedResponse(c, fiber.StatusConflict, err)
		}

		return response.FailedResponse(c, fiber.StatusInternalServerError, err)
	}

	return response.SuccessResponse(c, fiber.StatusCreated, "Reward Type Successfully Created", dto.ToRewardTypeResponse(res), nil)
}

func (h *RewardHandler) UpdateRewardType(c *fiber.Ctx) error {
	id, err := helper.GetParam[int64](c, "id")
	if err != nil {
		return response.FailedResponse(c, fiber.StatusBadRequest, err)
	}

	var request dto.UpdateRewardTypeRequest
	if err := c.BodyParser(&request); err != nil {
		return response.FailedResponse(c, fiber.StatusBadRequest, err)
	}

	res, err := h.RewardUseCase.UpdateRewardTypes(c.Context(), id, *request.ToEntity())
	if err != nil {
		if errors.Is(err, apperror.ErrConflict) {
			return response.FailedResponse(c, fiber.StatusConflict, err)
		}

		return response.FailedResponse(c, fiber.StatusInternalServerError, err)
	}

	return response.SuccessResponse(c, fiber.StatusCreated, "Reward Type Successfully Created", dto.ToRewardTypeResponse(res), nil)
}

func (h *RewardHandler) GetRewardTypes(c *fiber.Ctx) error {
	res, err := h.RewardUseCase.GetRewardTypes(c.Context())
	if err != nil {
		return response.FailedResponse(c, fiber.StatusInternalServerError, err)
	}

	return response.SuccessResponse(c, fiber.StatusOK, "Reward Type Successfully Retrieved", dto.ToRewardTypeResponses(res), nil)
}

func (h *RewardHandler) GetRewardTypeByID(c *fiber.Ctx) error {
	id, err := helper.GetParam[int64](c, "id")
	if err != nil {
		return response.FailedResponse(c, fiber.StatusBadRequest, apperror.ErrInvalidIDParam)
	}

	res, err := h.RewardUseCase.GetRewardTypeByID(c.Context(), id)
	if err != nil {
		if errors.Is(err, apperror.ErrRecordNotFound) {
			return response.FailedResponse(c, fiber.StatusNotFound, err)
		}

		return response.FailedResponse(c, fiber.StatusInternalServerError, err)
	}

	return response.SuccessResponse(c, fiber.StatusOK, "Reward Type Successfully Retrieved", dto.ToRewardTypeResponse(res), nil)
}

func (h *RewardHandler) CreateDailyReward(c *fiber.Ctx) error {
	var request dto.DailyRewardRequest
	if err := c.BodyParser(&request); err != nil {
		return response.FailedResponse(c, fiber.StatusBadRequest, err)
	}

	ctx := c.Context()
	res, err := h.DailyRewardUseCase.CreateDailyReward(ctx, *request.ToEntity(), request.Reward)
	if err != nil {
		return response.FailedResponse(c, fiber.StatusInternalServerError, err)
	}

	return response.SuccessResponse(c, fiber.StatusCreated, "Daily Reward Successfully Created", dto.ToDailyRewardResponse(res), nil)
}

func (h *RewardHandler) GetDailyRewards(c *fiber.Ctx) error {
	ctx := c.Context()
	params := helper.GetPaginationParams(c)

	res, totalRow, err := h.DailyRewardUseCase.GetDailyRewards(ctx, params.Limit, params.Offset)
	if err != nil {
		return response.FailedResponse(c, fiber.StatusInternalServerError, err)
	}

	data := dto.ToDailyRewardResponses(res)
	meta := helper.CreatePaginationMeta(params.Page, params.Limit, totalRow)

	return response.SuccessResponse(c, fiber.StatusOK, "Daily Reward Successfully Retrieved", data, meta)
}

func (h *RewardHandler) GetDailyRewardByID(c *fiber.Ctx) error {
	id, err := helper.GetParam[int64](c, "id")
	if err != nil {
		return response.FailedResponse(c, fiber.StatusBadRequest, err)
	}

	ctx := c.Context()

	res, err := h.DailyRewardUseCase.GetDailyRewardID(ctx, id)
	if err != nil {
		if errors.Is(err, apperror.ErrRecordNotFound) {
			return response.FailedResponse(c, fiber.StatusNotFound, err)
		}

		return response.FailedResponse(c, fiber.StatusInternalServerError, err)
	}

	return response.SuccessResponse(c, fiber.StatusOK, "Daily Reward Successfully Retrieved", dto.ToDailyRewardResponse(res), nil)
}

func (h *RewardHandler) UpdateDailyReward(c *fiber.Ctx) error {
	id, err := helper.GetParam[int64](c, "id")
	if err != nil {
		return response.FailedResponse(c, fiber.StatusBadRequest, err)
	}

	var request dto.DailyRewardRequest
	if err := c.BodyParser(&request); err != nil {
		return response.FailedResponse(c, fiber.StatusBadRequest, err)
	}

	ctx := c.Context()
	res, err := h.DailyRewardUseCase.UpdateDailyReward(ctx, id, *request.ToEntity(), request.Reward)
	if err != nil {
		if errors.Is(err, apperror.ErrRecordNotFound) {
			return response.FailedResponse(c, fiber.StatusNotFound, err)
		}

		return response.FailedResponse(c, fiber.StatusInternalServerError, err)
	}

	return response.SuccessResponse(c, fiber.StatusOK, "Daily Reward Successfully Updated", dto.ToDailyRewardResponse(res), nil)
}

func (h *RewardHandler) UpdateReward(c *fiber.Ctx) error {
	// TODO: UPDATE REWARD HERE !!!
	return response.SuccessResponse(c, fiber.StatusOK, "Reward Successfully Updated", nil, nil)
}

func (h *RewardHandler) GetRewardByID(c *fiber.Ctx) error {
	// TODO: GET REWARD HERE !!!
	return response.SuccessResponse(c, fiber.StatusOK, "Reward Successfully Updated", nil, nil)
}

func (h *RewardHandler) ToggleStatusReward(c *fiber.Ctx) error {
	return response.SuccessResponse(c, fiber.StatusOK, "Reward Successfully x", nil, nil)
}

func (h *RewardHandler) ToggleStatusDailyReward(c *fiber.Ctx) error {
	return response.SuccessResponse(c, fiber.StatusOK, "Reward Successfully y", nil, nil)
}

func (h *RewardHandler) Route(r fiber.Router, m middleware.Middleware) error {
	reward := r.Group("/rewards")

	// Reward Types Management
	reward.Post("/types", h.CreateRewardType)
	reward.Get("/types", h.GetRewardTypes)
	reward.Get("/types/:id", h.GetRewardTypeByID)
	reward.Put("/types/:id", h.UpdateRewardType)

	// Daily Rewards Management
	reward.Post("/daily", h.CreateDailyReward)
	reward.Get("/daily", h.GetDailyRewards)
	reward.Get("/daily/:id", h.GetDailyRewardByID)
	reward.Put("/daily/:id", h.UpdateDailyReward)
	reward.Patch("/daily/:id/status", h.ToggleStatusDailyReward)

	// General Rewards Management
	reward.Post("/", h.CreateReward)
	reward.Get("/", h.GetRewards)
	reward.Get("/:id", h.GetRewardByID)
	reward.Put("/:id", h.UpdateReward)
	reward.Patch("/:id/status", h.ToggleStatusReward)

	return nil
}
