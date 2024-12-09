package domain

import (
	"time"
)

type User struct {
	ID       uint   `gorm:"primaryKey"`
	Username string `gorm:"unique;not null" json:"username"`
	Email    string `gorm:"unique;not null" json:"email"`
	Password string `gorm:"not null" json:"password"`

	//INI HARUS NYA FIELD NYA DI PISAH
	Token     string `gorm:"unique"`
	NewPass   string
	CreatedAt time.Time `gorm:"autoCreateTime"`
}
