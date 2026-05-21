package dto

import "time"

// CreateTileRequest DTO для создания плитки
type CreateTileRequest struct {
	Name        string  `json:"name" binding:"required,min=2,max=200"`
	Shape       string  `json:"shape" binding:"required,oneof=square rectangle hexagon circle"`
	Color       string  `json:"color" binding:"required,min=2,max=50"`
	Size        string  `json:"size" binding:"required,min=2,max=50"`
	Material    string  `json:"material" binding:"omitempty,max=100"`
	PricePerM2  float64 `json:"price_per_m2" binding:"required,min=0"`
	Stock       int     `json:"stock" binding:"min=0"`
	Description string  `json:"description" binding:"max=1000"`
	ImageURL    string  `json:"image_url" binding:"omitempty,url,max=500"`
}

// UpdateTileRequest DTO для полного обновления
type UpdateTileRequest struct {
	Name        string  `json:"name" binding:"required,min=2,max=200"`
	Shape       string  `json:"shape" binding:"required,oneof=square rectangle hexagon circle"`
	Color       string  `json:"color" binding:"required,min=2,max=50"`
	Size        string  `json:"size" binding:"required,min=2,max=50"`
	Material    string  `json:"material" binding:"omitempty,max=100"`
	PricePerM2  float64 `json:"price_per_m2" binding:"required,min=0"`
	Stock       int     `json:"stock" binding:"min=0"`
	Description string  `json:"description" binding:"max=1000"`
	ImageURL    string  `json:"image_url" binding:"omitempty,url,max=500"`
}

// PatchTileRequest DTO для частичного обновления
type PatchTileRequest struct {
	Name        *string  `json:"name" binding:"omitempty,min=2,max=200"`
	Shape       *string  `json:"shape" binding:"omitempty,oneof=square rectangle hexagon circle"`
	Color       *string  `json:"color" binding:"omitempty,min=2,max=50"`
	Size        *string  `json:"size" binding:"omitempty,min=2,max=50"`
	Material    *string  `json:"material" binding:"omitempty,max=100"`
	PricePerM2  *float64 `json:"price_per_m2" binding:"omitempty,min=0"`
	Stock       *int     `json:"stock" binding:"omitempty,min=0"`
	Description *string  `json:"description" binding:"omitempty,max=1000"`
	ImageURL    *string  `json:"image_url" binding:"omitempty,url,max=500"`
}

// TileResponse DTO для ответа
type TileResponse struct {
	ID          uint      `json:"id"`
	Name        string    `json:"name"`
	Shape       string    `json:"shape"`
	Color       string    `json:"color"`
	Size        string    `json:"size"`
	Material    string    `json:"material"`
	PricePerM2  float64   `json:"price_per_m2"`
	Stock       int       `json:"stock"`
	Description string    `json:"description"`
	ImageURL    string    `json:"image_url"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

// PaginationRequest DTO для пагинации
type PaginationRequest struct {
	Page  int `form:"page" binding:"omitempty,min=1"`
	Limit int `form:"limit" binding:"omitempty,min=1,max=100"`
}

// PaginationResponse DTO для ответа с пагинацией
type PaginationResponse struct {
	Data       []TileResponse `json:"data"`
	Pagination Meta           `json:"meta"`
}

type Meta struct {
	Total      int64 `json:"total"`
	Page       int   `json:"page"`
	Limit      int   `json:"limit"`
	TotalPages int64 `json:"total_pages"`
}
