package middleware

import (
    "net/http"
    "strings"

    "github.com/gin-gonic/gin"
    "github.com/guttosm/url-shortener/internal/auth"
    "github.com/guttosm/url-shortener/internal/dto"
)

// AuthMiddleware validates the JWT token and extracts user information.
//
// Behavior:
// - Checks for the Authorization header.
// - Validates the JWT token.
// - Stores the user ID in the Gin context for further use.
func AuthMiddleware() gin.HandlerFunc {
    return func(c *gin.Context) {
        authHeader := c.GetHeader("Authorization")
        if authHeader == "" {
            _ = c.Error(dto.NewErrorResponse("Authorization header is required", nil)) // Register the error
            c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Authorization header is required"})
            return
        }

        tokenString := strings.TrimPrefix(authHeader, "Bearer ")
        claims, err := auth.ValidateToken(tokenString)
        if err != nil {
            c.Error(dto.NewErrorResponse("Invalid or expired token", err)) // Register the error
            c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Invalid or expired token"})
            return
        }

        // Store user ID in the context
        c.Set("user_id", claims["user_id"])
        c.Next()
    }
}