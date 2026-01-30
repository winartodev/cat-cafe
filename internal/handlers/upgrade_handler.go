package handlers

import (
	"fmt"

	"github.com/gofiber/fiber/v2"
	"github.com/winartodev/cat-cafe/internal/dto"
	"github.com/winartodev/cat-cafe/internal/usecase"
	"github.com/winartodev/cat-cafe/pkg/apperror"
	"github.com/winartodev/cat-cafe/pkg/helper"
	"github.com/winartodev/cat-cafe/pkg/response"
)

type UpgradeHandler struct {
	errorHandler *apperror.ErrorHandler
	upgradeUC    usecase.UpgradeUseCase
}

func NewUpgradeHandler(upgradeUC usecase.UpgradeUseCase) *UpgradeHandler {
	return &UpgradeHandler{
		errorHandler: apperror.NewErrorHandler(),
		upgradeUC:    upgradeUC,
	}
}

func (h *UpgradeHandler) CreateUpgrade(c *fiber.Ctx) error {
	var request dto.CreateUpgradeDTO
	if err := c.BodyParser(&request); err != nil {
		return response.FailedResponse(c, h.errorHandler, err)
	}

	if err := request.ValidateRequest(); err != nil {
		return response.FailedResponse(c, h.errorHandler, err)
	}

	data := request.ToEntity()
	upgrade, err := h.upgradeUC.CreateUpgrade(c.Context(), data)
	if err != nil {
		return response.FailedResponse(c, h.errorHandler, err)
	}

	return response.SuccessResponse(c, fiber.StatusOK, "Upgrade Successfully Created", dto.ToCreateUpgradeResponseDTO(upgrade), nil)
}

func (h *UpgradeHandler) GetUpgrades(c *fiber.Ctx) error {
	params := helper.GetPaginationParams(c)

	upgrades, totalRows, err := h.upgradeUC.GetUpgrades(c.Context(), params.Limit, params.Offset)
	if err != nil {
		return response.FailedResponse(c, h.errorHandler, err)
	}

	data := dto.ToGetUpgradesResponseDTO(upgrades)
	meta := helper.CreatePaginationMeta(params.Page, params.Limit, totalRows)

	return response.SuccessResponse(c, fiber.StatusOK, "Upgrades Successfully Retrieved", data, meta)
}

func (h *UpgradeHandler) GetUpgradeByID(c *fiber.Ctx) error {
	id, err := helper.GetParam[int64](c, "id")
	if err != nil {
		return response.FailedResponse(c, h.errorHandler, err)
	}

	upgrade, err := h.upgradeUC.GetUpgradeByID(c.Context(), id)
	if err != nil {
		return response.FailedResponse(c, h.errorHandler, err)
	}

	if upgrade == nil {
		return response.FailedResponse(c, h.errorHandler, apperror.ErrorNotFound("upgrade", fmt.Sprint(id)))
	}

	return response.SuccessResponse(c, fiber.StatusOK, "Upgrade Successfully Retrieved", dto.ToDetailUpgradeResponseDTO(upgrade), nil)
}

func (h *UpgradeHandler) UpdateUpgrade(c *fiber.Ctx) error {
	id, err := helper.GetParam[int64](c, "id")
	if err != nil {
		return response.FailedResponse(c, h.errorHandler, err)
	}

	var request dto.UpdateUpgradeDTO
	if err := c.BodyParser(&request); err != nil {
		return response.FailedResponse(c, h.errorHandler, err)
	}

	if err := request.ValidateRequest(); err != nil {
		return response.FailedResponse(c, h.errorHandler, err)
	}

	upgrades, err := h.upgradeUC.UpdateUpgrade(c.Context(), id, request.ToEntity())
	if err != nil {
		return response.FailedResponse(c, h.errorHandler, err)
	}

	data := dto.ToDetailUpgradeResponseDTO(upgrades)
	return response.SuccessResponse(c, fiber.StatusOK, "Upgrade Successfully Updated", data, nil)
}

func (h *UpgradeHandler) Route(open fiber.Router, userAuth fiber.Router, internalAuth fiber.Router) error {
	upgrade := internalAuth.Group("/upgrades")

	upgrade.Post("/", h.CreateUpgrade)
	upgrade.Get("/", h.GetUpgrades)
	upgrade.Get("/:id", h.GetUpgradeByID)
	upgrade.Put("/:id", h.UpdateUpgrade)

	return nil
}
