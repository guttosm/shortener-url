package http_test

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/guttosm/url-shortener/config"
	"github.com/guttosm/url-shortener/internal/auth"
	"github.com/guttosm/url-shortener/internal/entity"
	apphttp "github.com/guttosm/url-shortener/internal/http"
	"github.com/stretchr/testify/assert"
)

// --- Mock para URLService (usado pelo handler) ---

type mockURLService struct{}

func (m *mockURLService) Shorten(ctx context.Context, originalURL string) (*entity.URL, error) {
	return &entity.URL{
		ShortID:  "abc123",
		Original: originalURL,
	}, nil
}

func (m *mockURLService) FindByShortID(ctx context.Context, shortID string) (*entity.URL, error) {
	return &entity.URL{
		ShortID:  shortID,
		Original: "https://original.url",
	}, nil
}

// --- Mock para TokenValidator (usado pelo middleware) ---

type mockTokenValidator struct{}

func (m *mockTokenValidator) ValidateToken(token string) (map[string]interface{}, error) {
	if token == "valid-token" {
		return map[string]interface{}{"user_id": "123"}, nil
	}
	return nil, auth.ErrInvalidToken
}

// --- Testes ---

func TestRouter_PublicLoginRoute(t *testing.T) {
	authCfg := config.AuthConfig{
		Username: "any",
		Password: "any",
		UserID:   "user123",
	}

	router := apphttp.NewRouter(
		apphttp.NewHandler(&mockURLService{}, authCfg),
		&mockTokenValidator{},
	)

	loginPayload := map[string]string{
		"username": "any",
		"password": "any",
	}
	body, _ := json.Marshal(loginPayload)

	req := httptest.NewRequest(http.MethodPost, "/api/login", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	assert.Contains(t, w.Body.String(), "token")
}

func TestRouter_ProtectedShorten_Unauthorized(t *testing.T) {
	router := apphttp.NewRouter(
		apphttp.NewHandler(&mockURLService{}, config.AuthConfig{}),
		&mockTokenValidator{},
	)

	req := httptest.NewRequest(http.MethodPost, "/api/shorten", nil)
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusUnauthorized, w.Code)
	assert.Contains(t, w.Body.String(), "Authorization header is required")
}

func TestRouter_ProtectedShorten_Authorized(t *testing.T) {
	router := apphttp.NewRouter(
		apphttp.NewHandler(&mockURLService{}, config.AuthConfig{}),
		&mockTokenValidator{},
	)

	shortenPayload := map[string]string{
		"url": "https://example.com",
	}
	body, _ := json.Marshal(shortenPayload)

	req := httptest.NewRequest(http.MethodPost, "/api/shorten", bytes.NewBuffer(body))
	req.Header.Set("Authorization", "Bearer valid-token")
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	assert.Contains(t, w.Body.String(), "short_url")
}

func TestRouter_SwaggerRoute(t *testing.T) {
	router := apphttp.NewRouter(
		apphttp.NewHandler(&mockURLService{}, config.AuthConfig{}),
		&mockTokenValidator{},
	)

	req := httptest.NewRequest(http.MethodGet, "/swagger/index.html", nil)
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	assert.Contains(t, []int{
		http.StatusOK,
		http.StatusMovedPermanently,
		http.StatusFound,
		http.StatusNotFound, // caso os arquivos não estejam sendo servidos
	}, w.Code)
}
