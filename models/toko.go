package models

import "time"

type Toko struct {
	ID        uint      `gorm:"primaryKey" json:"id"`
	NamaToko  string    `gorm:"type:varchar(255);not null" json:"nama_toko"`
	Photo     string    `gorm:"type:text" json:"photo"`
	UserID    uint      `gorm:"uniqueIndex;not null" json:"user_id"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`

	User      *User     `json:"user,omitempty" gorm:"foreignKey:UserID"`
}
