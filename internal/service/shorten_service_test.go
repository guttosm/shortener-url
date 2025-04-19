
package service_test

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/guttosm/url-shortener/internal/entity"
	"github.com/guttosm/url-shortener/internal/service"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

type MockURLRepository struct {
	mock.Mock
}

func (m *MockURLRepository) Save(ctx context.Context, url *entity.URL) error {
	args := m.Called(ctx, url)
	return args.Error(0)
}

func (m *MockURLRepository) FindByOriginalURL(ctx context.Context, originalURL string) (*entity.URL, error) {
	args := m.Called(ctx, originalURL)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*entity.URL), args.Error(1)
}

func (m *MockURLRepository) FindByShortID(ctx context.Context, shortID string) (*entity.URL, error) {
	args := m.Called(ctx, shortID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*entity.URL), args.Error(1)
}

type MockURLCacheRepository struct {
	mock.Mock
}

func (m *MockURLCacheRepository) GetByOriginalURL(ctx context.Context, originalURL string) (*entity.URL, error) {
	args := m.Called(ctx, originalURL)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*entity.URL), args.Error(1)
}

func (m *MockURLCacheRepository) SetByOriginalURL(ctx context.Context, url *entity.URL) error {
	args := m.Called(ctx, url)
	return args.Error(0)
}

func (m *MockURLCacheRepository) GetByShortID(ctx context.Context, shortID string) (*entity.URL, error) {
	args := m.Called(ctx, shortID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*entity.URL), args.Error(1)
}

func (m *MockURLCacheRepository) SetByShortID(ctx context.Context, url *entity.URL) error {
	args := m.Called(ctx, url)
	return args.Error(0)
}

func TestShorten_URLExistsInCache(t *testing.T) {
	ctx := context.Background()
	original := "https://example.com"
	expected := &entity.URL{
		ShortID:   "abc123",
		Original:  original,
		CreatedAt: time.Now(),
	}

	cache := new(MockURLCacheRepository)
	repo := new(MockURLRepository)

	cache.On("GetByOriginalURL", ctx, original).Return(expected, nil)

	svc := service.NewURLService(repo, cache)
	result, err := svc.Shorten(ctx, original)

	assert.NoError(t, err)
	assert.Equal(t, expected, result)
	cache.AssertCalled(t, "GetByOriginalURL", ctx, original)
}

func TestShorten_URLExistsInDB(t *testing.T) {
	ctx := context.Background()
	original := "https://example.com"
	expected := &entity.URL{
		ShortID:   "abc123",
		Original:  original,
		CreatedAt: time.Now(),
	}

	cache := new(MockURLCacheRepository)
	repo := new(MockURLRepository)

	cache.On("GetByOriginalURL", ctx, original).Return(nil, errors.New("cache miss"))
	repo.On("FindByOriginalURL", ctx, original).Return(expected, nil)
	cache.On("SetByOriginalURL", ctx, expected).Return(nil)
	cache.On("SetByShortID", ctx, expected).Return(nil)

	svc := service.NewURLService(repo, cache)
	result, err := svc.Shorten(ctx, original)

	assert.NoError(t, err)
	assert.Equal(t, expected, result)
}

func TestShorten_URLDoesNotExist(t *testing.T) {
	ctx := context.Background()
	original := "https://example.com"

	cache := new(MockURLCacheRepository)
	repo := new(MockURLRepository)

	cache.On("GetByOriginalURL", ctx, original).Return(nil, errors.New("cache miss"))
	repo.On("FindByOriginalURL", ctx, original).Return((*entity.URL)(nil), nil)
	repo.On("Save", mock.Anything, mock.AnythingOfType("*entity.URL")).Return(nil)
	cache.On("SetByOriginalURL", mock.Anything, mock.AnythingOfType("*entity.URL")).Return(nil)
	cache.On("SetByShortID", mock.Anything, mock.AnythingOfType("*entity.URL")).Return(nil)

	svc := service.NewURLService(repo, cache)
	result, err := svc.Shorten(ctx, original)

	assert.NoError(t, err)
	assert.Equal(t, original, result.Original)
	assert.Len(t, result.ShortID, 6)
}
