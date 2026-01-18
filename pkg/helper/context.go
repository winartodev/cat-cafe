package helper

import (
	"context"
	"github.com/gofiber/fiber/v2"
)

const (
	ContextUserKey   = "user"
	ContextTokenKey  = "token"
	ContextEmailKey  = "email"
	ContextUserIDKey = "userID"
)

func GetUserID(c *fiber.Ctx) int64 {
	val := c.Locals(ContextUserIDKey)
	if id, ok := val.(int64); ok {
		return id
	}
	return 0
}

func GetEmail(c *fiber.Ctx) string {
	val, _ := c.Locals(ContextEmailKey).(string)
	return val
}

func GetToken(c *fiber.Ctx) string {
	val, _ := c.Locals(ContextTokenKey).(string)
	return val
}

func GetUserIDFromContext(ctx context.Context) (int64, error) {
	userID, ok := ctx.Value(ContextUserIDKey).(int64)
	if !ok || userID <= 0 {
		return 0, nil
	}
	return userID, nil
}
