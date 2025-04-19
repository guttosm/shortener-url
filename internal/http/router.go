package http

import (
	"github.com/gin-gonic/gin"
	_ "github.com/guttosm/url-shortener/docs"
	"github.com/guttosm/url-shortener/internal/middleware"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

// NewRouter sets up the HTTP routes for the application.
//
// Parameters:
// - handler (*Handler): The HTTP handler containing the logic for URL shortening and redirection.
//
// Returns:
// - *gin.Engine: The configured Gin router with all the necessary routes.
//
// Routes:
// - POST /api/shorten: Handles URL shortening requests.
// - GET /:shortID: Redirects to the original URL based on the shortened ID.
// - GET /swagger/*any: Serves the Swagger API documentation.
func NewRouter(handler *Handler) *gin.Engine {
	router := gin.Default()

	api := router.Group("/api")
	{
		api.POST("/login", handler.Login)
	}

	router.GET("/:shortID", handler.Redirect)
	router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	protected := router.Group("/api")
	protected.Use(middleware.AuthMiddleware())
	{
		protected.POST("/shorten", handler.ShortenURL)
	}

	return router
}
