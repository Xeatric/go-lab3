package controller

import (
	"net/http"
	"strconv"

	"paving-tiles-api/internal/auth/middleware"
	"paving-tiles-api/internal/dto"
	"paving-tiles-api/internal/service"

	"github.com/gin-gonic/gin"
)

type TileController struct {
	service service.TileService
}

func NewTileController(service service.TileService) *TileController {
	return &TileController{service: service}
}

// GetTiles - получение списка (только своих плиток)
func (c *TileController) GetTiles(ctx *gin.Context) {
	userID := middleware.GetCurrentUserID(ctx)
	if userID == 0 {
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}

	var paginationReq dto.PaginationRequest
	if err := ctx.ShouldBindQuery(&paginationReq); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	response, err := c.service.GetTiles(userID, paginationReq.Page, paginationReq.Limit)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, response)
}

// GetTileByID - получение по ID (только своей)
func (c *TileController) GetTileByID(ctx *gin.Context) {
	userID := middleware.GetCurrentUserID(ctx)
	if userID == 0 {
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}

	id, err := strconv.ParseUint(ctx.Param("id"), 10, 32)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "invalid id"})
		return
	}

	tile, err := c.service.GetTileByID(userID, uint(id))
	if err != nil {
		ctx.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, tile)
}

// CreateTile - создание (привязывается к текущему пользователю)
func (c *TileController) CreateTile(ctx *gin.Context) {
	userID := middleware.GetCurrentUserID(ctx)
	if userID == 0 {
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}

	var req dto.CreateTileRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	tile, err := c.service.CreateTile(userID, &req)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusCreated, tile)
}

// UpdateTile - полное обновление плитки
func (c *TileController) UpdateTile(ctx *gin.Context) {
	userID := middleware.GetCurrentUserID(ctx)
	if userID == 0 {
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}

	id, err := strconv.ParseUint(ctx.Param("id"), 10, 32)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "invalid id"})
		return
	}

	var req dto.UpdateTileRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	tile, err := c.service.UpdateTile(userID, uint(id), &req)
	if err != nil {
		ctx.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, tile)
}

// PatchTile - частичное обновление плитки
func (c *TileController) PatchTile(ctx *gin.Context) {
	userID := middleware.GetCurrentUserID(ctx)
	if userID == 0 {
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}

	id, err := strconv.ParseUint(ctx.Param("id"), 10, 32)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "invalid id"})
		return
	}

	var req dto.PatchTileRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	tile, err := c.service.PatchTile(userID, uint(id), &req)
	if err != nil {
		ctx.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, tile)
}

// DeleteTile - удаление (только своей)
func (c *TileController) DeleteTile(ctx *gin.Context) {
	userID := middleware.GetCurrentUserID(ctx)
	if userID == 0 {
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}

	id, err := strconv.ParseUint(ctx.Param("id"), 10, 32)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "invalid id"})
		return
	}

	if err := c.service.DeleteTile(userID, uint(id)); err != nil {
		ctx.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}

	ctx.Status(http.StatusNoContent)
}
