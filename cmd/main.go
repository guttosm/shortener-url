package main

import (
    "context"
    "log"
    "os"
    "os/signal"
    "syscall"
    "time"
    "net/http"

    "github.com/guttosm/url-shortener/config"
    "github.com/guttosm/url-shortener/internal/app"
)

// main is the entry point of the application.
//
// Behavior:
// - Loads the application configuration.
// - Initializes the application by setting up dependencies (e.g., database, Redis, router).
// - Starts the HTTP server to handle incoming requests.
// - Listens for OS signals (e.g., SIGINT, SIGTERM) to gracefully shut down the server.
//
// Logs:
// - Logs the server startup and shutdown process.
// - Logs errors if the server fails to start or shut down gracefully.
func main() {
    config.LoadConfig()

    router, cleanup, err := app.InitializeApp()
    if err != nil {
        log.Fatal("Error on start up application:", err)
    }
    defer cleanup()

    // Create a channel to listen for OS signals
    quit := make(chan os.Signal, 1)
    signal.Notify(quit, os.Interrupt, syscall.SIGTERM)

    // Start the server in a goroutine
    server := &http.Server{
        Addr:    ":" + config.AppConfig.ServerPort,
        Handler: router,
    }

    go func() {
        log.Println("Server running on port " + config.AppConfig.ServerPort)
        if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
            log.Fatalf("Error starting server: %v", err)
        }
    }()

    <-quit
    log.Println("Shutting down server...")

    ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
    defer cancel()

    if err := server.Shutdown(ctx); err != nil {
        log.Fatalf("Server forced to shutdown: %v", err)
    }

    log.Println("Server exited gracefully")
}
