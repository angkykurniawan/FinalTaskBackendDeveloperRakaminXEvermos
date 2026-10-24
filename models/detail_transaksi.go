package models

import "time"

type DetailTransaksi struct {
	ID          uint        `gorm:"primaryKey" json:"id"`
	LogProdukID uint        `gorm:"not null;index" json:"log_produk_id"`
	TransaksiID uint        `gorm:"not null;index" json:"transaksi_id"`
	Kuantitas   int         `gorm:"not null" json:"kuantitas"`
	HargaTotal  float64     `gorm:"type:decimal(15,2);not null" json:"harga_total"`
	CreatedAt   time.Time   `json:"created_at"`
	UpdatedAt   time.Time   `json:"updated_at"`

	LogProduk   *LogProduk  `json:"log_produk,omitempty" gorm:"foreignKey:LogProdukID"`
	Transaksi   *Transaksi  `json:"transaksi,omitempty" gorm:"foreignKey:TransaksiID"`
}
