package repositories

import (
	"evermos-mini/models"
	"strings"

	"gorm.io/gorm"
)

type ProductRepository interface {
	Create(product *models.Produk) error
	FindByID(id uint) (*models.Produk, error)
	FindAll(page, limit int, nama, categoryID, tokoID, minHarga, maxHarga string) ([]models.Produk, int64, error)
	Update(product *models.Produk) error
	Delete(id uint) error
	CreateLogProduk(log *models.LogProduk) error
	FindLogByID(id uint) (*models.LogProduk, error)
}

type productRepository struct {
	db *gorm.DB
}

func NewProductRepository(db *gorm.DB) ProductRepository {
	return &productRepository{db: db}
}

func (r *productRepository) Create(product *models.Produk) error {
	return r.db.Create(product).Error
}

func (r *productRepository) FindByID(id uint) (*models.Produk, error) {
	var product models.Produk
	err := r.db.Preload("Toko").Preload("Category").First(&product, id).Error
	if err != nil {
		return nil, err
	}
	return &product, nil
}

func (r *productRepository) FindAll(page, limit int, nama, categoryID, tokoID, minHarga, maxHarga string) ([]models.Produk, int64, error) {
	if page < 1 {
		page = 1
	}
	if limit < 1 {
		limit = 10
	}
	offset := (page - 1) * limit

	var products []models.Produk
	var total int64

	query := r.db.Model(&models.Produk{}).Preload("Toko").Preload("Category")

	if nama != "" {
		query = query.Where("LOWER(nama_produk) LIKE ?", "%"+strings.ToLower(nama)+"%")
	}
	if categoryID != "" {
		query = query.Where("category_id = ?", categoryID)
	}
	if tokoID != "" {
		query = query.Where("toko_id = ?", tokoID)
	}
	if minHarga != "" {
		query = query.Where("harga_konsumen >= ?", minHarga)
	}
	if maxHarga != "" {
		query = query.Where("harga_konsumen <= ?", maxHarga)
	}

	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	if err := query.Offset(offset).Limit(limit).Find(&products).Error; err != nil {
		return nil, 0, err
	}

	return products, total, nil
}

func (r *productRepository) Update(product *models.Produk) error {
	return r.db.Save(product).Error
}

func (r *productRepository) Delete(id uint) error {
	// Hapus log produk terkait atau unscoped jika dibutuhkan
	r.db.Where("produk_id = ?", id).Delete(&models.LogProduk{})
	return r.db.Delete(&models.Produk{}, id).Error
}

func (r *productRepository) CreateLogProduk(log *models.LogProduk) error {
	return r.db.Create(log).Error
}

func (r *productRepository) FindLogByID(id uint) (*models.LogProduk, error) {
	var log models.LogProduk
	err := r.db.Preload("Produk").Preload("Toko").Preload("Category").First(&log, id).Error
	if err != nil {
		return nil, err
	}
	return &log, nil
}
