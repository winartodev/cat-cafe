package handlers

import (
	"errors"

	"github.com/gofiber/fiber/v2"
	"github.com/winartodev/cat-cafe/internal/dto"
	"github.com/winartodev/cat-cafe/internal/usecase"
	"github.com/winartodev/cat-cafe/pkg/apperror"
	"github.com/winartodev/cat-cafe/pkg/helper"
	"github.com/winartodev/cat-cafe/pkg/response"
)

type DailyRewardHandler struct {
	DailyRewardUseCase usecase.DailyRewardUseCase
}

func NewDailyRewardHandler(dailyRewardUseCase usecase.DailyRewardUseCase) *DailyRewardHandler {
	return &DailyRewardHandler{
		DailyRewardUseCase: dailyRewardUseCase,
	}
}

func (d *DailyRewardHandler) CreateRewardType(c *fiber.Ctx) error {
	var request dto.CreateRewardTypeRequest
	if err := c.BodyParser(&request); err != nil {
		return response.FailedResponse(c, fiber.StatusBadRequest, err)
	}

	res, err := d.DailyRewardUseCase.CreateRewardType(c.Context(), *request.ToEntity())
	if err != nil {
		if errors.Is(err, apperror.ErrConflict) {
			return response.FailedResponse(c, fiber.StatusConflict, err)
		}

		return response.FailedResponse(c, fiber.StatusInternalServerError, err)
	}

	return response.SuccessResponse(c, fiber.StatusCreated, "Reward Type Successfully Created", dto.ToRewardTypeResponse(res), nil)
}

func (d *DailyRewardHandler) GetRewardTypes(c *fiber.Ctx) error {
	res, err := d.DailyRewardUseCase.GetRewardTypes(c.Context())
	if err != nil {
		return response.FailedResponse(c, fiber.StatusInternalServerError, err)
	}

	return response.SuccessResponse(c, fiber.StatusOK, "Reward Type Successfully Retrieved", dto.ToRewardTypeResponses(res), nil)
}

func (d *DailyRewardHandler) GetRewardTypeByID(c *fiber.Ctx) error {
	id, err := helper.GetParam[int64](c, "id")
	if err != nil {
		return response.FailedResponse(c, fiber.StatusBadRequest, apperror.ErrInvalidIDParam)
	}

	res, err := d.DailyRewardUseCase.GetRewardTypeByID(c.Context(), id)
	if err != nil {
		if errors.Is(err, apperror.ErrRecordNotFound) {
			return response.FailedResponse(c, fiber.StatusNotFound, err)
		}

		return response.FailedResponse(c, fiber.StatusInternalServerError, err)
	}

	return response.SuccessResponse(c, fiber.StatusOK, "Reward Type Successfully Retrieved", dto.ToRewardTypeResponse(res), nil)
}

func (d *DailyRewardHandler) UpdateRewardTypes(c *fiber.Ctx) error {
	id, err := helper.GetParam[int64](c, "id")
	if err != nil {
		return response.FailedResponse(c, fiber.StatusBadRequest, apperror.ErrInvalidIDParam)
	}

	var request dto.UpdateRewardTypeRequest
	if err := c.BodyParser(&request); err != nil {
		return response.FailedResponse(c, fiber.StatusBadRequest, err)
	}

	res, err := d.DailyRewardUseCase.UpdateRewardTypes(c.Context(), id, *request.ToEntity())
	if err != nil {
		if errors.Is(err, apperror.ErrNoUpdateRecord) {
			return response.FailedResponse(c, fiber.StatusNotFound, err)
		}

		return response.FailedResponse(c, fiber.StatusInternalServerError, err)
	}

	return response.SuccessResponse(c, fiber.StatusOK, "Reward Type Successfully Updated", dto.ToRewardTypeResponse(res), nil)
}

func (d *DailyRewardHandler) CreateDailyReward(c *fiber.Ctx) error {
	var request dto.DailyRewardRequest
	if err := c.BodyParser(&request); err != nil {
		return response.FailedResponse(c, fiber.StatusBadRequest, err)
	}

	ctx := c.Context()
	res, err := d.DailyRewardUseCase.CreateDailyReward(ctx, *request.ToEntity(), request.RewardType)
	if err != nil {
		return response.FailedResponse(c, fiber.StatusInternalServerError, err)
	}

	return response.SuccessResponse(c, fiber.StatusCreated, "Daily Reward Successfully Created", dto.ToDailyRewardResponse(res), nil)
}

func (d *DailyRewardHandler) GetDailyRewards(c *fiber.Ctx) error {
	ctx := c.Context()

	res, err := d.DailyRewardUseCase.GetDailyRewards(ctx)
	if err != nil {
		return response.FailedResponse(c, fiber.StatusInternalServerError, err)
	}

	return response.SuccessResponse(c, fiber.StatusOK, "Daily Reward Successfully Retrieved", dto.ToDailyRewardResponses(res), nil)
}

func (d *DailyRewardHandler) GetDailyRewardByID(c *fiber.Ctx) error {
	id, err := helper.GetParam[int64](c, "id")
	if err != nil {
		return response.FailedResponse(c, fiber.StatusBadRequest, err)
	}

	ctx := c.Context()

	res, err := d.DailyRewardUseCase.GetDailyRewardID(ctx, id)
	if err != nil {
		if errors.Is(err, apperror.ErrRecordNotFound) {
			return response.FailedResponse(c, fiber.StatusNotFound, err)
		}

		return response.FailedResponse(c, fiber.StatusInternalServerError, err)
	}

	return response.SuccessResponse(c, fiber.StatusOK, "Daily Reward Successfully Retrieved", dto.ToDailyRewardResponse(res), nil)
}

func (d *DailyRewardHandler) UpdateDailyReward(c *fiber.Ctx) error {
	id, err := helper.GetParam[int64](c, "id")
	if err != nil {
		return response.FailedResponse(c, fiber.StatusBadRequest, err)
	}

	var request dto.DailyRewardRequest
	if err := c.BodyParser(&request); err != nil {
		return response.FailedResponse(c, fiber.StatusBadRequest, err)
	}

	ctx := c.Context()
	res, err := d.DailyRewardUseCase.UpdateDailyReward(ctx, id, *request.ToEntity(), request.RewardType)
	if err != nil {
		if errors.Is(err, apperror.ErrRecordNotFound) {
			return response.FailedResponse(c, fiber.StatusNotFound, err)
		}

		return response.FailedResponse(c, fiber.StatusInternalServerError, err)
	}

	return response.SuccessResponse(c, fiber.StatusOK, "Daily Reward Successfully Updated", dto.ToDailyRewardResponse(res), nil)
}

func (d *DailyRewardHandler) GetDailyRewardStatus(c *fiber.Ctx) error {
	ctx := c.Context()
	rewards, dailyRewardIdx, isNewDay, err := d.DailyRewardUseCase.GetRewardStatus(ctx)
	if err != nil {
		if errors.Is(err, apperror.ErrRecordNotFound) {
			return response.FailedResponse(c, fiber.StatusNotFound, err)
		}

		return response.FailedResponse(c, fiber.StatusInternalServerError, err)
	}

	return response.SuccessResponse(c, fiber.StatusOK, "Daily Reward Status Successfully Retrieved", dto.ToDailyRewardStatus(rewards, dailyRewardIdx, isNewDay), nil)
}

func (d *DailyRewardHandler) ClaimReward(c *fiber.Ctx) error {
	ctx := c.Context()
	reward, newBalance, err := d.DailyRewardUseCase.ClaimReward(ctx)
	if err != nil {
		return response.FailedResponse(c, fiber.StatusInternalServerError, err)
	}

	return response.SuccessResponse(c, fiber.StatusOK, "Daily Reward Claimed Successfully", dto.ToClaimDailyRewardResponse(reward, newBalance), nil)
}

func (d *DailyRewardHandler) Route(route fiber.Router) error {
	reward := route.Group("/rewards")

	reward.Post("/types", d.CreateRewardType)
	reward.Get("/types", d.GetRewardTypes)

	reward.Get("/types/:id", d.GetRewardTypeByID)
	reward.Put("/types/:id", d.UpdateRewardTypes)

	reward.Post("/daily", d.CreateDailyReward)
	reward.Get("/daily", d.GetDailyRewards)
	reward.Get("/daily/status", d.GetDailyRewardStatus)
	reward.Post("/daily/claim", d.ClaimReward)

	reward.Get("/daily/:id", d.GetDailyRewardByID)
	reward.Put("/daily/:id", d.UpdateDailyReward)

	return nil
}
