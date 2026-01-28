package middleware

import (
	"errors"
	"github.com/gofiber/fiber/v2"
	"github.com/winartodev/cat-cafe/internal/repositories"
	"github.com/winartodev/cat-cafe/pkg/apperror"
	"github.com/winartodev/cat-cafe/pkg/helper"
	"github.com/winartodev/cat-cafe/pkg/jwt"
	"github.com/winartodev/cat-cafe/pkg/response"
	"strings"
)

type Middleware interface {
	WithUserAuth() fiber.Handler
}

type middleware struct {
	jwtManager     *jwt.JWT
	userRepository repositories.UserRepository
	errorHandler   *apperror.ErrorHandler
}

func NewMiddleware(jwtManager *jwt.JWT, userRepository repositories.UserRepository) Middleware {
	return &middleware{
		jwtManager:     jwtManager,
		userRepository: userRepository,
		errorHandler:   apperror.NewErrorHandler(),
	}
}

func (m *middleware) WithUserAuth() fiber.Handler {
	return func(c *fiber.Ctx) error {
		ctx := c.Context()
		authHeader := c.Get("Authorization")
		if authHeader == "" {
			return response.FailedResponse(c, m.errorHandler, apperror.ErrMissingAuthHeader)
		}

		tokenString := strings.TrimPrefix(authHeader, "Bearer ")
		if tokenString == authHeader {
			return response.FailedResponse(c, m.errorHandler, apperror.ErrInvalidToken)
		}

		// Check blacklist
		if m.userRepository.IsTokenBlacklisted(ctx, tokenString) {
			return response.FailedResponse(c, m.errorHandler, apperror.ErrTokenRevoked)
		}

		// Validate token
		claims, err := m.jwtManager.ValidateToken(tokenString)
		if err != nil {
			if errors.Is(err, apperror.ErrTokenExpired) {
				return response.FailedResponse(c, m.errorHandler, apperror.ErrTokenExpired)
			}
			return response.FailedResponse(c, m.errorHandler, apperror.ErrInvalidToken)
		}

		userCache, err := m.userRepository.GetUserByIDDB(ctx, claims.UserID)
		if userCache == nil {
			return response.FailedResponse(c, m.errorHandler, apperror.ErrInvalidToken)
		}

		if err == nil {
			c.Locals(helper.ContextUserKey, userCache)
			c.Locals(helper.ContextUserIDKey, userCache.ID)
			c.Locals(helper.ContextEmailKey, userCache.Email)
		} else {
			c.Locals(helper.ContextUserIDKey, claims.UserID)
			c.Locals(helper.ContextEmailKey, claims.Email)
		}

		c.Locals(helper.ContextTokenKey, tokenString)

		return c.Next()
	}
}
