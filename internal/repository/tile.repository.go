package repository

import (
	"paving-tiles-api/internal/models"

	"gorm.io/gorm"
)

type TileRepository interface {
	Create(tile *models.Tile) error
	FindByID(id uint) (*models.Tile, error)
	FindAll(offset, limit int) ([]models.Tile, int64, error)
	Update(tile *models.Tile) error
	Delete(id uint) error
}

type tileRepository struct {
	db *gorm.DB
}

func NewTileRepository(db *gorm.DB) TileRepository {
	return &tileRepository{db: db}
}

func (r *tileRepository) Create(tile *models.Tile) error {
	return r.db.Create(tile).Error
}

func (r *tileRepository) FindByID(id uint) (*models.Tile, error) {
	var tile models.Tile
	// GORM автоматически исключает soft deleted записи
	err := r.db.First(&tile, id).Error
	if err != nil {
		return nil, err
	}
	return &tile, nil
}

func (r *tileRepository) FindAll(offset, limit int) ([]models.Tile, int64, error) {
	var tiles []models.Tile
	var total int64

	// Подсчет общего количества (исключая soft deleted)
	if err := r.db.Model(&models.Tile{}).Count(&total).Error; err != nil {
		return nil, 0, err
	}

	// Получение пагинированных данных
	err := r.db.Offset(offset).Limit(limit).Find(&tiles).Error
	if err != nil {
		return nil, 0, err
	}

	return tiles, total, nil
}

func (r *tileRepository) Update(tile *models.Tile) error {
	return r.db.Save(tile).Error
}

func (r *tileRepository) Delete(id uint) error {
	// Soft delete - GORM автоматически устанавливает DeletedAt
	return r.db.Delete(&models.Tile{}, id).Error
}
