package usecases

import (
	"errors"
	"evermos-mini/models"
	"evermos-mini/repositories"
	"strings"
)

type AlamatInput struct {
	JudulAlamat  string `json:"judul_alamat"`
	NamaPenerima string `json:"nama_penerima"`
	NoTelp       string `json:"no_telp"`
	Provinsi     string `json:"provinsi"`
	Kota         string `json:"kota"`
	Kecamatan    string `json:"kecamatan"`
	Kelurahan    string `json:"kelurahan"`
	DetailAlamat string `json:"detail_alamat"`
}

type AlamatUsecase interface {
	GetAllAlamat(userID uint) ([]models.Alamat, error)
	GetAlamatByID(id, userID uint) (*models.Alamat, error)
	CreateAlamat(userID uint, input AlamatInput) (*models.Alamat, error)
	UpdateAlamat(id, userID uint, input AlamatInput) (*models.Alamat, error)
	DeleteAlamat(id, userID uint) error
}

type alamatUsecase struct {
	alamatRepo repositories.AlamatRepository
}

func NewAlamatUsecase(alamatRepo repositories.AlamatRepository) AlamatUsecase {
	return &alamatUsecase{alamatRepo: alamatRepo}
}

func (u *alamatUsecase) GetAllAlamat(userID uint) ([]models.Alamat, error) {
	return u.alamatRepo.FindAllByUserID(userID)
}

func (u *alamatUsecase) GetAlamatByID(id, userID uint) (*models.Alamat, error) {
	alamat, err := u.alamatRepo.FindByIDAndUserID(id, userID)
	if err != nil {
		return nil, errors.New("alamat tidak ditemukan atau Anda tidak memiliki akses")
	}
	return alamat, nil
}

func (u *alamatUsecase) CreateAlamat(userID uint, input AlamatInput) (*models.Alamat, error) {
	if strings.TrimSpace(input.NamaPenerima) == "" || strings.TrimSpace(input.DetailAlamat) == "" {
		return nil, errors.New("nama penerima dan detail alamat wajib diisi")
	}

	alamat := models.Alamat{
		JudulAlamat:  strings.TrimSpace(input.JudulAlamat),
		NamaPenerima: strings.TrimSpace(input.NamaPenerima),
		NoTelp:       strings.TrimSpace(input.NoTelp),
		Provinsi:     strings.TrimSpace(input.Provinsi),
		Kota:         strings.TrimSpace(input.Kota),
		Kecamatan:    strings.TrimSpace(input.Kecamatan),
		Kelurahan:    strings.TrimSpace(input.Kelurahan),
		DetailAlamat: strings.TrimSpace(input.DetailAlamat),
		UserID:       userID,
	}

	if err := u.alamatRepo.Create(&alamat); err != nil {
		return nil, errors.New("gagal menyimpan data alamat")
	}

	return &alamat, nil
}

func (u *alamatUsecase) UpdateAlamat(id, userID uint, input AlamatInput) (*models.Alamat, error) {
	alamat, err := u.alamatRepo.FindByIDAndUserID(id, userID)
	if err != nil {
		return nil, errors.New("alamat tidak ditemukan atau Anda tidak memiliki hak akses")
	}

	if strings.TrimSpace(input.JudulAlamat) != "" {
		alamat.JudulAlamat = strings.TrimSpace(input.JudulAlamat)
	}
	if strings.TrimSpace(input.NamaPenerima) != "" {
		alamat.NamaPenerima = strings.TrimSpace(input.NamaPenerima)
	}
	if strings.TrimSpace(input.NoTelp) != "" {
		alamat.NoTelp = strings.TrimSpace(input.NoTelp)
	}
	if strings.TrimSpace(input.Provinsi) != "" {
		alamat.Provinsi = strings.TrimSpace(input.Provinsi)
	}
	if strings.TrimSpace(input.Kota) != "" {
		alamat.Kota = strings.TrimSpace(input.Kota)
	}
	if strings.TrimSpace(input.Kecamatan) != "" {
		alamat.Kecamatan = strings.TrimSpace(input.Kecamatan)
	}
	if strings.TrimSpace(input.Kelurahan) != "" {
		alamat.Kelurahan = strings.TrimSpace(input.Kelurahan)
	}
	if strings.TrimSpace(input.DetailAlamat) != "" {
		alamat.DetailAlamat = strings.TrimSpace(input.DetailAlamat)
	}

	if err := u.alamatRepo.Update(alamat); err != nil {
		return nil, errors.New("gagal memperbarui alamat")
	}

	return alamat, nil
}

func (u *alamatUsecase) DeleteAlamat(id, userID uint) error {
	// Pastikan alamat ada dan milik user
	_, err := u.alamatRepo.FindByIDAndUserID(id, userID)
	if err != nil {
		return errors.New("alamat tidak ditemukan atau Anda tidak berhak menghapus")
	}

	if err := u.alamatRepo.Delete(id, userID); err != nil {
		return errors.New("gagal menghapus alamat")
	}
	return nil
}
