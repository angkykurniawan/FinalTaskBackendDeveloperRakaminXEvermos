package controllers

import (
	"evermos-mini/usecases"
	"evermos-mini/utils"
	"fmt"
	"os"
	"path/filepath"
	"strconv"
	"time"

	"github.com/gofiber/fiber/v2"
)

type TokoController struct {
	tokoUsecase usecases.TokoUsecase
}

func NewTokoController(tokoUsecase usecases.TokoUsecase) *TokoController {
	return &TokoController{tokoUsecase: tokoUsecase}
}

// GET /toko
func (ctrl *TokoController) GetAllToko(c *fiber.Ctx) error {
	page, _ := strconv.Atoi(c.Query("page", "1"))
	limit, _ := strconv.Atoi(c.Query("limit", "10"))
	nama := c.Query("nama", "")

	tokos, total, err := ctrl.tokoUsecase.GetAllToko(page, limit, nama)
	if err != nil {
		return utils.ErrorResponse(c, fiber.StatusInternalServerError, err.Error())
	}

	return utils.PaginationResponse(c, fiber.StatusOK, page, limit, total, tokos)
}

// GET /toko/my
func (ctrl *TokoController) GetMyToko(c *fiber.Ctx) error {
	userID, ok := c.Locals("user_id").(uint)
	if !ok {
		return utils.ErrorResponse(c, fiber.StatusUnauthorized, "Sesi tidak valid")
	}

	toko, err := ctrl.tokoUsecase.GetMyToko(userID)
	if err != nil {
		return utils.ErrorResponse(c, fiber.StatusNotFound, err.Error())
	}

	return utils.SuccessResponse(c, fiber.StatusOK, "", toko)
}

// GET /toko/:id_toko
func (ctrl *TokoController) GetTokoByID(c *fiber.Ctx) error {
	idParam := c.Params("id_toko")
	id, err := strconv.ParseUint(idParam, 10, 32)
	if err != nil {
		return utils.ErrorResponse(c, fiber.StatusBadRequest, "ID toko tidak valid")
	}

	toko, err := ctrl.tokoUsecase.GetTokoByID(uint(id))
	if err != nil {
		return utils.ErrorResponse(c, fiber.StatusNotFound, err.Error())
	}

	return utils.SuccessResponse(c, fiber.StatusOK, "", toko)
}

// PUT /toko/:id_toko (form-data: nama_toko, photo)
func (ctrl *TokoController) UpdateToko(c *fiber.Ctx) error {
	userID, ok := c.Locals("user_id").(uint)
	if !ok {
		return utils.ErrorResponse(c, fiber.StatusUnauthorized, "Sesi tidak valid")
	}

	idParam := c.Params("id_toko")
	id, err := strconv.ParseUint(idParam, 10, 32)
	if err != nil {
		return utils.ErrorResponse(c, fiber.StatusBadRequest, "ID toko tidak valid")
	}

	namaToko := c.FormValue("nama_toko")
	var photoURL string

	file, err := c.FormFile("photo")
	if err == nil && file != nil {
		_ = os.MkdirAll("./uploads", os.ModePerm)
		ext := filepath.Ext(file.Filename)
		filename := fmt.Sprintf("toko-%d-%d%s", userID, time.Now().Unix(), ext)
		savePath := fmt.Sprintf("./uploads/%s", filename)
		if err := c.SaveFile(file, savePath); err == nil {
			photoURL = fmt.Sprintf("http://localhost:8080/uploads/%s", filename)
		}
	}

	toko, err := ctrl.tokoUsecase.UpdateToko(userID, uint(id), namaToko, photoURL)
	if err != nil {
		return utils.ErrorResponse(c, fiber.StatusBadRequest, err.Error())
	}

	return utils.SuccessResponse(c, fiber.StatusOK, "Toko berhasil diperbarui", toko)
}
