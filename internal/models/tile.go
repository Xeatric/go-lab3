package models

import (
	"time"

	"gorm.io/gorm"
)

// Tile - модель тротуарной плитки
type Tile struct {
	ID          uint           `gorm:"primaryKey" json:"id"`
	UserID      uint           `gorm:"not null;index" json:"user_id"`                   // владелец
	Name        string         `gorm:"not null;size:200" json:"name"`                   // Название
	Shape       string         `gorm:"not null;size:50" json:"shape"`                   // Форма (квадрат, прямоугольник, шестигранник)
	Color       string         `gorm:"not null;size:50" json:"color"`                   // Цвет
	Size        string         `gorm:"not null;size:50" json:"size"`                    // Размер (например, "30x30", "40x60")
	Material    string         `gorm:"size:100" json:"material"`                        // Материал (бетон, керамогранит, резина)
	PricePerM2  float64        `gorm:"not null;type:numeric(10,2)" json:"price_per_m2"` // Цена за м²
	Stock       int            `gorm:"default:0" json:"stock"`                          // Количество на складе
	Description string         `gorm:"type:text" json:"description"`                    // Описание
	ImageURL    string         `gorm:"size:500" json:"image_url"`                       // Ссылка на изображение
	CreatedAt   time.Time      `json:"created_at"`
	UpdatedAt   time.Time      `json:"updated_at"`
	DeletedAt   gorm.DeletedAt `gorm:"index" json:"-"` // Soft delete
}
