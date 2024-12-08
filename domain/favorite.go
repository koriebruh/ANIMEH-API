package domain

import "time"

type Favorite struct {
	ID      uint      `gorm:"primaryKey"`
	UserID  uint      `gorm:"not null"`
	AnimeID string    `gorm:"not null"`
	AddedAt time.Time `gorm:"autoCreateTime"`
	User    User      `gorm:"foreignKey:UserID;constraint:OnDelete:CASCADE"`
}
