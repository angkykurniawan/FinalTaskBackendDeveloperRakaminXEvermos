package repositories

import (
	"evermos-mini/models"

	"gorm.io/gorm"
)

type AlamatRepository interface {
	Create(alamat *models.Alamat) error
	FindAllByUserID(userID uint) ([]models.Alamat, error)
	FindByIDAndUserID(id, userID uint) (*models.Alamat, error)
	Update(alamat *models.Alamat) error
	Delete(id, userID uint) error
}

type alamatRepository struct {
	db *gorm.DB
}

func NewAlamatRepository(db *gorm.DB) AlamatRepository {
	return &alamatRepository{db: db}
}

func (r *alamatRepository) Create(alamat *models.Alamat) error {
	return r.db.Create(alamat).Error
}

func (r *alamatRepository) FindAllByUserID(userID uint) ([]models.Alamat, error) {
	var list []models.Alamat
	err := r.db.Where("user_id = ?", userID).Find(&list).Error
	return list, err
}

func (r *alamatRepository) FindByIDAndUserID(id, userID uint) (*models.Alamat, error) {
	var alamat models.Alamat
	err := r.db.Where("id = ? AND user_id = ?", id, userID).First(&alamat).Error
	if err != nil {
		return nil, err
	}
	return &alamat, nil
}

func (r *alamatRepository) Update(alamat *models.Alamat) error {
	return r.db.Save(alamat).Error
}

func (r *alamatRepository) Delete(id, userID uint) error {
	return r.db.Where("id = ? AND user_id = ?", id, userID).Delete(&models.Alamat{}).Error
}
