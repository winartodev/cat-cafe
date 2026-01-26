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

type FoodItemHandler struct {
	foodItemUseCase usecase.FoodItemUseCase
	errorHandler    *apperror.ErrorHandler
}

func NewFoodItemHandler(foodItemUC usecase.FoodItemUseCase) *FoodItemHandler {
	return &FoodItemHandler{
		foodItemUseCase: foodItemUC,
		errorHandler:    apperror.NewErrorHandler(),
	}
}

func (h *FoodItemHandler) CreateFood(c *fiber.Ctx) error {
	var request dto.FoodItemRequest
	if err := c.BodyParser(&request); err != nil {
		return response.FailedResponse(c, h.errorHandler, err)
	}

	res, err := h.foodItemUseCase.CreateFood(c.Context(), request.ToEntity())
	if err != nil {
		if errors.Is(err, apperror.ErrNoUpdateRecord) {
			return response.FailedResponse(c, h.errorHandler, err)
		}

		return response.FailedResponse(c, h.errorHandler, err)
	}

	return response.SuccessResponse(c, fiber.StatusOK, "Food Successfully Created", dto.ToFoodItemResponse(res), nil)
}

func (h *FoodItemHandler) GetRewards(c *fiber.Ctx) error {
	params := helper.GetPaginationParams(c)

	res, totalRows, err := h.foodItemUseCase.GetFoods(c.Context(), params.Limit, params.Offset)
	if err != nil {
		if errors.Is(err, apperror.ErrNoUpdateRecord) {
			return response.FailedResponse(c, h.errorHandler, err)
		}

		return response.FailedResponse(c, h.errorHandler, err)
	}

	data := dto.ToFoodItemsResponse(res)
	meta := helper.CreatePaginationMeta(params.Page, params.Limit, totalRows)

	return response.SuccessResponse(c, fiber.StatusOK, "Food Successfully Retrieved", data, meta)
}

func (h *FoodItemHandler) GetReward(c *fiber.Ctx) error {
	id, err := helper.GetParam[int64](c, "id")
	if err != nil {
		return response.FailedResponse(c, h.errorHandler, apperror.ErrInvalidParam)
	}

	res, err := h.foodItemUseCase.GetFoodByID(c.Context(), id)
	if err != nil {
		if errors.Is(err, apperror.ErrNoUpdateRecord) {
			return response.FailedResponse(c, h.errorHandler, err)
		}

		return response.FailedResponse(c, h.errorHandler, err)
	}

	data := dto.ToFoodItemResponse(res)

	return response.SuccessResponse(c, fiber.StatusOK, "Food Successfully Retrieved", data, nil)
}

func (h *FoodItemHandler) UpdateReward(c *fiber.Ctx) error {
	id, err := helper.GetParam[int64](c, "id")
	if err != nil {
		return response.FailedResponse(c, h.errorHandler, apperror.ErrInvalidParam)
	}

	var request dto.FoodItemRequest
	if err := c.BodyParser(&request); err != nil {
		return response.FailedResponse(c, h.errorHandler, err)
	}

	res, err := h.foodItemUseCase.UpdateFood(c.Context(), id, request.ToEntity())
	if err != nil {
		if errors.Is(err, apperror.ErrNoUpdateRecord) {
			return response.FailedResponse(c, h.errorHandler, err)
		}

		return response.FailedResponse(c, h.errorHandler, err)
	}

	data := dto.ToFoodItemResponse(res)

	return response.SuccessResponse(c, fiber.StatusOK, "Food Successfully Updated", data, nil)
}

func (h *FoodItemHandler) Route(open fiber.Router, userAuth fiber.Router, internalAuth fiber.Router) error {
	foodItem := internalAuth.Group("/foods")

	foodItem.Post("/", h.CreateFood)
	foodItem.Get("/", h.GetRewards)
	foodItem.Get("/:id", h.GetReward)
	foodItem.Put("/:id", h.UpdateReward)

	return nil
}
