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

type PaginationParams struct {
	Page   int
	Limit  int
	Offset int
}

type PaginationMeta struct {
	Page       int   `json:"page"`
	Limit      int   `json:"limit"`
	TotalRows  int64 `json:"total_rows"`
	TotalPages int   `json:"total_pages"`
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

func GetPaginationParams(c *fiber.Ctx) *PaginationParams {
	page := c.QueryInt("page", 1)
	limit := c.QueryInt("limit", 10)

	// Validate page
	if page < 1 {
		page = 1
	}

	// Validate limit (min: 1, max: 100)
	if limit < 1 || limit > 100 {
		limit = 10
	}

	// Calculate offset
	offset := (page - 1) * limit

	return &PaginationParams{
		Page:   page,
		Limit:  limit,
		Offset: offset,
	}
}

func CreatePaginationMeta(page, limit int, totalRows int64) *PaginationMeta {
	totalPages := int(totalRows) / limit
	if int(totalRows)%limit > 0 {
		totalPages++
	}

	return &PaginationMeta{
		Page:       page,
		Limit:      limit,
		TotalRows:  totalRows,
		TotalPages: totalPages,
	}
}
