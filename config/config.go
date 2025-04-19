package config

import (
    "log"
    "strings"

    "github.com/spf13/viper"
)

// Config holds the application configuration.
//
// Fields:
// - MongoURI (string): The MongoDB connection URI.
// - MongoDB (string): The name of the MongoDB database.
// - RedisURI (string): The Redis connection URI.
// - ServerPort (string): The port on which the server will run.
// - Auth (AuthConfig): The authentication configuration.
type Config struct {
    MongoURI   string
    MongoDB    string
    RedisURI   string
    ServerPort string
    Auth       AuthConfig
}

// AuthConfig holds the authentication configuration.
//
// Fields:
// - UserID (string): The user ID for authentication.
// - Username (string): The username for authentication.
// - Password (string): The password for authentication.
type AuthConfig struct {
    UserID   string
    Username string
    Password string
}

// AppConfig is the global instance of the application configuration.
var AppConfig *Config

// LoadConfig loads the application configuration from the .env file or environment variables.
//
// Behavior:
// - Reads configuration from a .env file if it exists.
// - Falls back to environment variables if the .env file is not found.
// - Populates the AppConfig global variable with the loaded configuration.
//
// Logs:
// - Logs whether the .env file was loaded or if environment variables are being used.
// - Logs a warning if required variables are missing.
func LoadConfig() {
    viper.SetConfigFile(".env")
    viper.SetConfigType("env")
    viper.AutomaticEnv()
    viper.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))

    // Attempt to read the .env file
    if err := viper.ReadInConfig(); err == nil {
        log.Println("File .env loaded")
    } else {
        log.Println(".env not found. Using environment variables")
    }

    // Populate the AppConfig struct
    AppConfig = &Config{
        MongoURI:   viper.GetString("MONGO_URI"),
        MongoDB:    viper.GetString("MONGO_DB"),
        RedisURI:   viper.GetString("REDIS_URI"),
        ServerPort: viper.GetString("SERVER_PORT"),
        Auth: AuthConfig{
            UserID:   viper.GetString("AUTH_USER_ID"),
            Username: viper.GetString("AUTH_USERNAME"),
            Password: viper.GetString("AUTH_PASSWORD"),
        },
    }

    // Check for missing required variables
    if AppConfig.MongoURI == "" || AppConfig.ServerPort == "" {
        log.Println("Some variables were not declared! Check file docker-compose.yml or environment variables.")
    }
}