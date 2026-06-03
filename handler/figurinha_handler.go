package handler

import (
	"errors"
	"ponderada-3/domain"
	"ponderada-3/service"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

type FigureHandler struct {
	svc service.FigureService
}

func NewFigureHandler(svc service.FigureService) *FigureHandler {
	return &FigureHandler{svc: svc}
}

func (h *FigureHandler) RegisterRoutes(r *gin.Engine) {
	r.POST("/figurinha", h.Create)
	r.GET("/figurinha", h.ListAll)
	r.GET("/figurinha/:id", h.GetByID)
	r.PUT("/figurinha/:id", h.Update)
	r.DELETE("/figurinha/:id", h.Delete)
}

func (h *FigureHandler) Create(c *gin.Context) {
	var req domain.CreateFigureRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	figurinha, err := h.svc.Create(req)

	if err != nil {
		c.JSON(mapErrorToStatus(err), gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, figurinha)
}


func (h *FigureHandler) ListAll(c *gin.Context) {
	tipo := c.Query("tipo")
	posicao := c.Query("posicao")
	req := domain.FindAllFigureRequest{
		Tipo:     domain.FigurinhaType(tipo),
		Posicao:  domain.FigurinhaPosition(posicao),
	}

	figurinhas, err := h.svc.ListAll(req)

	if err != nil {
		c.JSON(mapErrorToStatus(err), gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, figurinhas)
}

func (h *FigureHandler) GetByID(c *gin.Context) {
	id, err := parseID(c)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "id inválido"})
		return
	}

	figurinha, err := h.svc.GetByID(id)
	if err != nil {
		c.JSON(mapErrorToStatus(err), gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, figurinha)
}

func (h *FigureHandler) Update(c *gin.Context) {
	id, err := parseID(c)

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "id inválido"})
		return
	}

	var req domain.UpdateFigureRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	figurinha, err := h.svc.Update(id, req)
	if err != nil {
		c.JSON(mapErrorToStatus(err), gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, figurinha)
}

func (h *FigureHandler) Delete(c *gin.Context) {
	id, err := parseID(c)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "id inválido"})
		return
	}

	if err := h.svc.Delete(id); err != nil {
		c.JSON(mapErrorToStatus(err), gin.H{"error": err.Error()})
		return
	}

	c.Status(http.StatusNoContent)
}

func parseID(c *gin.Context) (uint, error) {
	raw, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		return 0, err
	}
	return uint(raw), nil
}

func mapErrorToStatus(err error) int {
	switch {
	case errors.Is(err, service.ErrFigureNotFound):
		return http.StatusNotFound
	case errors.Is(err, service.ErrInvalidNumber),
		errors.Is(err, service.ErrInvalidType),
		errors.Is(err, service.ErrInvalidPosition):
		return http.StatusBadRequest
	default:
		return http.StatusInternalServerError
	}
}