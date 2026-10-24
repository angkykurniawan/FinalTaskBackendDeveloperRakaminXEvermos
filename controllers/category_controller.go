package controllers

import (
	"evermos-mini/usecases"
	"evermos-mini/utils"
	"strconv"

	"github.com/gofiber/fiber/v2"
)

type CategoryController struct {
	categoryUsecase usecases.CategoryUsecase
}

func NewCategoryController(categoryUsecase usecases.CategoryUsecase) *CategoryController {
	return &CategoryController{categoryUsecase: categoryUsecase}
}

// GET /category
func (ctrl *CategoryController) GetAllCategory(c *fiber.Ctx) error {
	page, _ := strconv.Atoi(c.Query("page", "1"))
	limit, _ := strconv.Atoi(c.Query("limit", "10"))
	search := c.Query("search", "")

	categories, total, err := ctrl.categoryUsecase.GetAllCategory(page, limit, search)
	if err != nil {
		return utils.ErrorResponse(c, fiber.StatusInternalServerError, err.Error())
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"status": true,
		"data":   categories,
		"meta": fiber.Map{
			"page":  page,
			"limit": limit,
			"total": total,
		},
	})
}

// GET /category/:id
func (ctrl *CategoryController) GetCategoryByID(c *fiber.Ctx) error {
	idParam := c.Params("id")
	id, err := strconv.ParseUint(idParam, 10, 32)
	if err != nil {
		return utils.ErrorResponse(c, fiber.StatusBadRequest, "ID kategori tidak valid")
	}

	cat, err := ctrl.categoryUsecase.GetCategoryByID(uint(id))
	if err != nil {
		return utils.ErrorResponse(c, fiber.StatusNotFound, err.Error())
	}

	return utils.SuccessResponse(c, fiber.StatusOK, "", cat)
}

// POST /category (Admin Only)
func (ctrl *CategoryController) CreateCategory(c *fiber.Ctx) error {
	var input struct {
		NamaCategory string `json:"nama_category"`
	}

	if err := c.BodyParser(&input); err != nil {
		return utils.ErrorResponse(c, fiber.StatusBadRequest, "Format payload kategori tidak valid")
	}

	cat, err := ctrl.categoryUsecase.CreateCategory(input.NamaCategory)
	if err != nil {
		return utils.ErrorResponse(c, fiber.StatusBadRequest, err.Error())
	}

	return utils.SuccessResponse(c, fiber.StatusOK, "Kategori berhasil ditambahkan", cat)
}

// PUT /category/:id (Admin Only)
func (ctrl *CategoryController) UpdateCategory(c *fiber.Ctx) error {
	idParam := c.Params("id")
	id, err := strconv.ParseUint(idParam, 10, 32)
	if err != nil {
		return utils.ErrorResponse(c, fiber.StatusBadRequest, "ID kategori tidak valid")
	}

	var input struct {
		NamaCategory string `json:"nama_category"`
	}

	if err := c.BodyParser(&input); err != nil {
		return utils.ErrorResponse(c, fiber.StatusBadRequest, "Format payload kategori tidak valid")
	}

	cat, err := ctrl.categoryUsecase.UpdateCategory(uint(id), input.NamaCategory)
	if err != nil {
		return utils.ErrorResponse(c, fiber.StatusBadRequest, err.Error())
	}

	return utils.SuccessResponse(c, fiber.StatusOK, "Kategori berhasil diperbarui", cat)
}

// DELETE /category/:id (Admin Only)
func (ctrl *CategoryController) DeleteCategory(c *fiber.Ctx) error {
	idParam := c.Params("id")
	id, err := strconv.ParseUint(idParam, 10, 32)
	if err != nil {
		return utils.ErrorResponse(c, fiber.StatusBadRequest, "ID kategori tidak valid")
	}

	if err := ctrl.categoryUsecase.DeleteCategory(uint(id)); err != nil {
		return utils.ErrorResponse(c, fiber.StatusBadRequest, err.Error())
	}

	return utils.SuccessResponse(c, fiber.StatusOK, "Kategori berhasil dihapus", nil)
}
