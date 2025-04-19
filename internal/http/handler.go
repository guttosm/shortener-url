package http

import (
	"context"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/guttosm/url-shortener/internal/dto"
	"github.com/guttosm/url-shortener/internal/service"
)

// Handler handles HTTP requests for URL shortening and redirection.
//
// Fields:
// - urlService (service.URLService): The service responsible for URL-related operations.
type Handler struct {
	urlService service.URLService
}

// NewHandler creates a new instance of Handler.
//
// Parameters:
// - s (service.URLService): The URL service to be used by the handler.
//
// Returns:
// - *Handler: A new Handler instance.
func NewHandler(s service.URLService) *Handler {
	return &Handler{urlService: s}
}

// ShortenURL handles the shortening of a URL.
//
// @Summary Shorten a URL
// @Description Receives a long URL and returns a shortened version.
// @Tags URLs
// @Accept json
// @Produce json
// @Param url body dto.ShortenRequest true "URL" example({"url": "https://www.someurl.com"})
// @Success 200 {object} dto.ShortenResponse "Response with shorter URL"
// @Router /api/shorten [post]
func (h *Handler) ShortenURL(c *gin.Context) {
	var req dto.ShortenRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid URL"})
		return
	}

	urlEntity, err := h.urlService.Shorten(context.Background(), req.URL)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error creating short URL"})
		return
	}

	resp := dto.ShortenResponse{
		ShortID:  urlEntity.ShortID,
		ShortURL: c.Request.Host + "/" + urlEntity.ShortID,
	}

	c.JSON(http.StatusOK, resp)
}

// Redirect handles the redirection of a shortened URL to its original URL.
//
// Parameters:
// - c (*gin.Context): The Gin context containing the HTTP request and response.
func (h *Handler) Redirect(c *gin.Context) {
	shortID := c.Param("shortID")

	urlEntity, err := h.urlService.FindByShortID(context.Background(), shortID)
	if err != nil || urlEntity == nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "URL not found"})
		return
	}

	c.Redirect(http.StatusFound, urlEntity.Original)
}
