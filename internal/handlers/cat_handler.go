package handlers

import (
	"net/http"

	"github.com/gofiber/fiber/v2"
	"github.com/winartodev/cat-cafe/internal/controllers"
	"github.com/winartodev/cat-cafe/pkg/response"
)

type CatHandler struct {
	CatController controllers.CatController
}

func NewCatHandler(catController controllers.CatController) CatHandler {
	return CatHandler{
		CatController: catController,
	}
}

func (ch *CatHandler) GetCats(c *fiber.Ctx) error {
	res, err := ch.CatController.GetCatController()
	if err != nil {
		return response.FailedResponse(c, http.StatusInternalServerError, err)
	}

	return response.SuccessResponse(c, http.StatusOK, "SUCCESS", res, nil)
}

func (ch CatHandler) Route(route fiber.Router) error {

	route.Get("/cats", ch.GetCats)

	return nil
}
