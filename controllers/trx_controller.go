package controllers

import (
	"evermos-mini/usecases"
	"evermos-mini/utils"
	"strconv"

	"github.com/gofiber/fiber/v2"
)

type TrxController struct {
	trxUsecase usecases.TrxUsecase
}

func NewTrxController(trxUsecase usecases.TrxUsecase) *TrxController {
	return &TrxController{trxUsecase: trxUsecase}
}

// POST /trx
func (ctrl *TrxController) CreateTrx(c *fiber.Ctx) error {
	userID, ok := c.Locals("user_id").(uint)
	if !ok {
		return utils.ErrorResponse(c, fiber.StatusUnauthorized, "Sesi tidak valid")
	}

	var input usecases.TrxCreateInput
	if err := c.BodyParser(&input); err != nil {
		return utils.ErrorResponse(c, fiber.StatusBadRequest, "Format payload transaksi tidak valid")
	}

	trx, err := ctrl.trxUsecase.CreateTrx(userID, input)
	if err != nil {
		return utils.ErrorResponse(c, fiber.StatusBadRequest, err.Error())
	}

	return utils.SuccessResponse(c, fiber.StatusOK, "Transaksi berhasil dibuat", trx)
}

// GET /trx
func (ctrl *TrxController) GetAllTrx(c *fiber.Ctx) error {
	userID, ok := c.Locals("user_id").(uint)
	if !ok {
		return utils.ErrorResponse(c, fiber.StatusUnauthorized, "Sesi tidak valid")
	}

	list, err := ctrl.trxUsecase.GetAllTrx(userID)
	if err != nil {
		return utils.ErrorResponse(c, fiber.StatusInternalServerError, err.Error())
	}

	return utils.SuccessResponse(c, fiber.StatusOK, "", list)
}

// GET /trx/:id
func (ctrl *TrxController) GetTrxByID(c *fiber.Ctx) error {
	userID, ok := c.Locals("user_id").(uint)
	if !ok {
		return utils.ErrorResponse(c, fiber.StatusUnauthorized, "Sesi tidak valid")
	}

	idParam := c.Params("id")
	id, err := strconv.ParseUint(idParam, 10, 32)
	if err != nil {
		return utils.ErrorResponse(c, fiber.StatusBadRequest, "ID transaksi tidak valid")
	}

	trx, err := ctrl.trxUsecase.GetTrxByID(uint(id), userID)
	if err != nil {
		return utils.ErrorResponse(c, fiber.StatusNotFound, err.Error())
	}

	return utils.SuccessResponse(c, fiber.StatusOK, "", trx)
}

// PUT /trx/:id/status
func (ctrl *TrxController) UpdateStatusTrx(c *fiber.Ctx) error {
	idParam := c.Params("id")
	id, err := strconv.ParseUint(idParam, 10, 32)
	if err != nil {
		return utils.ErrorResponse(c, fiber.StatusBadRequest, "ID transaksi tidak valid")
	}

	var input struct {
		Status string `json:"status"`
	}
	if err := c.BodyParser(&input); err != nil || input.Status == "" {
		return utils.ErrorResponse(c, fiber.StatusBadRequest, "Status transaksi wajib diisi")
	}

	trx, err := ctrl.trxUsecase.UpdateStatusTrx(uint(id), input.Status)
	if err != nil {
		return utils.ErrorResponse(c, fiber.StatusBadRequest, err.Error())
	}

	return utils.SuccessResponse(c, fiber.StatusOK, "Status transaksi berhasil diperbarui", trx)
}
