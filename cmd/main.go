package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/guttosm/url-shortener/config"
	"github.com/guttosm/url-shortener/internal/app"
)

// startServer starts the HTTP server and listens for incoming requests.
//
// Parameters:
// - router (http.Handler): The HTTP router to handle requests.
// - port (string): The port on which the server will run.
//
// Returns:
// - *http.Server: The initialized HTTP server.
func startServer(router http.Handler, port string) *http.Server {
	server := &http.Server{
		Addr:    ":" + port,
		Handler: router,
	}

	go func() {
		log.Println("Server running on port " + port)
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("Error starting server: %v", err)
		}
	}()

	return server
}

// gracefulShutdown shuts down the server gracefully when an OS signal is received.
//
// Parameters:
// - server (*http.Server): The HTTP server to shut down.
// - cleanup (func()): A cleanup function to release resources.
func gracefulShutdown(server *http.Server, cleanup func()) {
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt, syscall.SIGTERM)

	<-quit
	log.Println("Shutting down server...")

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if err := server.Shutdown(ctx); err != nil {
		log.Fatalf("Server forced to shutdown: %v", err)
	}

	cleanup()
	log.Println("Server exited gracefully")
}

func main() {
	config.LoadConfig()

	router, cleanup, err := app.InitializeApp()
	if err != nil {
		log.Fatal("Error on start up application:", err)
	}

	server := startServer(router, config.AppConfig.ServerPort)
	gracefulShutdown(server, cleanup)
}
