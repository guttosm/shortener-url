package config

import (
	"os"
	"testing"

	"github.com/spf13/viper"
)

func TestLoadConfig_WithEnvVariables(t *testing.T) {
	// Set environment variables for testing
	os.Setenv("MONGO_URI", "mongodb://localhost:27017")
	os.Setenv("MONGO_DB", "test_db")
	os.Setenv("REDIS_URI", "redis://localhost:6379")
	os.Setenv("SERVER_PORT", "8080")
	os.Setenv("AUTH_USER_ID", "test-user-id")
	os.Setenv("AUTH_USERNAME", "test-user")
	os.Setenv("AUTH_PASSWORD", "test-password")

	// Reset viper to avoid conflicts with other tests
	viper.Reset()

	// Call LoadConfig
	LoadConfig()

	// Assertions
	if AppConfig.MongoURI != "mongodb://localhost:27017" {
		t.Errorf("Expected MongoURI to be 'mongodb://localhost:27017', got '%s'", AppConfig.MongoURI)
	}
	if AppConfig.MongoDB != "test_db" {
		t.Errorf("Expected MongoDB to be 'test_db', got '%s'", AppConfig.MongoDB)
	}
	if AppConfig.RedisURI != "redis://localhost:6379" {
		t.Errorf("Expected RedisURI to be 'redis://localhost:6379', got '%s'", AppConfig.RedisURI)
	}
	if AppConfig.ServerPort != "8080" {
		t.Errorf("Expected ServerPort to be '8080', got '%s'", AppConfig.ServerPort)
	}
	if AppConfig.Auth.UserID != "test-user-id" {
		t.Errorf("Expected Auth.UserID to be 'test-user-id', got '%s'", AppConfig.Auth.UserID)
	}
	if AppConfig.Auth.Username != "test-user" {
		t.Errorf("Expected Auth.Username to be 'test-user', got '%s'", AppConfig.Auth.Username)
	}
	if AppConfig.Auth.Password != "test-password" {
		t.Errorf("Expected Auth.Password to be 'test-password', got '%s'", AppConfig.Auth.Password)
	}

	// Clean up environment variables
	os.Unsetenv("MONGO_URI")
	os.Unsetenv("MONGO_DB")
	os.Unsetenv("REDIS_URI")
	os.Unsetenv("SERVER_PORT")
	os.Unsetenv("AUTH_USER_ID")
	os.Unsetenv("AUTH_USERNAME")
	os.Unsetenv("AUTH_PASSWORD")
}

func TestLoadConfig_WithoutEnvVariables(t *testing.T) {
	// Reset viper and environment variables
	viper.Reset()
	os.Clearenv()

	// Call LoadConfig
	LoadConfig()

	// Assertions for missing variables
	if AppConfig.MongoURI != "" {
		t.Errorf("Expected MongoURI to be empty, got '%s'", AppConfig.MongoURI)
	}
	if AppConfig.ServerPort != "" {
		t.Errorf("Expected ServerPort to be empty, got '%s'", AppConfig.ServerPort)
	}
}
