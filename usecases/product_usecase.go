package usecases

import (
	"errors"
	"evermos-mini/models"
	"evermos-mini/repositories"
	"fmt"
	"strings"
)

type CreateProductInput struct {
	NamaProduk    string
	CategoryID    uint
	HargaReseller float64
	HargaKonsumen float64
	Stok          int
	Deskripsi     string
	Gambar        string
}

type UpdateProductInput struct {
	NamaProduk    string
	CategoryID    uint
	HargaReseller float64
	HargaKonsumen float64
	Stok          int
	Deskripsi     string
	Gambar        string
}

type ProductUsecase interface {
	GetAllProduct(page, limit int, nama, categoryID, tokoID, minHarga, maxHarga string) ([]models.Produk, int64, error)
	GetProductByID(id uint) (*models.Produk, error)
	CreateProduct(userID uint, input CreateProductInput) (*models.Produk, error)
	UpdateProduct(userID, id uint, input UpdateProductInput) (*models.Produk, error)
	DeleteProduct(userID, id uint) error
}

type productUsecase struct {
	productRepo  repositories.ProductRepository
	tokoRepo     repositories.TokoRepository
	categoryRepo repositories.CategoryRepository
}

func NewProductUsecase(
	productRepo repositories.ProductRepository,
	tokoRepo repositories.TokoRepository,
	categoryRepo repositories.CategoryRepository,
) ProductUsecase {
	return &productUsecase{
		productRepo:  productRepo,
		tokoRepo:     tokoRepo,
		categoryRepo: categoryRepo,
	}
}

func (u *productUsecase) GetAllProduct(page, limit int, nama, categoryID, tokoID, minHarga, maxHarga string) ([]models.Produk, int64, error) {
	return u.productRepo.FindAll(page, limit, nama, categoryID, tokoID, minHarga, maxHarga)
}

func (u *productUsecase) GetProductByID(id uint) (*models.Produk, error) {
	prod, err := u.productRepo.FindByID(id)
	if err != nil {
		return nil, errors.New("produk tidak ditemukan")
	}
	return prod, nil
}

func (u *productUsecase) CreateProduct(userID uint, input CreateProductInput) (*models.Produk, error) {
	toko, err := u.tokoRepo.FindByUserID(userID)
	if err != nil {
		return nil, errors.New("toko milik Anda tidak ditemukan")
	}

	input.NamaProduk = strings.TrimSpace(input.NamaProduk)
	if input.NamaProduk == "" {
		return nil, errors.New("nama produk wajib diisi")
	}

	if input.CategoryID > 0 {
		_, err := u.categoryRepo.FindByID(input.CategoryID)
		if err != nil {
			return nil, errors.New("kategori yang dipilih tidak valid")
		}
	} else {
		return nil, errors.New("kategori produk wajib dipilih")
	}

	slug := strings.ToLower(strings.ReplaceAll(input.NamaProduk, " ", "-"))

	product := models.Produk{
		NamaProduk:    input.NamaProduk,
		Slug:          slug,
		HargaReseller: input.HargaReseller,
		HargaKonsumen: input.HargaKonsumen,
		Stok:          input.Stok,
		Deskripsi:     input.Deskripsi,
		Gambar:        input.Gambar,
		TokoID:        toko.ID,
		CategoryID:    input.CategoryID,
	}

	if err := u.productRepo.Create(&product); err != nil {
		return nil, fmt.Errorf("gagal menambahkan produk: %v", err)
	}

	// Catat log snapshot produk
	logSnapshot := models.LogProduk{
		ProdukID:      product.ID,
		NamaProduk:    product.NamaProduk,
		Slug:          product.Slug,
		HargaReseller: product.HargaReseller,
		HargaKonsumen: product.HargaKonsumen,
		Deskripsi:     product.Deskripsi,
		TokoID:        product.TokoID,
		CategoryID:    product.CategoryID,
	}
	_ = u.productRepo.CreateLogProduk(&logSnapshot)

	// Muat relasi lengkap
	fullProd, err := u.productRepo.FindByID(product.ID)
	if err == nil {
		return fullProd, nil
	}
	return &product, nil
}

func (u *productUsecase) UpdateProduct(userID, id uint, input UpdateProductInput) (*models.Produk, error) {
	toko, err := u.tokoRepo.FindByUserID(userID)
	if err != nil {
		return nil, errors.New("toko milik Anda tidak ditemukan")
	}

	product, err := u.productRepo.FindByID(id)
	if err != nil {
		return nil, errors.New("produk tidak ditemukan")
	}

	// Proteksi: User tidak dapat mengelola produk dari toko user lain
	if product.TokoID != toko.ID {
		return nil, errors.New("akses ditolak: Anda tidak memiliki wewenang untuk mengubah produk ini")
	}

	if strings.TrimSpace(input.NamaProduk) != "" {
		product.NamaProduk = strings.TrimSpace(input.NamaProduk)
	}
	if input.CategoryID > 0 {
		_, err := u.categoryRepo.FindByID(input.CategoryID)
		if err == nil {
			product.CategoryID = input.CategoryID
		}
	}
	if input.HargaReseller > 0 {
		product.HargaReseller = input.HargaReseller
	}
	if input.HargaKonsumen > 0 {
		product.HargaKonsumen = input.HargaKonsumen
	}
	if input.Stok >= 0 {
		product.Stok = input.Stok
	}
	if input.Deskripsi != "" {
		product.Deskripsi = input.Deskripsi
	}
	if input.Gambar != "" {
		product.Gambar = input.Gambar
	}

	if err := u.productRepo.Update(product); err != nil {
		return nil, errors.New("gagal memperbarui produk")
	}

	// Buat log produk snapshot terbaru
	logSnapshot := models.LogProduk{
		ProdukID:      product.ID,
		NamaProduk:    product.NamaProduk,
		Slug:          product.Slug,
		HargaReseller: product.HargaReseller,
		HargaKonsumen: product.HargaKonsumen,
		Deskripsi:     product.Deskripsi,
		TokoID:        product.TokoID,
		CategoryID:    product.CategoryID,
	}
	_ = u.productRepo.CreateLogProduk(&logSnapshot)

	fullProd, _ := u.productRepo.FindByID(product.ID)
	if fullProd != nil {
		return fullProd, nil
	}
	return product, nil
}

func (u *productUsecase) DeleteProduct(userID, id uint) error {
	toko, err := u.tokoRepo.FindByUserID(userID)
	if err != nil {
		return errors.New("toko tidak ditemukan")
	}

	product, err := u.productRepo.FindByID(id)
	if err != nil {
		return errors.New("produk tidak ditemukan")
	}

	// Proteksi: User tidak dapat mengelola produk toko user lain
	if product.TokoID != toko.ID {
		return errors.New("akses ditolak: Anda tidak memiliki wewenang untuk menghapus produk ini")
	}

	return u.productRepo.Delete(id)
}
