package models

import "time"

type LogProduk struct {
	ID            uint      `gorm:"primaryKey" json:"id"`
	ProdukID      uint      `gorm:"index" json:"produk_id"`
	NamaProduk    string    `gorm:"type:varchar(255);not null" json:"nama_produk"`
	Slug          string    `gorm:"type:varchar(255)" json:"slug"`
	HargaReseller float64   `gorm:"type:decimal(15,2);default:0" json:"harga_reseller"`
	HargaKonsumen float64   `gorm:"type:decimal(15,2);default:0" json:"harga_konsumen"`
	Deskripsi     string    `gorm:"type:text" json:"deskripsi"`
	TokoID        uint      `gorm:"index" json:"toko_id"`
	CategoryID    uint      `gorm:"index" json:"category_id"`
	CreatedAt     time.Time `json:"created_at"`
	UpdatedAt     time.Time `json:"updated_at"`

	Produk        *Produk   `gorm:"foreignKey:ProdukID" json:"produk,omitempty"`
	Toko          *Toko     `gorm:"foreignKey:TokoID" json:"toko,omitempty"`
	Category      *Category `gorm:"foreignKey:CategoryID" json:"category,omitempty"`
}
