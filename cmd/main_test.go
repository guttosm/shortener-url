package main

import (
	"context"
	"net/http"
	"syscall"
	"testing"
	"time"
)

// TestStartServer tests the startServer function.
func TestStartServer(t *testing.T) {
	router := http.NewServeMux()
	router.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	})

	server := startServer(router, "8081")
	defer server.Shutdown(context.Background())

	// Give the server some time to start
	time.Sleep(100 * time.Millisecond)

	// Send a test request to the server
	resp, err := http.Get("http://localhost:8081/health")
	if err != nil {
		t.Fatalf("Failed to send request: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		t.Errorf("Expected status code 200, got %d", resp.StatusCode)
	}
}

// TestGracefulShutdown tests the gracefulShutdown function.
func TestGracefulShutdown(t *testing.T) {
	router := http.NewServeMux()
	server := startServer(router, "8082")

	cleanupCalled := false
	cleanup := func() {
		cleanupCalled = true
	}

	go func() {
		time.Sleep(100 * time.Millisecond)
		// Simulate sending a shutdown signal
		syscall.Kill(syscall.Getpid(), syscall.SIGTERM)
	}()

	gracefulShutdown(server, cleanup)

	if !cleanupCalled {
		t.Error("Expected cleanup to be called, but it was not")
	}
}
