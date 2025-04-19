package http

import (
	"github.com/gin-gonic/gin"
	_ "github.com/guttosm/url-shortener/docs"
	"github.com/guttosm/url-shortener/internal/auth"
	"github.com/guttosm/url-shortener/internal/middleware"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

// NewRouter sets up the HTTP routes for the application.
//
// Parameters:
// - handler (*Handler): The HTTP handler containing the logic for URL shortening and login.
// - validator (auth.TokenValidator): The service used to validate JWT tokens.
//
// Returns:
// - *gin.Engine: The configured Gin router with public and protected endpoints.
func NewRouter(handler *Handler, validator auth.TokenValidator) *gin.Engine {
	router := gin.Default()

	// Swagger documentation
	router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	// Public routes
	public := router.Group("/api")
	{
		public.POST("/login", handler.Login)
	}

	// Protected routes
	protected := router.Group("/api")
	protected.Use(middleware.AuthMiddleware(validator))
	{
		protected.POST("/shorten", handler.ShortenURL)
	}

	return router
}
