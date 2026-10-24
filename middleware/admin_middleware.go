package middleware

import (
	"evermos-mini/config"
	"evermos-mini/models"
	"evermos-mini/utils"

	"github.com/gofiber/fiber/v2"
)

func AdminOnly(c *fiber.Ctx) error {
	userLocals := c.Locals("user_id")
	if userLocals == nil {
		return utils.ErrorResponse(c, fiber.StatusUnauthorized, "Sesi pengguna tidak valid")
	}

	userID, ok := userLocals.(uint)
	if !ok {
		return utils.ErrorResponse(c, fiber.StatusUnauthorized, "Format identitas pengguna tidak valid")
	}

	// Cek langsung status role admin terkini di database
	var user models.User
	if err := config.DB.Select("id", "role").First(&user, userID).Error; err != nil {
		return utils.ErrorResponse(c, fiber.StatusUnauthorized, "Data pengguna tidak ditemukan")
	}

	if user.Role != "admin" {
		return utils.ErrorResponse(c, fiber.StatusForbidden, "Akses ditolak: Fitur ini hanya dapat diakses oleh Admin")
	}

	return c.Next()
}
