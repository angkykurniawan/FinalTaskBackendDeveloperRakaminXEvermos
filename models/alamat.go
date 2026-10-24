package models

import "time"

type Alamat struct {
	ID           uint      `gorm:"primaryKey" json:"id"`
	JudulAlamat  string    `gorm:"type:varchar(255)" json:"judul_alamat"`
	NamaPenerima string    `gorm:"type:varchar(255);not null" json:"nama_penerima"`
	NoTelp       string    `gorm:"type:varchar(50);not null" json:"no_telp"`
	Provinsi     string    `gorm:"type:varchar(100)" json:"provinsi"`
	Kota         string    `gorm:"type:varchar(100)" json:"kota"`
	Kecamatan    string    `gorm:"type:varchar(100)" json:"kecamatan"`
	Kelurahan    string    `gorm:"type:varchar(100)" json:"kelurahan"`
	DetailAlamat string    `gorm:"type:text;not null" json:"detail_alamat"`
	UserID       uint      `gorm:"not null;index" json:"user_id"`
	CreatedAt    time.Time `json:"created_at"`
	UpdatedAt    time.Time `json:"updated_at"`

	User         *User     `json:"user,omitempty" gorm:"foreignKey:UserID"`
}
