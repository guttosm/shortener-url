package http

import (
    "context"
    "net/http"

    "github.com/gin-gonic/gin"
    "github.com/guttosm/url-shortener/internal/auth"
    "github.com/guttosm/url-shortener/internal/dto"
    "github.com/guttosm/url-shortener/internal/service"
    "github.com/guttosm/url-shortener/config"
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

    // Retrieve credentials from the configuration
    configUsername := config.AppConfig.Auth.Username
    configPassword := config.AppConfig.Auth.Password
    configUserID := config.AppConfig.Auth.UserID

    // Authenticate the user
    if req.Username != configUsername || req.Password != configPassword {
        c.Error(dto.NewErrorResponse("Invalid credentials", nil))
        return
    }

    // Generate a JWT token using the user ID from the configuration
    token, err := auth.GenerateToken(configUserID)
    if err != nil {
        c.Error(dto.NewErrorResponse("Failed to generate token", err))
        return
    }

    // Return the token in the response
    c.JSON(http.StatusOK, gin.H{"token": token})
}
