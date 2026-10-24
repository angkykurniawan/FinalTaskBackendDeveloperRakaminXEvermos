package models

import "time"

type User struct {
	ID           uint        `gorm:"primaryKey" json:"id"`
	NamaUser     string      `gorm:"type:varchar(255);not null" json:"nama_user"`
	KataSandi    string      `gorm:"type:varchar(255);not null" json:"-"`
	NoTelp       string      `gorm:"type:varchar(50);uniqueIndex;not null" json:"no_telp"`
	TanggalLahir time.Time   `json:"tanggal_lahir"`
	JenisKelamin string      `gorm:"type:varchar(20)" json:"jenis_kelamin"`
	Email        string      `gorm:"type:varchar(255);uniqueIndex;not null" json:"email"`
	Pekerjaan    string      `gorm:"type:varchar(100)" json:"pekerjaan"`
	IDProvinsi   string      `gorm:"type:varchar(50)" json:"id_provinsi"`
	IDKota       string      `gorm:"type:varchar(50)" json:"id_kota"`
	Role         string      `gorm:"type:varchar(50);default:'user'" json:"role"`
	CreatedAt    time.Time   `json:"created_at"`
	UpdatedAt    time.Time   `json:"updated_at"`

	Toko         *Toko       `json:"toko,omitempty" gorm:"foreignKey:UserID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;"`

	Alamat       []Alamat    `json:"alamat,omitempty" gorm:"foreignKey:UserID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;"`
	Transaksi    []Transaksi `json:"transaksi,omitempty" gorm:"foreignKey:UserID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;"`
}
