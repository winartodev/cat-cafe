package handlers

import (
	"github.com/gofiber/fiber/v2"
	"github.com/winartodev/cat-cafe/internal/dto"
	"github.com/winartodev/cat-cafe/internal/middleware"
	"github.com/winartodev/cat-cafe/internal/usecase"
	"github.com/winartodev/cat-cafe/pkg/apperror"
	"github.com/winartodev/cat-cafe/pkg/helper"
	"github.com/winartodev/cat-cafe/pkg/response"
)

type AuthHandler struct {
	AuthUseCase usecase.AuthUseCase
}

func NewAuthHandler(authUseCase usecase.AuthUseCase) *AuthHandler {
	return &AuthHandler{
		AuthUseCase: authUseCase,
	}
}

func (a *AuthHandler) Login(c *fiber.Ctx) error {
	authCode := c.Query("auth_code")
	if authCode == "" {
		return response.FailedResponse(c, fiber.StatusBadRequest, apperror.ErrBadRequest)
	}

	token, user, gameData, err := a.AuthUseCase.Login(c.Context(), authCode)
	if err != nil {
		return response.FailedResponse(c, fiber.StatusInternalServerError, err)
	}

	return response.SuccessResponse(c, fiber.StatusOK, "Login Success", dto.ToLoginResponse(token, user, gameData), nil)
}

func (a *AuthHandler) Logout(c *fiber.Ctx) error {
	token := helper.GetToken(c)
	userID := helper.GetUserID(c)

	if err := a.AuthUseCase.Logout(c.Context(), token, userID); err != nil {
		return response.FailedResponse(c, fiber.StatusInternalServerError, err)
	}

	return response.SuccessResponse(c, fiber.StatusOK, "Logout Success", nil, nil)
}

func (a *AuthHandler) GetUserData(c *fiber.Ctx) error {
	userID := helper.GetUserID(c)

	res, err := a.AuthUseCase.GetUserByID(c.Context(), userID)
	if err != nil {
		return response.FailedResponse(c, fiber.StatusInternalServerError, err)
	}

	return response.SuccessResponse(c, fiber.StatusOK, "Success Get User Data", dto.ToUserResponse(res), nil)
}

func (a *AuthHandler) Route(route fiber.Router, m middleware.Middleware) error {
	auth := route.Group("/auth")

	auth.Post("/login", a.Login)
	auth.Post("/logout", m.WithUserAuth(a.Logout))
	auth.Get("/me", m.WithUserAuth(a.GetUserData))

	return nil
}
