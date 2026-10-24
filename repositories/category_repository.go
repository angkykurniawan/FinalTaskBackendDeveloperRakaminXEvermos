package repositories

import (
	"evermos-mini/models"
	"strings"

	"gorm.io/gorm"
)

type CategoryRepository interface {
	Create(cat *models.Category) error
	FindAll(page, limit int, search string) ([]models.Category, int64, error)
	FindByID(id uint) (*models.Category, error)
	Update(cat *models.Category) error
	Delete(id uint) error
}

type categoryRepository struct {
	db *gorm.DB
}

func NewCategoryRepository(db *gorm.DB) CategoryRepository {
	return &categoryRepository{db: db}
}

func (r *categoryRepository) Create(cat *models.Category) error {
	return r.db.Create(cat).Error
}

func (r *categoryRepository) FindAll(page, limit int, search string) ([]models.Category, int64, error) {
	var categories []models.Category
	var total int64

	query := r.db.Model(&models.Category{})
	if search != "" {
		query = query.Where("LOWER(nama_category) LIKE ?", "%"+strings.ToLower(search)+"%")
	}

	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	if page < 1 {
		page = 1
	}
	if limit < 1 {
		limit = 10
	}
	offset := (page - 1) * limit

	if err := query.Offset(offset).Limit(limit).Find(&categories).Error; err != nil {
		return nil, 0, err
	}

	return categories, total, nil
}

func (r *categoryRepository) FindByID(id uint) (*models.Category, error) {
	var cat models.Category
	err := r.db.First(&cat, id).Error
	if err != nil {
		return nil, err
	}
	return &cat, nil
}

func (r *categoryRepository) Update(cat *models.Category) error {
	return r.db.Save(cat).Error
}

func (r *categoryRepository) Delete(id uint) error {
	return r.db.Delete(&models.Category{}, id).Error
}
