package usecases

import (
	"errors"
	"evermos-mini/models"
	"evermos-mini/repositories"
	"strings"
)

type CategoryUsecase interface {
	GetAllCategory(page, limit int, search string) ([]models.Category, int64, error)
	GetCategoryByID(id uint) (*models.Category, error)
	CreateCategory(nama string) (*models.Category, error)
	UpdateCategory(id uint, nama string) (*models.Category, error)
	DeleteCategory(id uint) error
}

type categoryUsecase struct {
	categoryRepo repositories.CategoryRepository
}

func NewCategoryUsecase(categoryRepo repositories.CategoryRepository) CategoryUsecase {
	return &categoryUsecase{categoryRepo: categoryRepo}
}

func (u *categoryUsecase) GetAllCategory(page, limit int, search string) ([]models.Category, int64, error) {
	return u.categoryRepo.FindAll(page, limit, search)
}

func (u *categoryUsecase) GetCategoryByID(id uint) (*models.Category, error) {
	cat, err := u.categoryRepo.FindByID(id)
	if err != nil {
		return nil, errors.New("kategori tidak ditemukan")
	}
	return cat, nil
}

func (u *categoryUsecase) CreateCategory(nama string) (*models.Category, error) {
	nama = strings.TrimSpace(nama)
	if nama == "" {
		return nil, errors.New("nama kategori wajib diisi")
	}

	cat := models.Category{NamaCategory: nama}
	if err := u.categoryRepo.Create(&cat); err != nil {
		return nil, errors.New("gagal menambahkan kategori")
	}
	return &cat, nil
}

func (u *categoryUsecase) UpdateCategory(id uint, nama string) (*models.Category, error) {
	cat, err := u.categoryRepo.FindByID(id)
	if err != nil {
		return nil, errors.New("kategori tidak ditemukan")
	}

	nama = strings.TrimSpace(nama)
	if nama == "" {
		return nil, errors.New("nama kategori wajib diisi")
	}

	cat.NamaCategory = nama
	if err := u.categoryRepo.Update(cat); err != nil {
		return nil, errors.New("gagal memperbarui kategori")
	}
	return cat, nil
}

func (u *categoryUsecase) DeleteCategory(id uint) error {
	_, err := u.categoryRepo.FindByID(id)
	if err != nil {
		return errors.New("kategori tidak ditemukan")
	}

	if err := u.categoryRepo.Delete(id); err != nil {
		return errors.New("gagal menghapus kategori")
	}
	return nil
}
