package repositories

import (
	"evermos-mini/models"
	"strings"

	"gorm.io/gorm"
)

type TokoRepository interface {
	Create(toko *models.Toko) error
	FindByUserID(userID uint) (*models.Toko, error)
	FindByID(id uint) (*models.Toko, error)
	FindAll(page, limit int, nama string) ([]models.Toko, int64, error)
	Update(toko *models.Toko) error
}

type tokoRepository struct {
	db *gorm.DB
}

func NewTokoRepository(db *gorm.DB) TokoRepository {
	return &tokoRepository{db: db}
}

func (r *tokoRepository) Create(toko *models.Toko) error {
	return r.db.Create(toko).Error
}

func (r *tokoRepository) FindByUserID(userID uint) (*models.Toko, error) {
	var toko models.Toko
	err := r.db.Where("user_id = ?", userID).Preload("User").First(&toko).Error
	if err != nil {
		return nil, err
	}
	return &toko, nil
}

func (r *tokoRepository) FindByID(id uint) (*models.Toko, error) {
	var toko models.Toko
	err := r.db.Preload("User").First(&toko, id).Error
	if err != nil {
		return nil, err
	}
	return &toko, nil
}

func (r *tokoRepository) FindAll(page, limit int, nama string) ([]models.Toko, int64, error) {
	if page < 1 {
		page = 1
	}
	if limit < 1 {
		limit = 10
	}
	offset := (page - 1) * limit

	var tokos []models.Toko
	var total int64

	query := r.db.Model(&models.Toko{}).Preload("User")
	if nama != "" {
		query = query.Where("LOWER(nama_toko) LIKE ?", "%"+strings.ToLower(nama)+"%")
	}

	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	if err := query.Offset(offset).Limit(limit).Find(&tokos).Error; err != nil {
		return nil, 0, err
	}

	return tokos, total, nil
}

func (r *tokoRepository) Update(toko *models.Toko) error {
	return r.db.Save(toko).Error
}
