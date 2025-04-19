package http

import (
	"context"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/guttosm/url-shortener/config"
	"github.com/guttosm/url-shortener/internal/auth"
	"github.com/guttosm/url-shortener/internal/dto"
	"github.com/guttosm/url-shortener/internal/middleware"
	"github.com/guttosm/url-shortener/internal/service"
)

type Handler struct {
	urlService service.URLService
	authConfig config.AuthConfig
}

func NewHandler(s service.URLService, authCfg config.AuthConfig) *Handler {
	return &Handler{
		urlService: s,
		authConfig: authCfg,
	}
}

// ShortenURL handles the shortening of a URL.
func (h *Handler) ShortenURL(c *gin.Context) {
	var req dto.ShortenRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		middleware.AbortWithError(c, http.StatusBadRequest, "Invalid request", err)
		return
	}

	urlEntity, err := h.urlService.Shorten(context.Background(), req.URL)
	if err != nil {
		middleware.AbortWithError(c, http.StatusInternalServerError, "Failed to shorten URL", err)
		return
	}

	resp := dto.ShortenResponse{
		ShortID:  urlEntity.ShortID,
		ShortURL: c.Request.Host + "/" + urlEntity.ShortID,
	}
	c.JSON(http.StatusOK, resp)
}

// Login handles user login and generates a JWT token.
func (h *Handler) Login(c *gin.Context) {
	var req dto.LoginRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		middleware.AbortWithError(c, http.StatusBadRequest, "Invalid request", err)
		return
	}

	if req.Username != h.authConfig.Username || req.Password != h.authConfig.Password {
		middleware.AbortWithError(c, http.StatusUnauthorized, "Invalid credentials", nil)
		return
	}

	token, err := auth.GenerateToken(h.authConfig.UserID)
	if err != nil {
		middleware.AbortWithError(c, http.StatusInternalServerError, "Failed to generate token", err)
		return
	}

	c.JSON(http.StatusOK, gin.H{"token": token})
}
