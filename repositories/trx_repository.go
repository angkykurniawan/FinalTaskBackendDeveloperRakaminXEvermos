package repositories

import (
	"evermos-mini/models"
	"time"

	"gorm.io/gorm"
)

type TrxRepository interface {
	Create(trx *models.Transaksi, details []models.DetailTransaksi) error
	FindAllByUserID(userID uint) ([]models.Transaksi, error)
	FindByIDAndUserID(id, userID uint) (*models.Transaksi, error)
	FindByID(id uint) (*models.Transaksi, error)
	UpdateStatus(id uint, status string) (*models.Transaksi, error)
	GetDB() *gorm.DB
}

type trxRepository struct {
	db *gorm.DB
}

func NewTrxRepository(db *gorm.DB) TrxRepository {
	return &trxRepository{db: db}
}

func (r *trxRepository) GetDB() *gorm.DB {
	return r.db
}

func (r *trxRepository) Create(trx *models.Transaksi, details []models.DetailTransaksi) error {
	return r.db.Transaction(func(tx *gorm.DB) error {
		if err := tx.Create(trx).Error; err != nil {
			return err
		}

		for i := range details {
			details[i].TransaksiID = trx.ID
		}

		if len(details) > 0 {
			if err := tx.Create(&details).Error; err != nil {
				return err
			}
		}

		return nil
	})
}

func (r *trxRepository) FindAllByUserID(userID uint) ([]models.Transaksi, error) {
	var list []models.Transaksi
	err := r.db.Where("user_id = ?", userID).
		Preload("DetailTransaksi.LogProduk.Produk").
		Preload("DetailTransaksi.LogProduk.Toko").
		Preload("DetailTransaksi.LogProduk.Category").
		Find(&list).Error
	return list, err
}

func (r *trxRepository) FindByIDAndUserID(id, userID uint) (*models.Transaksi, error) {
	var trx models.Transaksi
	err := r.db.Where("id = ? AND user_id = ?", id, userID).
		Preload("DetailTransaksi.LogProduk.Produk").
		Preload("DetailTransaksi.LogProduk.Toko").
		Preload("DetailTransaksi.LogProduk.Category").
		Preload("User").
		First(&trx).Error
	if err != nil {
		return nil, err
	}
	return &trx, nil
}

func (r *trxRepository) FindByID(id uint) (*models.Transaksi, error) {
	var trx models.Transaksi
	err := r.db.Preload("DetailTransaksi.LogProduk.Produk").
		Preload("DetailTransaksi.LogProduk.Toko").
		Preload("DetailTransaksi.LogProduk.Category").
		Preload("User").
		First(&trx, id).Error
	if err != nil {
		return nil, err
	}
	return &trx, nil
}

func (r *trxRepository) UpdateStatus(id uint, status string) (*models.Transaksi, error) {
	var trx models.Transaksi
	if err := r.db.First(&trx, id).Error; err != nil {
		return nil, err
	}

	trx.Status = status
	trx.UpdatedAt = time.Now()
	if err := r.db.Save(&trx).Error; err != nil {
		return nil, err
	}
	return &trx, nil
}
