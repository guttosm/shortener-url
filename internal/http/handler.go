package http

import (
    "context"
    "net/http"
    "fmt"
    "github.com/guttosm/url-shortener/internal/auth"

    "github.com/gin-gonic/gin"
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
// @Router /api/shorten [post]
//
// Behavior:
// - Validates the request body to ensure it contains a valid URL.
// - Calls the URL service to generate a shortened URL.
// - Returns the shortened URL in the response or an error if the operation fails.
func (h *Handler) ShortenURL(c *gin.Context) {
    var req dto.ShortenRequest

    if err := c.ShouldBindJSON(&req); err != nil {
        c.Error(err)
        return
    }

    urlEntity, err := h.urlService.Shorten(context.Background(), req.URL)
    if err != nil {
        c.Error(err)
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
//
// Behavior:
// - Extracts the short ID from the request parameters.
// - Calls the URL service to find the original URL associated with the short ID.
// - Redirects the user to the original URL or returns a 404 error if the short ID is not found.
func (h *Handler) Redirect(c *gin.Context) {
    shortID := c.Param("shortID")

    urlEntity, err := h.urlService.FindByShortID(context.Background(), shortID)
    if err != nil || urlEntity == nil {
        c.Error(fmt.Errorf("URL not found"))
        return
    }

    c.Redirect(http.StatusFound, urlEntity.Original)
}

// ProtectedResource handles requests to a protected resource.
//
// Parameters:
// - c (*gin.Context): The Gin context containing the HTTP request and response.
//
// Behavior:
// - Extracts the user ID from the context (set by the authentication middleware).
// - Returns the user ID and a success message if the user is authenticated.
// - Returns a 401 error if the user ID is not found in the context.
func (h *Handler) ProtectedResource(c *gin.Context) {
    userID, exists := c.Get("user_id")
    if !exists {
        c.JSON(http.StatusUnauthorized, gin.H{"error": "User ID not found"})
        return
    }

    c.JSON(http.StatusOK, gin.H{
        "message": "This is a protected resource",
        "user_id": userID,
    })
}

// Login handles user login and generates a JWT token.
//
// Parameters:
// - c (*gin.Context): The Gin context containing the HTTP request and response.
//
// Behavior:
// - Validates the request body to ensure it contains a username and password.
// - Authenticates the user (dummy logic for now).
// - Generates a JWT token if the credentials are valid.
// - Returns the token in the response or an error if authentication fails.
func (h *Handler) Login(c *gin.Context) {
    var req struct {
        Username string `json:"username" binding:"required"`
        Password string `json:"password" binding:"required"`
    }

    if err := c.ShouldBindJSON(&req); err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request"})
        return
    }


    if req.Username != "admin" || req.Password != "password" {
        c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid credentials"})
        return
    }

    token, err := auth.GenerateToken("user-id-123")
    if err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to generate token"})
        return
    }

    c.JSON(http.StatusOK, gin.H{"token": token})
}