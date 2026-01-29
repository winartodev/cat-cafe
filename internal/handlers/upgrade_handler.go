package handlers

import (
	"github.com/gofiber/fiber/v2"
	"github.com/winartodev/cat-cafe/internal/dto"
	"github.com/winartodev/cat-cafe/internal/usecase"
	"github.com/winartodev/cat-cafe/pkg/apperror"
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

	if err := h.upgradeUC.CreateUpgrade(c.Context(), request.ToEntity()); err != nil {
		return response.FailedResponse(c, h.errorHandler, err)
	}

	return response.SuccessResponse(c, fiber.StatusOK, "Upgrade Successfully Created", nil, nil)
}

func (h *UpgradeHandler) Route(open fiber.Router, userAuth fiber.Router, internalAuth fiber.Router) error {
	internalAuth.Post("/upgrades", h.CreateUpgrade)

	return nil
}
