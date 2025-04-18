package http

import (
	"github.com/gin-gonic/gin"
	_ "github.com/guttosm/url-shortener/docs"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

func NewRouter(handler *Handler) *gin.Engine {
	router := gin.Default()

	api := router.Group("/api")
	{
		api.POST("/shorten", handler.ShortenURL)
	}

	router.GET("/:shortID", handler.Redirect)
	router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	return router
}
