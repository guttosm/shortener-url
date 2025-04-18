package http

import (
	"context"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/guttosm/url-shortener/internal/service"
)

type Handler struct {
	urlService service.URLService
}

func NewHandler(s service.URLService) *Handler {
	return &Handler{urlService: s}
}

type shortenRequest struct {
	URL string `json:"url" binding:"required,url"`
}

// ShortenURL encurta uma URL.
// @Summary Encurtar URL
// @Description Recebe uma URL longa e retorna uma versão encurtada.
// @Tags URLs
// @Accept json
// @Produce json
// @Param url body entity.URL true "URL para encurtar"
// @Success 200 {object} map[string]string "Response with shorter URL"
// @Router /api/shorten [post]
func (h *Handler) ShortenURL(c *gin.Context) {
	var req shortenRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid URL"})
		return
	}

	urlEntity, err := h.urlService.Shorten(context.Background(), req.URL)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error creating short URL"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"short_id":  urlEntity.ShortID,
		"short_url": c.Request.Host + "/" + urlEntity.ShortID,
	})
}

func (h *Handler) Redirect(c *gin.Context) {
	shortID := c.Param("shortID")

	urlEntity, err := h.urlService.FindByShortID(context.Background(), shortID)
	if err != nil || urlEntity == nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "URL not found"})
		return
	}

	c.Redirect(http.StatusFound, urlEntity.Original)
}
