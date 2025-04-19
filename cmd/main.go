// @title Shortener URL - API
// @version 1.0
// @description This is a URL shortener API.
// @description It allows you to shorten long URLs and retrieve the original URL using the shortened version.
// @host localhost:8080
// @BasePath /
package main

import (
	"log"

	"github.com/guttosm/url-shortener/config"
	"github.com/guttosm/url-shortener/internal/app"
)

func main() {
	config.LoadConfig()

	router, cleanup, err := app.InitializeApp()
	if err != nil {
		log.Fatal("Error on start up application:", err)
	}
	defer cleanup()

	log.Println("Server running in port " + config.AppConfig.ServerPort)

	if err := router.Run(":" + config.AppConfig.ServerPort); err != nil {
		log.Fatal("Erro ao iniciar o servidor HTTP:", err)
	}
}
