package middleware

import (
	"fmt"
	"net/http"
	"runtime/debug"

	"github.com/gin-gonic/gin"
	"github.com/guttosm/url-shortener/internal/dto"
)

// RecoveryMiddleware is a middleware that recovers from panics,
// logs the stack trace, and returns a standardized 500 error response.
func RecoveryMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		defer func() {
			if r := recover(); r != nil {
				// Log para debugging
				fmt.Printf("[PANIC RECOVERED] %v\n%s\n", r, debug.Stack())

				errResponse := dto.NewErrorResponse("Internal server error", fmt.Errorf("%v", r))
				c.AbortWithStatusJSON(http.StatusInternalServerError, errResponse)
			}
		}()

		c.Next()
	}
}
