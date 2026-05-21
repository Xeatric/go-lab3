package controller

import (
	"errors"
	"net/http"
	"paving-tiles-api/internal/dto"
	"paving-tiles-api/internal/service"
	"strconv"

	"github.com/gin-gonic/gin"
)

type TileController struct {
	service service.TileService
}

func NewTileController(service service.TileService) *TileController {
	return &TileController{service: service}
}

// GetTiles возвращает список плитки с пагинацией
// @GET /api/v1/tiles?page=1&limit=10
func (c *TileController) GetTiles(ctx *gin.Context) {
	var paginationReq dto.PaginationRequest
	if err := ctx.ShouldBindQuery(&paginationReq); err != nil {
		ctx.Error(err)
		return
	}

	response, err := c.service.GetTiles(paginationReq.Page, paginationReq.Limit)
	if err != nil {
		ctx.Error(err)
		return
	}

	ctx.JSON(http.StatusOK, response)
}

// GetTileByID возвращает плитку по ID
// @GET /api/v1/tiles/:id
func (c *TileController) GetTileByID(ctx *gin.Context) {
	id, err := parseID(ctx.Param("id"))
	if err != nil {
		ctx.Error(err)
		return
	}

	tile, err := c.service.GetTileByID(id)
	if err != nil {
		ctx.Error(err)
		return
	}

	ctx.JSON(http.StatusOK, tile)
}

// CreateTile создает новую плитку
// @POST /api/v1/tiles
func (c *TileController) CreateTile(ctx *gin.Context) {
	var req dto.CreateTileRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.Error(err)
		return
	}

	tile, err := c.service.CreateTile(&req)
	if err != nil {
		ctx.Error(err)
		return
	}

	ctx.JSON(http.StatusCreated, tile)
}

// UpdateTile полностью обновляет плитку
// @PUT /api/v1/tiles/:id
func (c *TileController) UpdateTile(ctx *gin.Context) {
	id, err := parseID(ctx.Param("id"))
	if err != nil {
		ctx.Error(err)
		return
	}

	var req dto.UpdateTileRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.Error(err)
		return
	}

	tile, err := c.service.UpdateTile(id, &req)
	if err != nil {
		ctx.Error(err)
		return
	}

	ctx.JSON(http.StatusOK, tile)
}

// PatchTile частично обновляет плитку
// @PATCH /api/v1/tiles/:id
func (c *TileController) PatchTile(ctx *gin.Context) {
	id, err := parseID(ctx.Param("id"))
	if err != nil {
		ctx.Error(err)
		return
	}

	var req dto.PatchTileRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.Error(err)
		return
	}

	tile, err := c.service.PatchTile(id, &req)
	if err != nil {
		ctx.Error(err)
		return
	}

	ctx.JSON(http.StatusOK, tile)
}

// DeleteTile мягко удаляет плитку
// @DELETE /api/v1/tiles/:id
func (c *TileController) DeleteTile(ctx *gin.Context) {
	id, err := parseID(ctx.Param("id"))
	if err != nil {
		ctx.Error(err)
		return
	}

	if err := c.service.DeleteTile(id); err != nil {
		ctx.Error(err)
		return
	}

	ctx.Status(http.StatusNoContent)
}

func parseID(idStr string) (uint, error) {
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		return 0, errors.New("invalid id format")
	}
	return uint(id), nil
}
