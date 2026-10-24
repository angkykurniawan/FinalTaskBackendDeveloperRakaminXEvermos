package controllers

import (
	"evermos-mini/usecases"
	"evermos-mini/utils"

	"github.com/gofiber/fiber/v2"
)

type AuthController struct {
	authUsecase usecases.AuthUsecase
}

func NewAuthController(authUsecase usecases.AuthUsecase) *AuthController {
	return &AuthController{authUsecase: authUsecase}
}

// POST /auth/register
func (ctrl *AuthController) Register(c *fiber.Ctx) error {
	var input usecases.RegisterInput
	if err := c.BodyParser(&input); err != nil {
		return utils.ErrorResponse(c, fiber.StatusBadRequest, "Format payload pendaftaran tidak valid")
	}

	user, err := ctrl.authUsecase.Register(input)
	if err != nil {
		return utils.ErrorResponse(c, fiber.StatusBadRequest, err.Error())
	}

	return utils.SuccessResponse(c, fiber.StatusOK, "Registrasi berhasil", user)
}

// POST /auth/login
func (ctrl *AuthController) Login(c *fiber.Ctx) error {
	var input struct {
		NoTelp    string `json:"no_telp"`
		KataSandi string `json:"kata_sandi"`
	}

	if err := c.BodyParser(&input); err != nil {
		return utils.ErrorResponse(c, fiber.StatusBadRequest, "Format payload login tidak valid")
	}

	token, user, err := ctrl.authUsecase.Login(input.NoTelp, input.KataSandi)
	if err != nil {
		return utils.ErrorResponse(c, fiber.StatusUnauthorized, err.Error())
	}

	return utils.SuccessResponse(c, fiber.StatusOK, "Login berhasil", fiber.Map{
		"token": token,
		"user":  user,
	})
}
