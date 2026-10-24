package utils

import (
	"math"

	"github.com/gofiber/fiber/v2"
)

func SuccessResponse(c *fiber.Ctx, statusCode int, message string, data interface{}) error {
	res := fiber.Map{
		"status": true,
	}
	if message != "" {
		res["message"] = message
	}
	if data != nil {
		res["data"] = data
	}
	return c.Status(statusCode).JSON(res)
}

func ErrorResponse(c *fiber.Ctx, statusCode int, message string) error {
	return c.Status(statusCode).JSON(fiber.Map{
		"status":  false,
		"message": message,
	})
}

func PaginationResponse(c *fiber.Ctx, statusCode int, page, limit int, total int64, data interface{}) error {
	totalPages := 0
	if limit > 0 {
		totalPages = int(math.Ceil(float64(total) / float64(limit)))
	}

	return c.Status(statusCode).JSON(fiber.Map{
		"status": true,
		"pagination": fiber.Map{
			"page":        page,
			"limit":       limit,
			"total_data":  total,
			"total_pages": totalPages,
		},
		"data": data,
	})
}
