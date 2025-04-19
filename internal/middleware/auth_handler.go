// filepath: /Users/gmoraes/Documents/personal/it/url-shortener/internal/middleware/auth_middleware.go
package middleware

import (
    "net/http"
    "strings"

    "github.com/gin-gonic/gin"
    "github.com/guttosm/url-shortener/internal/auth"
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
            c.JSON(http.StatusUnauthorized, gin.H{"error": "Authorization header is required"})
            c.Abort()
            return
        }

        tokenString := strings.TrimPrefix(authHeader, "Bearer ")
        claims, err := auth.ValidateToken(tokenString)
        if err != nil {
            c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid or expired token"})
            c.Abort()
            return
        }

        // Store user ID in the context
        c.Set("user_id", claims["user_id"])
        c.Next()
    }
}