package service

import (
	"errors"
	"paving-tiles-api/internal/dto"
	"paving-tiles-api/internal/models"
	"paving-tiles-api/internal/repository"

	"gorm.io/gorm"
)

type TileService interface {
	CreateTile(req *dto.CreateTileRequest) (*dto.TileResponse, error)
	GetTileByID(id uint) (*dto.TileResponse, error)
	GetTiles(page, limit int) (*dto.PaginationResponse, error)
	UpdateTile(id uint, req *dto.UpdateTileRequest) (*dto.TileResponse, error)
	PatchTile(id uint, req *dto.PatchTileRequest) (*dto.TileResponse, error)
	DeleteTile(id uint) error
}

type tileService struct {
	repo repository.TileRepository
}

func NewTileService(repo repository.TileRepository) TileService {
	return &tileService{repo: repo}
}

func (s *tileService) CreateTile(req *dto.CreateTileRequest) (*dto.TileResponse, error) {
	tile := &models.Tile{
		Name:        req.Name,
		Shape:       req.Shape,
		Color:       req.Color,
		Size:        req.Size,
		Material:    req.Material,
		PricePerM2:  req.PricePerM2,
		Stock:       req.Stock,
		Description: req.Description,
		ImageURL:    req.ImageURL,
	}

	if err := s.repo.Create(tile); err != nil {
		return nil, err
	}

	return s.toResponse(tile), nil
}

func (s *tileService) GetTileByID(id uint) (*dto.TileResponse, error) {
	tile, err := s.repo.FindByID(id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("tile not found")
		}
		return nil, err
	}

	return s.toResponse(tile), nil
}

func (s *tileService) GetTiles(page, limit int) (*dto.PaginationResponse, error) {
	if page < 1 {
		page = 1
	}
	if limit < 1 || limit > 100 {
		limit = 10
	}

	offset := (page - 1) * limit
	tiles, total, err := s.repo.FindAll(offset, limit)
	if err != nil {
		return nil, err
	}

	responses := make([]dto.TileResponse, len(tiles))
	for i, tile := range tiles {
		responses[i] = *s.toResponse(&tile)
	}

	totalPages := (total + int64(limit) - 1) / int64(limit)

	return &dto.PaginationResponse{
		Data: responses,
		Pagination: dto.Meta{
			Total:      total,
			Page:       page,
			Limit:      limit,
			TotalPages: totalPages,
		},
	}, nil
}

func (s *tileService) UpdateTile(id uint, req *dto.UpdateTileRequest) (*dto.TileResponse, error) {
	tile, err := s.repo.FindByID(id)
	if err != nil {
		return nil, err
	}

	tile.Name = req.Name
	tile.Shape = req.Shape
	tile.Color = req.Color
	tile.Size = req.Size
	tile.Material = req.Material
	tile.PricePerM2 = req.PricePerM2
	tile.Stock = req.Stock
	tile.Description = req.Description
	tile.ImageURL = req.ImageURL

	if err := s.repo.Update(tile); err != nil {
		return nil, err
	}

	return s.toResponse(tile), nil
}

func (s *tileService) PatchTile(id uint, req *dto.PatchTileRequest) (*dto.TileResponse, error) {
	tile, err := s.repo.FindByID(id)
	if err != nil {
		return nil, err
	}

	if req.Name != nil {
		tile.Name = *req.Name
	}
	if req.Shape != nil {
		tile.Shape = *req.Shape
	}
	if req.Color != nil {
		tile.Color = *req.Color
	}
	if req.Size != nil {
		tile.Size = *req.Size
	}
	if req.Material != nil {
		tile.Material = *req.Material
	}
	if req.PricePerM2 != nil {
		tile.PricePerM2 = *req.PricePerM2
	}
	if req.Stock != nil {
		tile.Stock = *req.Stock
	}
	if req.Description != nil {
		tile.Description = *req.Description
	}
	if req.ImageURL != nil {
		tile.ImageURL = *req.ImageURL
	}

	if err := s.repo.Update(tile); err != nil {
		return nil, err
	}

	return s.toResponse(tile), nil
}

func (s *tileService) DeleteTile(id uint) error {
	// Проверяем существование плитки
	_, err := s.repo.FindByID(id)
	if err != nil {
		return err
	}

	// Soft delete
	return s.repo.Delete(id)
}

func (s *tileService) toResponse(tile *models.Tile) *dto.TileResponse {
	return &dto.TileResponse{
		ID:          tile.ID,
		Name:        tile.Name,
		Shape:       tile.Shape,
		Color:       tile.Color,
		Size:        tile.Size,
		Material:    tile.Material,
		PricePerM2:  tile.PricePerM2,
		Stock:       tile.Stock,
		Description: tile.Description,
		ImageURL:    tile.ImageURL,
		CreatedAt:   tile.CreatedAt,
		UpdatedAt:   tile.UpdatedAt,
	}
}
