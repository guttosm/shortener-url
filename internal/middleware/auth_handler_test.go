package middleware_test

import (
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/guttosm/url-shortener/internal/middleware"
	"github.com/stretchr/testify/assert"
)

// MockTokenValidator implements the auth.TokenValidator interface
type MockTokenValidator struct {
	ValidateTokenFunc func(token string) (map[string]interface{}, error)
}

func (m *MockTokenValidator) ValidateToken(token string) (map[string]interface{}, error) {
	return m.ValidateTokenFunc(token)
}

func TestAuthMiddleware_NoAuthHeader(t *testing.T) {
	gin.SetMode(gin.TestMode)
	router := gin.New()

	mockValidator := &MockTokenValidator{
		ValidateTokenFunc: func(token string) (map[string]interface{}, error) {
			return nil, nil
		},
	}

	router.Use(middleware.AuthMiddleware(mockValidator))
	router.GET("/protected", func(c *gin.Context) {
		c.Status(http.StatusOK)
	})

	req := httptest.NewRequest("GET", "/protected", nil) // No Authorization header
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusUnauthorized, w.Code)
}

func TestAuthMiddleware_InvalidToken(t *testing.T) {
	gin.SetMode(gin.TestMode)
	router := gin.New()

	mockValidator := &MockTokenValidator{
		ValidateTokenFunc: func(token string) (map[string]interface{}, error) {
			return nil, errors.New("invalid token")
		},
	}

	router.Use(middleware.AuthMiddleware(mockValidator))
	router.GET("/protected", func(c *gin.Context) {
		c.Status(http.StatusOK)
	})

	req := httptest.NewRequest("GET", "/protected", nil)
	req.Header.Set("Authorization", "Bearer invalid_token")
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusUnauthorized, w.Code)
}

func TestAuthMiddleware_ValidToken(t *testing.T) {
	gin.SetMode(gin.TestMode)
	router := gin.New()

	mockValidator := &MockTokenValidator{
		ValidateTokenFunc: func(token string) (map[string]interface{}, error) {
			return map[string]interface{}{"user_id": "123"}, nil
		},
	}

	router.Use(middleware.AuthMiddleware(mockValidator))
	router.GET("/protected", func(c *gin.Context) {
		userID, exists := c.Get("user_id")
		assert.True(t, exists)
		assert.Equal(t, "123", userID)
		c.JSON(http.StatusOK, gin.H{"user_id": userID})
	})

	req := httptest.NewRequest("GET", "/protected", nil)
	req.Header.Set("Authorization", "Bearer valid_token")
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	assert.Contains(t, w.Body.String(), "123")
}
