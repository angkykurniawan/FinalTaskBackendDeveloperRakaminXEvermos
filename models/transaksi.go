package models

import "time"

type Transaksi struct {
	ID               uint              `gorm:"primaryKey" json:"id"`
	AlamatPengiriman uint              `gorm:"not null" json:"alamat_pengiriman"`
	HargaTotal       float64           `gorm:"type:decimal(15,2);not null" json:"harga_total"`
	Invoice          string            `gorm:"type:varchar(100);uniqueIndex;not null" json:"invoice"`
	MetodeBayar      string            `gorm:"type:varchar(50);not null" json:"metode_bayar"`
	Status           string            `gorm:"type:varchar(50);default:'Menunggu Pembayaran'" json:"status"`
	UserID           uint              `gorm:"not null;index" json:"user_id"`
	CreatedAt        time.Time         `json:"created_at"`
	UpdatedAt        time.Time         `json:"updated_at"`

	User             *User             `json:"user,omitempty" gorm:"foreignKey:UserID"`
	DetailTransaksi  []DetailTransaksi `json:"detail_transaksi" gorm:"foreignKey:TransaksiID"`
}
