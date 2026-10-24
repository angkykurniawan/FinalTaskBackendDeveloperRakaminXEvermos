package middleware

import (
	"evermos-mini/utils"
	"strings"

	"github.com/gofiber/fiber/v2"
)

func JWTProtected(c *fiber.Ctx) error {
	authHeader := c.Get("Authorization")
	if authHeader == "" {
		authHeader = c.Get("token")
	}

	if authHeader == "" {
		return utils.ErrorResponse(c, fiber.StatusUnauthorized, "Token otentikasi tidak ditemukan")
	}

	tokenString := strings.TrimSpace(authHeader)
	if strings.HasPrefix(tokenString, "Bearer ") {
		tokenString = strings.TrimPrefix(tokenString, "Bearer ")
	} else if strings.HasPrefix(tokenString, "bearer ") {
		tokenString = strings.TrimPrefix(tokenString, "bearer ")
	}

	if tokenString == "" {
		return utils.ErrorResponse(c, fiber.StatusUnauthorized, "Token otentikasi kosong")
	}

	userID, role, err := utils.ValidateToken(tokenString)
	if err != nil {
		return utils.ErrorResponse(c, fiber.StatusUnauthorized, "Token otentikasi tidak valid atau telah kadaluarsa")
	}

	c.Locals("user_id", userID)
	c.Locals("role", role)
	return c.Next()
}
