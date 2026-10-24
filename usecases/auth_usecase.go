package usecases

import (
	"errors"
	"evermos-mini/models"
	"evermos-mini/repositories"
	"evermos-mini/utils"
	"strings"
	"time"
)

type RegisterInput struct {
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

type AuthUsecase interface {
	Register(input RegisterInput) (*models.User, error)
	Login(noTelp, kataSandi string) (string, *models.User, error)
}

type authUsecase struct {
	userRepo repositories.UserRepository
	tokoRepo repositories.TokoRepository
}

func NewAuthUsecase(userRepo repositories.UserRepository, tokoRepo repositories.TokoRepository) AuthUsecase {
	return &authUsecase{
		userRepo: userRepo,
		tokoRepo: tokoRepo,
	}
}

func (u *authUsecase) Register(input RegisterInput) (*models.User, error) {
	input.Nama = strings.TrimSpace(input.Nama)
	input.Email = strings.TrimSpace(input.Email)
	input.NoTelp = strings.TrimSpace(input.NoTelp)
	input.KataSandi = strings.TrimSpace(input.KataSandi)

	if input.Nama == "" || input.Email == "" || input.NoTelp == "" || input.KataSandi == "" {
		return nil, errors.New("nama, email, nomor telepon, dan kata sandi wajib diisi")
	}

	if len(input.KataSandi) < 6 {
		return nil, errors.New("kata sandi minimal 6 karakter")
	}

	// Validasi email dan nomor telepon unik
	existing, _ := u.userRepo.FindByEmailOrPhone(input.Email, input.NoTelp)
	if existing != nil {
		if existing.Email == input.Email {
			return nil, errors.New("email sudah terdaftar, gunakan email lain")
		}
		if existing.NoTelp == input.NoTelp {
			return nil, errors.New("nomor telepon sudah terdaftar, gunakan nomor lain")
		}
	}

	hash, err := utils.HashPassword(input.KataSandi)
	if err != nil {
		return nil, errors.New("gagal memproses enkripsi kata sandi")
	}

	var tgl time.Time
	if input.TanggalLahir != "" {
		tgl, err = time.Parse("02/01/2006", input.TanggalLahir)
		if err != nil {
			tgl, err = time.Parse("2006-01-02", input.TanggalLahir)
			if err != nil {
				tgl = time.Time{}
			}
		}
	}

	user := models.User{
		NamaUser:     input.Nama,
		Email:        input.Email,
		NoTelp:       input.NoTelp,
		KataSandi:    hash,
		TanggalLahir: tgl,
		JenisKelamin: input.JenisKelamin,
		Pekerjaan:    input.Pekerjaan,
		IDProvinsi:   input.IDProvinsi,
		IDKota:       input.IDKota,
		Role:         "user",
	}

	if err := u.userRepo.Create(&user); err != nil {
		return nil, errors.New("gagal menyimpan data pengguna baru")
	}

	// Toko otomatis terbuat ketika user mendaftar
	namaToko := "Toko " + input.Nama
	toko := models.Toko{
		NamaToko: namaToko,
		UserID:   user.ID,
		Photo:    "",
	}
	_ = u.tokoRepo.Create(&toko)

	user.Toko = &toko
	return &user, nil
}

func (u *authUsecase) Login(noTelp, kataSandi string) (string, *models.User, error) {
	noTelp = strings.TrimSpace(noTelp)
	kataSandi = strings.TrimSpace(kataSandi)

	if noTelp == "" || kataSandi == "" {
		return "", nil, errors.New("nomor telepon dan kata sandi wajib diisi")
	}

	user, err := u.userRepo.FindByPhone(noTelp)
	if err != nil {
		return "", nil, errors.New("nomor telepon tidak ditemukan")
	}

	if !utils.CheckPasswordHash(kataSandi, user.KataSandi) {
		return "", nil, errors.New("kata sandi tidak sesuai")
	}

	token, err := utils.GenerateJWT(user.ID, user.Role)
	if err != nil {
		return "", nil, errors.New("gagal membuat token otorisasi")
	}

	return token, user, nil
}
