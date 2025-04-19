package middleware

import (
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/guttosm/url-shortener/internal/auth"
	"github.com/guttosm/url-shortener/internal/dto"
)

// AuthMiddleware returns a Gin middleware that validates JWT tokens
// and injects the user ID into the context if the token is valid.
//
// Parameters:
// - validator (auth.TokenValidator): The service responsible for validating JWT tokens.
//
// Behavior:
// - Extracts the Bearer token from the Authorization header.
// - Validates the token and retrieves claims.
// - Stores the user_id in the Gin context for use in handlers.
// - Returns 401 Unauthorized if validation fails.
func AuthMiddleware(validator auth.TokenValidator) gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" || !strings.HasPrefix(authHeader, "Bearer ") {
			_ = c.Error(dto.NewErrorResponse("Authorization header is required", nil))
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Authorization header is required"})
			return
		}

		tokenString := strings.TrimPrefix(authHeader, "Bearer ")
		claims, err := validator.ValidateToken(tokenString)
		if err != nil {
			c.Error(dto.NewErrorResponse("Invalid or expired token", err))
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Invalid or expired token"})
			return
		}

		userID, ok := claims["user_id"]
		if !ok {
			c.Error(dto.NewErrorResponse("user_id claim is missing", nil))
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Invalid token claims"})
			return
		}

		c.Set("user_id", userID)
		c.Next()
	}
}
