package helper

import (
	"errors"
	"fmt"
	"strconv"

	"github.com/gofiber/fiber/v2"
)

type Scalar interface {
	~int | ~int64 | ~string
}

func GetParam[T Scalar](c *fiber.Ctx, name string) (T, error) {
	p := c.Params(name)
	if p == "" {
		var zero T
		return zero, fmt.Errorf("parameter %s is missing", name)
	}

	var target T
	switch any(target).(type) {
	case string:
		return any(p).(T), nil
	case int:
		v, err := strconv.Atoi(p)
		if err != nil {
			return target, fmt.Errorf("invalid int: %w", err)
		}
		return any(v).(T), nil
	case int64:
		v, err := strconv.ParseInt(p, 10, 64)
		if err != nil {
			return target, fmt.Errorf("invalid int64: %w", err)
		}
		return any(v).(T), nil
	default:
		return target, errors.New("unsupported type")
	}
}
