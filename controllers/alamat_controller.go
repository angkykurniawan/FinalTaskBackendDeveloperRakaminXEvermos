package controllers

import (
	"evermos-mini/usecases"
	"evermos-mini/utils"
	"strconv"

	"github.com/gofiber/fiber/v2"
)

type AlamatController struct {
	alamatUsecase usecases.AlamatUsecase
}

func NewAlamatController(alamatUsecase usecases.AlamatUsecase) *AlamatController {
	return &AlamatController{alamatUsecase: alamatUsecase}
}

// GET /user/alamat
func (ctrl *AlamatController) GetAllAlamat(c *fiber.Ctx) error {
	userID, ok := c.Locals("user_id").(uint)
	if !ok {
		return utils.ErrorResponse(c, fiber.StatusUnauthorized, "Sesi tidak valid")
	}

	list, err := ctrl.alamatUsecase.GetAllAlamat(userID)
	if err != nil {
		return utils.ErrorResponse(c, fiber.StatusInternalServerError, err.Error())
	}

	return utils.SuccessResponse(c, fiber.StatusOK, "", list)
}

// GET /user/alamat/:id
func (ctrl *AlamatController) GetAlamatByID(c *fiber.Ctx) error {
	userID, ok := c.Locals("user_id").(uint)
	if !ok {
		return utils.ErrorResponse(c, fiber.StatusUnauthorized, "Sesi tidak valid")
	}

	idParam := c.Params("id")
	id, err := strconv.ParseUint(idParam, 10, 32)
	if err != nil {
		return utils.ErrorResponse(c, fiber.StatusBadRequest, "ID alamat tidak valid")
	}

	alamat, err := ctrl.alamatUsecase.GetAlamatByID(uint(id), userID)
	if err != nil {
		return utils.ErrorResponse(c, fiber.StatusNotFound, err.Error())
	}

	return utils.SuccessResponse(c, fiber.StatusOK, "", alamat)
}

// POST /user/alamat
func (ctrl *AlamatController) CreateAlamat(c *fiber.Ctx) error {
	userID, ok := c.Locals("user_id").(uint)
	if !ok {
		return utils.ErrorResponse(c, fiber.StatusUnauthorized, "Sesi tidak valid")
	}

	var input usecases.AlamatInput
	if err := c.BodyParser(&input); err != nil {
		return utils.ErrorResponse(c, fiber.StatusBadRequest, "Format payload alamat tidak valid")
	}

	alamat, err := ctrl.alamatUsecase.CreateAlamat(userID, input)
	if err != nil {
		return utils.ErrorResponse(c, fiber.StatusBadRequest, err.Error())
	}

	return utils.SuccessResponse(c, fiber.StatusOK, "Alamat berhasil ditambahkan", alamat)
}

// PUT /user/alamat/:id
func (ctrl *AlamatController) UpdateAlamat(c *fiber.Ctx) error {
	userID, ok := c.Locals("user_id").(uint)
	if !ok {
		return utils.ErrorResponse(c, fiber.StatusUnauthorized, "Sesi tidak valid")
	}

	idParam := c.Params("id")
	id, err := strconv.ParseUint(idParam, 10, 32)
	if err != nil {
		return utils.ErrorResponse(c, fiber.StatusBadRequest, "ID alamat tidak valid")
	}

	var input usecases.AlamatInput
	if err := c.BodyParser(&input); err != nil {
		return utils.ErrorResponse(c, fiber.StatusBadRequest, "Format payload pembaruan alamat tidak valid")
	}

	alamat, err := ctrl.alamatUsecase.UpdateAlamat(uint(id), userID, input)
	if err != nil {
		return utils.ErrorResponse(c, fiber.StatusBadRequest, err.Error())
	}

	return utils.SuccessResponse(c, fiber.StatusOK, "Alamat berhasil diperbarui", alamat)
}

// DELETE /user/alamat/:id
func (ctrl *AlamatController) DeleteAlamat(c *fiber.Ctx) error {
	userID, ok := c.Locals("user_id").(uint)
	if !ok {
		return utils.ErrorResponse(c, fiber.StatusUnauthorized, "Sesi tidak valid")
	}

	idParam := c.Params("id")
	id, err := strconv.ParseUint(idParam, 10, 32)
	if err != nil {
		return utils.ErrorResponse(c, fiber.StatusBadRequest, "ID alamat tidak valid")
	}

	if err := ctrl.alamatUsecase.DeleteAlamat(uint(id), userID); err != nil {
		return utils.ErrorResponse(c, fiber.StatusBadRequest, err.Error())
	}

	return utils.SuccessResponse(c, fiber.StatusOK, "Alamat berhasil dihapus", nil)
}
