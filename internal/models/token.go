package models

import (
	"time"

	"gorm.io/gorm"
)

type Token struct {
	ID        uint           `gorm:"primaryKey" json:"id"`
	UserID    uint           `gorm:"not null;index" json:"user_id"`
	User      User           `gorm:"foreignKey:UserID" json:"-"`
	TokenHash string         `gorm:"uniqueIndex;size:255;not null" json:"-"`
	TokenType string         `gorm:"size:20;not null" json:"-"`
	ExpiresAt time.Time      `gorm:"not null;index" json:"-"`
	IsRevoked bool           `gorm:"default:false;index" json:"-"`
	RevokedAt *time.Time     `json:"-"`
	UserAgent string         `gorm:"size:500" json:"-"`
	IPAddress string         `gorm:"size:45" json:"-"`
	CreatedAt time.Time      `json:"-"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"-"`
}
