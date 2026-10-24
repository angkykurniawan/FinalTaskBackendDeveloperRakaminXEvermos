package usecases

import (
	"errors"
	"evermos-mini/models"
	"evermos-mini/repositories"
	"strings"
)

type TokoUsecase interface {
	GetMyToko(userID uint) (*models.Toko, error)
	GetTokoByID(id uint) (*models.Toko, error)
	GetAllToko(page, limit int, nama string) ([]models.Toko, int64, error)
	UpdateToko(userID, id uint, namaToko, photoURL string) (*models.Toko, error)
}

type tokoUsecase struct {
	tokoRepo repositories.TokoRepository
}

func NewTokoUsecase(tokoRepo repositories.TokoRepository) TokoUsecase {
	return &tokoUsecase{tokoRepo: tokoRepo}
}

func (u *tokoUsecase) GetMyToko(userID uint) (*models.Toko, error) {
	toko, err := u.tokoRepo.FindByUserID(userID)
	if err != nil {
		return nil, errors.New("data toko Anda tidak ditemukan")
	}
	return toko, nil
}

func (u *tokoUsecase) GetTokoByID(id uint) (*models.Toko, error) {
	toko, err := u.tokoRepo.FindByID(id)
	if err != nil {
		return nil, errors.New("toko dengan ID tersebut tidak ditemukan")
	}
	return toko, nil
}

func (u *tokoUsecase) GetAllToko(page, limit int, nama string) ([]models.Toko, int64, error) {
	tokos, total, err := u.tokoRepo.FindAll(page, limit, nama)
	if err != nil {
		return nil, 0, errors.New("gagal memuat daftar toko")
	}

	// Fallback seperti pada referensi jika pencarian filter tidak menemukan data
	if len(tokos) == 0 && nama != "" {
		allTokos, allTotal, errAll := u.tokoRepo.FindAll(page, limit, "")
		if errAll == nil && len(allTokos) > 0 {
			return allTokos, allTotal, nil
		}
	}

	return tokos, total, nil
}

func (u *tokoUsecase) UpdateToko(userID, id uint, namaToko, photoURL string) (*models.Toko, error) {
	toko, err := u.tokoRepo.FindByID(id)
	if err != nil {
		return nil, errors.New("toko tidak ditemukan")
	}

	// Proteksi: User tidak dapat mengelola data toko user lain
	if toko.UserID != userID {
		return nil, errors.New("akses ditolak: Anda tidak memiliki wewenang untuk mengubah toko ini")
	}

	if strings.TrimSpace(namaToko) != "" {
		toko.NamaToko = strings.TrimSpace(namaToko)
	}

	if photoURL != "" {
		toko.Photo = photoURL
	}

	if err := u.tokoRepo.Update(toko); err != nil {
		return nil, errors.New("gagal memperbarui data toko")
	}

	return toko, nil
}
