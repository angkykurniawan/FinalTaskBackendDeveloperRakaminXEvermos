package models

import "time"

type Produk struct {
	ID            uint      `gorm:"primaryKey" json:"id"`
	NamaProduk    string    `gorm:"type:varchar(255);not null" json:"nama_produk"`
	Slug          string    `gorm:"type:varchar(255)" json:"slug"`
	HargaReseller float64   `gorm:"type:decimal(15,2);default:0" json:"harga_reseller"`
	HargaKonsumen float64   `gorm:"type:decimal(15,2);default:0" json:"harga_konsumen"`
	Stok          int       `gorm:"default:0" json:"stok"`
	Deskripsi     string    `gorm:"type:text" json:"deskripsi"`
	Gambar        string    `gorm:"type:text" json:"gambar"`
	TokoID        uint      `gorm:"not null;index" json:"toko_id"`
	CategoryID    uint      `gorm:"not null;index" json:"category_id"`
	CreatedAt     time.Time `json:"created_at"`
	UpdatedAt     time.Time `json:"updated_at"`

	Toko          Toko      `json:"toko" gorm:"foreignKey:TokoID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;"`
	Category      Category  `json:"category" gorm:"foreignKey:CategoryID;constraint:OnUpdate:CASCADE,OnDelete:RESTRICT;"`
}
