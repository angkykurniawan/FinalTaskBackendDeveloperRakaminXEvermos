package usecases

import (
	"errors"
	"evermos-mini/models"
	"evermos-mini/repositories"
	"evermos-mini/utils"
	"strings"
	"time"
)

type UpdateUserInput struct {
	Nama         string `json:"nama"`
	Email        string `json:"email"`
	NoTelp       string `json:"no_telp"`
	KataSandi    string `json:"kata_sandi"`
	TanggalLahir string `json:"tanggal_Lahir"`
	JenisKelamin string `json:"jenis_kelamin"`
	Pekerjaan    string `json:"pekerjaan"`
	IDProvinsi   string `json:"id_provinsi"`
	IDKota       string `json:"id_kota"`
}

type UserUsecase interface {
	GetProfile(userID uint) (*models.User, error)
	UpdateProfile(userID uint, input UpdateUserInput) (*models.User, error)
}

type userUsecase struct {
	userRepo repositories.UserRepository
}

func NewUserUsecase(userRepo repositories.UserRepository) UserUsecase {
	return &userUsecase{userRepo: userRepo}
}

func (u *userUsecase) GetProfile(userID uint) (*models.User, error) {
	user, err := u.userRepo.FindProfileByID(userID)
	if err != nil {
		return nil, errors.New("data pengguna tidak ditemukan")
	}
	return user, nil
}

func (u *userUsecase) UpdateProfile(userID uint, input UpdateUserInput) (*models.User, error) {
	user, err := u.userRepo.FindByID(userID)
	if err != nil {
		return nil, errors.New("pengguna tidak ditemukan")
	}

	input.Email = strings.TrimSpace(input.Email)
	input.NoTelp = strings.TrimSpace(input.NoTelp)

	// Validasi jika email atau no_telp diubah, pastikan tidak bertabrakan dengan akun lain
	if input.Email != "" && input.Email != user.Email {
		existing, _ := u.userRepo.FindByEmailOrPhone(input.Email, "")
		if existing != nil && existing.ID != userID {
			return nil, errors.New("email sudah digunakan oleh pengguna lain")
		}
		user.Email = input.Email
	}

	if input.NoTelp != "" && input.NoTelp != user.NoTelp {
		existing, _ := u.userRepo.FindByEmailOrPhone("", input.NoTelp)
		if existing != nil && existing.ID != userID {
			return nil, errors.New("nomor telepon sudah digunakan oleh pengguna lain")
		}
		user.NoTelp = input.NoTelp
	}

	if strings.TrimSpace(input.Nama) != "" {
		user.NamaUser = strings.TrimSpace(input.Nama)
	}

	if strings.TrimSpace(input.KataSandi) != "" {
		if len(input.KataSandi) < 6 {
			return nil, errors.New("kata sandi baru minimal 6 karakter")
		}
		hash, err := utils.HashPassword(input.KataSandi)
		if err == nil {
			user.KataSandi = hash
		}
	}

	if input.JenisKelamin != "" {
		user.JenisKelamin = input.JenisKelamin
	}
	if input.Pekerjaan != "" {
		user.Pekerjaan = input.Pekerjaan
	}
	if input.IDProvinsi != "" {
		user.IDProvinsi = input.IDProvinsi
	}
	if input.IDKota != "" {
		user.IDKota = input.IDKota
	}

	if input.TanggalLahir != "" {
		tgl, err := time.Parse("02/01/2006", input.TanggalLahir)
		if err != nil {
			tgl, err = time.Parse("2006-01-02", input.TanggalLahir)
		}
		if err == nil {
			user.TanggalLahir = tgl
		}
	}

	if err := u.userRepo.Update(user); err != nil {
		return nil, errors.New("gagal memperbarui profil pengguna")
	}

	return user, nil
}
