package http

import (
    "context"
    "net/http"

    "github.com/gin-gonic/gin"
    "github.com/guttosm/url-shortener/internal/auth"
    "github.com/guttosm/url-shortener/internal/dto"
    "github.com/guttosm/url-shortener/internal/service"
)

// Handler handles HTTP requests for URL shortening, redirection, and user authentication.
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
// @Failure 400 {object} dto.ErrorResponse "Invalid request"
// @Failure 500 {object} dto.ErrorResponse "Failed to shorten URL"
// @Router /api/shorten [post]
//
// Behavior:
// - Validates the request body to ensure it contains a valid URL.
// - Calls the URL service to generate a shortened URL.
// - Returns the shortened URL in the response or an error if the operation fails.
func (h *Handler) ShortenURL(c *gin.Context) {
    var req dto.ShortenRequest

    // Validate the request body
    if err := c.ShouldBindJSON(&req); err != nil {
        c.Error(dto.NewErrorResponse("Invalid request", err))
        return
    }

    // Call the URL service to shorten the URL
    urlEntity, err := h.urlService.Shorten(context.Background(), req.URL)
    if err != nil {
        c.Error(dto.NewErrorResponse("Failed to shorten URL", err))
        return
    }

    // Return the shortened URL in the response
    resp := dto.ShortenResponse{
        ShortID:  urlEntity.ShortID,
        ShortURL: c.Request.Host + "/" + urlEntity.ShortID,
    }
    c.JSON(http.StatusOK, resp)
}

// Redirect handles the redirection of a shortened URL to its original URL.
//
// @Summary Redirect to Original URL
// @Description Redirects the user to the original URL based on the shortened ID.
// @Tags URLs
// @Produce plain
// @Param shortID path string true "Shortened URL ID"
// @Success 302 "Redirects to the original URL"
// @Failure 404 {object} dto.ErrorResponse "URL not found"
// @Router /{shortID} [get]
//
// Behavior:
// - Extracts the short ID from the request parameters.
// - Calls the URL service to find the original URL associated with the short ID.
// - Redirects the user to the original URL or returns a 404 error if the short ID is not found.
func (h *Handler) Redirect(c *gin.Context) {
    shortID := c.Param("shortID")

    // Call the URL service to find the original URL
    urlEntity, err := h.urlService.FindByShortID(context.Background(), shortID)
    if err != nil || urlEntity == nil {
        c.Error(dto.NewErrorResponse("URL not found", err))
        return
    }

    // Redirect to the original URL
    c.Redirect(http.StatusFound, urlEntity.Original)
}

// Login handles user login and generates a JWT token.
//
// @Summary User Login
// @Description Authenticates the user and generates a JWT token.
// @Tags Authentication
// @Accept json
// @Produce json
// @Param credentials body dto.LoginRequest true "User credentials"
// @Success 200 {object} map[string]string "JWT token"
// @Failure 400 {object} dto.ErrorResponse "Invalid request"
// @Failure 401 {object} dto.ErrorResponse "Invalid credentials"
// @Failure 500 {object} dto.ErrorResponse "Failed to generate token"
// @Router /api/login [post]
func (h *Handler) Login(c *gin.Context) {
    var req dto.LoginRequest

    // Validate the request body
    if err := c.ShouldBindJSON(&req); err != nil {
        c.Error(dto.NewErrorResponse("Invalid request", err))
        return
    }

    // Authenticate the user
    if req.Username != "admin" || req.Password != "password" {
        c.Error(dto.NewErrorResponse("Invalid credentials", nil))
        return
    }

    // Generate a JWT token
    token, err := auth.GenerateToken("user-id-123")
    if err != nil {
        c.Error(dto.NewErrorResponse("Failed to generate token", err))
        return
    }

    // Return the token in the response
    c.JSON(http.StatusOK, gin.H{"token": token})
}
