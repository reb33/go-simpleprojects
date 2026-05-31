package configs

import (
	"log"
	"os"

	"github.com/joho/godotenv"
)

type SendlerConfig struct {
	Email    string
	Password string
}

type Config struct {
	SendlerConfig
}

func LoadConfig() *Config {
	err := godotenv.Load()
	if err != nil {
		log.Println("Error loading .env file")
	}
	return &Config{
		SendlerConfig: SendlerConfig{
			Email:    os.Getenv("EMAIL"),
			Password: os.Getenv("PASSWORD"),
		},
	}
}
