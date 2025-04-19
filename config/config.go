package config

import (
	"log"
	"strings"

	"github.com/spf13/viper"
)

type Config struct {
	MongoURI   string
	MongoDB    string
	RedisURI   string
	ServerPort string
}

var AppConfig *Config

func LoadConfig() {
	viper.SetConfigFile(".env")
	viper.SetConfigType("env")
	viper.AutomaticEnv()
	viper.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))

	if err := viper.ReadInConfig(); err == nil {
		log.Println("File .env loaded")
	} else {
		log.Println(".env not found. Using environment variables")
	}

	AppConfig = &Config{
		MongoURI:   viper.GetString("MONGO_URI"),
		MongoDB:    viper.GetString("MONGO_DB"),
		RedisURI:   viper.GetString("REDIS_URI"),
		ServerPort: viper.GetString("SERVER_PORT"),
	}

	if AppConfig.MongoURI == "" || AppConfig.ServerPort == "" {
		log.Println("Some variables was not declared! Check file docker-compose.yml")
	}
}
