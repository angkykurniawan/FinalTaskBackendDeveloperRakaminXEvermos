package config

import (
	"evermos-mini/models"
	"log"

	"gorm.io/gorm"
)

func SeedData(db *gorm.DB) {
	var count int64
	db.Model(&models.Category{}).Count(&count)
	if count == 0 {
		categories := []models.Category{
			{NamaCategory: "Fashion"},
			{NamaCategory: "Aksesoris"},
			{NamaCategory: "Rumah Tangga"},
		}
		if err := db.Create(&categories).Error; err != nil {
			log.Printf("Gagal seeding data kategori awal: %v\n", err)
		} else {
			log.Println("Seeder kategori awal berhasil dieksekusi.")
		}
	}
}
