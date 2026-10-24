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

type ProductController struct {
	productUsecase usecases.ProductUsecase
}

func NewProductController(productUsecase usecases.ProductUsecase) *ProductController {
	return &ProductController{productUsecase: productUsecase}
}

// GET /product?nama_produk=&limit=&page=&category_id=&toko_id=&min_harga=&max_harga=
func (ctrl *ProductController) GetAllProduct(c *fiber.Ctx) error {
	page, _ := strconv.Atoi(c.Query("page", "1"))
	limit, _ := strconv.Atoi(c.Query("limit", "10"))
	nama := c.Query("nama_produk", "")
	categoryID := c.Query("category_id", "")
	tokoID := c.Query("toko_id", "")
	minHarga := c.Query("min_harga", "")
	maxHarga := c.Query("max_harga", "")

	products, total, err := ctrl.productUsecase.GetAllProduct(page, limit, nama, categoryID, tokoID, minHarga, maxHarga)
	if err != nil {
		return utils.ErrorResponse(c, fiber.StatusInternalServerError, err.Error())
	}

	return utils.PaginationResponse(c, fiber.StatusOK, page, limit, total, products)
}

// GET /product/:id
func (ctrl *ProductController) GetProductByID(c *fiber.Ctx) error {
	idParam := c.Params("id")
	id, err := strconv.ParseUint(idParam, 10, 32)
	if err != nil {
		return utils.ErrorResponse(c, fiber.StatusBadRequest, "ID produk tidak valid")
	}

	product, err := ctrl.productUsecase.GetProductByID(uint(id))
	if err != nil {
		return utils.ErrorResponse(c, fiber.StatusNotFound, err.Error())
	}

	return utils.SuccessResponse(c, fiber.StatusOK, "", product)
}

// POST /product
func (ctrl *ProductController) CreateProduct(c *fiber.Ctx) error {
	userID, ok := c.Locals("user_id").(uint)
	if !ok {
		return utils.ErrorResponse(c, fiber.StatusUnauthorized, "Sesi tidak valid")
	}

	namaProduk := c.FormValue("nama_produk")
	hargaReseller, _ := strconv.ParseFloat(c.FormValue("harga_reseller"), 64)
	hargaKonsumen, _ := strconv.ParseFloat(c.FormValue("harga_konsumen"), 64)
	stok, _ := strconv.Atoi(c.FormValue("stok"))
	deskripsi := c.FormValue("deskripsi")
	categoryID, _ := strconv.ParseUint(c.FormValue("category_id"), 10, 32)

	var gambarURL string
	// Support upload field bernama "photos" atau "photo"
	file, err := c.FormFile("photos")
	if err != nil || file == nil {
		file, err = c.FormFile("photo")
	}

	if err == nil && file != nil {
		_ = os.MkdirAll("./uploads", os.ModePerm)
		ext := filepath.Ext(file.Filename)
		filename := fmt.Sprintf("prod-%d-%d%s", userID, time.Now().UnixNano(), ext)
		savePath := fmt.Sprintf("./uploads/%s", filename)
		if err := c.SaveFile(file, savePath); err == nil {
			gambarURL = fmt.Sprintf("http://localhost:8080/uploads/%s", filename)
		}
	}

	product, err := ctrl.productUsecase.CreateProduct(userID, usecases.CreateProductInput{
		NamaProduk:    namaProduk,
		CategoryID:    uint(categoryID),
		HargaReseller: hargaReseller,
		HargaKonsumen: hargaKonsumen,
		Stok:          stok,
		Deskripsi:     deskripsi,
		Gambar:        gambarURL,
	})

	if err != nil {
		return utils.ErrorResponse(c, fiber.StatusBadRequest, err.Error())
	}

	return utils.SuccessResponse(c, fiber.StatusOK, "Produk berhasil ditambahkan", product)
}

// PUT /product/:id
func (ctrl *ProductController) UpdateProduct(c *fiber.Ctx) error {
	userID, ok := c.Locals("user_id").(uint)
	if !ok {
		return utils.ErrorResponse(c, fiber.StatusUnauthorized, "Sesi tidak valid")
	}

	idParam := c.Params("id")
	id, err := strconv.ParseUint(idParam, 10, 32)
	if err != nil {
		return utils.ErrorResponse(c, fiber.StatusBadRequest, "ID produk tidak valid")
	}

	namaProduk := c.FormValue("nama_produk")
	hargaReseller, _ := strconv.ParseFloat(c.FormValue("harga_reseller"), 64)
	hargaKonsumen, _ := strconv.ParseFloat(c.FormValue("harga_konsumen"), 64)
	stok, _ := strconv.Atoi(c.FormValue("stok"))
	deskripsi := c.FormValue("deskripsi")
	categoryID, _ := strconv.ParseUint(c.FormValue("category_id"), 10, 32)

	var gambarURL string
	file, err := c.FormFile("photos")
	if err != nil || file == nil {
		file, err = c.FormFile("photo")
	}

	if err == nil && file != nil {
		_ = os.MkdirAll("./uploads", os.ModePerm)
		ext := filepath.Ext(file.Filename)
		filename := fmt.Sprintf("prod-%d-%d%s", userID, time.Now().UnixNano(), ext)
		savePath := fmt.Sprintf("./uploads/%s", filename)
		if err := c.SaveFile(file, savePath); err == nil {
			gambarURL = fmt.Sprintf("http://localhost:8080/uploads/%s", filename)
		}
	}

	product, err := ctrl.productUsecase.UpdateProduct(userID, uint(id), usecases.UpdateProductInput{
		NamaProduk:    namaProduk,
		CategoryID:    uint(categoryID),
		HargaReseller: hargaReseller,
		HargaKonsumen: hargaKonsumen,
		Stok:          stok,
		Deskripsi:     deskripsi,
		Gambar:        gambarURL,
	})

	if err != nil {
		return utils.ErrorResponse(c, fiber.StatusBadRequest, err.Error())
	}

	return utils.SuccessResponse(c, fiber.StatusOK, "Produk berhasil diperbarui", product)
}

// DELETE /product/:id
func (ctrl *ProductController) DeleteProduct(c *fiber.Ctx) error {
	userID, ok := c.Locals("user_id").(uint)
	if !ok {
		return utils.ErrorResponse(c, fiber.StatusUnauthorized, "Sesi tidak valid")
	}

	idParam := c.Params("id")
	id, err := strconv.ParseUint(idParam, 10, 32)
	if err != nil {
		return utils.ErrorResponse(c, fiber.StatusBadRequest, "ID produk tidak valid")
	}

	if err := ctrl.productUsecase.DeleteProduct(userID, uint(id)); err != nil {
		return utils.ErrorResponse(c, fiber.StatusBadRequest, err.Error())
	}

	return utils.SuccessResponse(c, fiber.StatusOK, "Produk berhasil dihapus", nil)
}
