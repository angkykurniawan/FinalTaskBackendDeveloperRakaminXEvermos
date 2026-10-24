package controllers

import (
	"evermos-mini/usecases"
	"evermos-mini/utils"

	"github.com/gofiber/fiber/v2"
)

type UserController struct {
	userUsecase usecases.UserUsecase
}

func NewUserController(userUsecase usecases.UserUsecase) *UserController {
	return &UserController{userUsecase: userUsecase}
}

// GET /user
func (ctrl *UserController) GetMyProfile(c *fiber.Ctx) error {
	userID, ok := c.Locals("user_id").(uint)
	if !ok {
		return utils.ErrorResponse(c, fiber.StatusUnauthorized, "Sesi pengguna tidak valid")
	}

	user, err := ctrl.userUsecase.GetProfile(userID)
	if err != nil {
		return utils.ErrorResponse(c, fiber.StatusNotFound, err.Error())
	}

	return utils.SuccessResponse(c, fiber.StatusOK, "", user)
}

// PUT /user
func (ctrl *UserController) UpdateMyProfile(c *fiber.Ctx) error {
	userID, ok := c.Locals("user_id").(uint)
	if !ok {
		return utils.ErrorResponse(c, fiber.StatusUnauthorized, "Sesi pengguna tidak valid")
	}

	var input usecases.UpdateUserInput
	if err := c.BodyParser(&input); err != nil {
		return utils.ErrorResponse(c, fiber.StatusBadRequest, "Format payload pembaruan profil tidak valid")
	}

	user, err := ctrl.userUsecase.UpdateProfile(userID, input)
	if err != nil {
		return utils.ErrorResponse(c, fiber.StatusBadRequest, err.Error())
	}

	return utils.SuccessResponse(c, fiber.StatusOK, "Profil berhasil diperbarui", user)
}
