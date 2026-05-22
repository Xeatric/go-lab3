package models

import (
	"time"

	"gorm.io/gorm"
)

type User struct {
	ID           uint           `gorm:"primaryKey" json:"id"`
	Email        string         `gorm:"uniqueIndex;size:255;not null" json:"email"`
	PasswordHash string         `gorm:"not null" json:"-"`
	Salt         string         `gorm:"not null" json:"-"`
	YandexID     string         `gorm:"size:255" json:"-"`
	VkID         string         `gorm:"size:255" json:"-"`
	Name         string         `gorm:"size:255" json:"name"`
	Role         string         `gorm:"default:user;size:50" json:"role"`
	IsActive     bool           `gorm:"default:true" json:"is_active"`
	LastLoginAt  *time.Time     `json:"last_login_at,omitempty"`
	CreatedAt    time.Time      `json:"created_at"`
	UpdatedAt    time.Time      `json:"updated_at"`
	DeletedAt    gorm.DeletedAt `gorm:"index" json:"-"`
}

// UserResponse DTO для безопасного ответа
type UserResponse struct {
	ID        uint      `json:"id"`
	Email     string    `json:"email"`
	Name      string    `json:"name"`
	Role      string    `json:"role"`
	CreatedAt time.Time `json:"created_at"`
}
