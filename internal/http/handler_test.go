package http_test

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/guttosm/url-shortener/config"
	"github.com/guttosm/url-shortener/internal/dto"
	"github.com/guttosm/url-shortener/internal/entity"
	apphttp "github.com/guttosm/url-shortener/internal/http"
	"github.com/stretchr/testify/assert"
)

// Mock para URLService
type mockURLServiceHandlerTest struct {
	shortenFunc       func(context.Context, string) (*entity.URL, error)
	findByShortIDFunc func(context.Context, string) (*entity.URL, error)
}

func (m *mockURLServiceHandlerTest) Shorten(ctx context.Context, originalURL string) (*entity.URL, error) {
	return m.shortenFunc(ctx, originalURL)
}

func (m *mockURLServiceHandlerTest) FindByShortID(ctx context.Context, shortID string) (*entity.URL, error) {
	return m.findByShortIDFunc(ctx, shortID)
}

func TestHandler_Login(t *testing.T) {
	authCfg := config.AuthConfig{
		Username: "admin",
		Password: "secret",
		UserID:   "user123",
	}

	handler := apphttp.NewHandler(&mockURLServiceHandlerTest{}, authCfg)

	router := gin.Default()
	router.POST("/login", handler.Login)

	t.Run("invalid JSON", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodPost, "/login", bytes.NewBufferString("invalid"))
		req.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()

		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusBadRequest, w.Code)
		assert.Contains(t, w.Body.String(), "Invalid request")
	})

	t.Run("invalid credentials", func(t *testing.T) {
		body, _ := json.Marshal(dto.LoginRequest{Username: "wrong", Password: "wrong"})
		req := httptest.NewRequest(http.MethodPost, "/login", bytes.NewBuffer(body))
		req.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()

		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusUnauthorized, w.Code)
		assert.Contains(t, w.Body.String(), "Invalid credentials")
	})

	t.Run("valid credentials", func(t *testing.T) {
		body, _ := json.Marshal(dto.LoginRequest{Username: "admin", Password: "secret"})
		req := httptest.NewRequest(http.MethodPost, "/login", bytes.NewBuffer(body))
		req.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()

		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)
		assert.Contains(t, w.Body.String(), "token")
	})
}

func TestHandler_ShortenURL(t *testing.T) {
	gin.SetMode(gin.TestMode)

	t.Run("invalid JSON", func(t *testing.T) {
		handler := apphttp.NewHandler(&mockURLService{}, config.AuthConfig{})
		router := gin.Default()
		router.POST("/shorten", handler.ShortenURL)

		req := httptest.NewRequest(http.MethodPost, "/shorten", bytes.NewBufferString("invalid"))
		req.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()

		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusBadRequest, w.Code)
		assert.Contains(t, w.Body.String(), "Invalid request")
	})

	t.Run("error from service", func(t *testing.T) {
		mockService := &mockURLServiceHandlerTest{
			shortenFunc: func(ctx context.Context, url string) (*entity.URL, error) {
				return nil, errors.New("fail")
			},
		}
		handler := apphttp.NewHandler(mockService, config.AuthConfig{})
		router := gin.Default()
		router.POST("/shorten", handler.ShortenURL)

		body, _ := json.Marshal(dto.ShortenRequest{URL: "https://example.com"})
		req := httptest.NewRequest(http.MethodPost, "/shorten", bytes.NewBuffer(body))
		req.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()

		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusInternalServerError, w.Code)
		assert.Contains(t, w.Body.String(), "Failed to shorten URL")
	})

	t.Run("success", func(t *testing.T) {
		mockService := &mockURLServiceHandlerTest{
			shortenFunc: func(ctx context.Context, url string) (*entity.URL, error) {
				return &entity.URL{
					ShortID:  "abc123",
					Original: url,
				}, nil
			},
		}
		handler := apphttp.NewHandler(mockService, config.AuthConfig{})
		router := gin.Default()
		router.POST("/shorten", handler.ShortenURL)

		body, _ := json.Marshal(dto.ShortenRequest{URL: "https://example.com"})
		req := httptest.NewRequest(http.MethodPost, "/shorten", bytes.NewBuffer(body))
		req.Header.Set("Content-Type", "application/json")
		req.Host = "localhost:8080"
		w := httptest.NewRecorder()

		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)
		assert.Contains(t, w.Body.String(), "abc123")
		assert.Contains(t, w.Body.String(), "localhost:8080/abc123")
	})
}
