package middleware

import (
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/guttosm/url-shortener/internal/auth"
)

// AuthMiddleware returns a Gin middleware that validates JWT tokens
// and injects the user ID into the context if the token is valid.
func AuthMiddleware(validator auth.TokenValidator) gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" || !strings.HasPrefix(authHeader, "Bearer ") {
			AbortWithError(c, http.StatusUnauthorized, "Authorization header is required", nil)
			return
		}

		tokenString := strings.TrimPrefix(authHeader, "Bearer ")
		claims, err := validator.ValidateToken(tokenString)
		if err != nil {
			AbortWithError(c, http.StatusUnauthorized, "Invalid or expired token", err)
			return
		}

		userID, ok := claims["user_id"]
		if !ok {
			AbortWithError(c, http.StatusUnauthorized, "user_id claim is missing", nil)
			return
		}

		c.Set("user_id", userID)
		c.Next()
	}
}
