package handlers

import (
	"github.com/gofiber/fiber/v2"
	"github.com/winartodev/cat-cafe/internal/dto"
	"github.com/winartodev/cat-cafe/internal/usecase"
	"github.com/winartodev/cat-cafe/pkg/apperror"
	"github.com/winartodev/cat-cafe/pkg/helper"
	"github.com/winartodev/cat-cafe/pkg/response"
	"net/http"
)

type TutorialHandler struct {
	errorHandler    *apperror.ErrorHandler
	tutorialUseCase usecase.TutorialUseCase
}

func NewTutorialHandler(tutorialUC usecase.TutorialUseCase) *TutorialHandler {
	return &TutorialHandler{
		errorHandler:    apperror.NewErrorHandler(),
		tutorialUseCase: tutorialUC,
	}
}

func (t *TutorialHandler) CreateTutorials(c *fiber.Ctx) error {
	var request dto.TutorialDTO
	if err := c.BodyParser(&request); err != nil {
		return response.FailedResponse(c, t.errorHandler, err)
	}

	ctx := c.Context()
	sequence, translations := request.ToEntity()
	resSeq, resTutorial, err := t.tutorialUseCase.CreateTutorial(ctx, sequence, translations)
	if err != nil {
		return response.FailedResponse(c, t.errorHandler, err)
	}

	return response.SuccessResponse(c, http.StatusCreated, "Tutorial Successfully Created", dto.ToCreateTutorialResponse(resSeq, resTutorial), nil)
}

func (t *TutorialHandler) GetTutorials(c *fiber.Ctx) error {
	ctx := c.Context()
	params := helper.GetPaginationParams(c)

	res, totalRows, err := t.tutorialUseCase.GetTutorials(ctx, params.Limit, params.Offset)
	if err != nil {
		return response.FailedResponse(c, t.errorHandler, err)
	}

	meta := helper.CreatePaginationMeta(params.Page, params.Limit, totalRows)

	return response.SuccessResponse(c, http.StatusOK, "Tutorial Successfully Retrieved", dto.ToDetailTutorialsResponse(res), meta)
}

func (t *TutorialHandler) GetTranslations(c *fiber.Ctx) error {
	key, err := helper.GetParam[string](c, "key")
	if err != nil {
		return response.FailedResponse(c, t.errorHandler, err)
	}

	ctx := c.Context()
	params := helper.GetPaginationParams(c)

	res, totalRows, err := t.tutorialUseCase.GetTranslationsByTutorialKey(ctx, key, params.Limit, params.Offset)
	if err != nil {
		return response.FailedResponse(c, t.errorHandler, err)
	}

	meta := helper.CreatePaginationMeta(params.Page, params.Limit, totalRows)

	return response.SuccessResponse(c, http.StatusOK, "Tutorial Successfully Retrieved", dto.ToDetailTutorialsResponse(res), meta)
}

func (t *TutorialHandler) GetTranslationByID(c *fiber.Ctx) error {
	//id, err := helper.GetParam[int64](c, "id")
	//if err != nil {
	//	return response.FailedResponse(c, t.errorHandler, err)
	//}
	//
	//ctx := c.Context()
	//params := helper.GetPaginationParams(c)

	return response.SuccessResponse(c, http.StatusOK, "", nil, nil)

}

func (t *TutorialHandler) UpdateTutorial(c *fiber.Ctx) error {
	return response.SuccessResponse(c, http.StatusOK, "", nil, nil)
}

func (t *TutorialHandler) Route(open fiber.Router, userAuth fiber.Router, internalAuth fiber.Router) error {
	tutorials := internalAuth.Group("/tutorials")

	tutorials.Post("/", t.CreateTutorials)
	tutorials.Get("/", t.GetTutorials)
	tutorials.Get("/:key/translations", t.GetTranslations)
	tutorials.Get("/:key/translations/:id", t.GetTranslationByID)
	tutorials.Put("/:key/translations/:id", t.GetTranslationByID)
	return nil
}
