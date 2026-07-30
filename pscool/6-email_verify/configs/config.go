package configs

import (
	"log"
	"os"

	"github.com/joho/godotenv"
)

type SendlerConfig struct {
	Email    string
	Password string
	Address  string
	PORT     string
}

type Config struct {
	SendlerConfig
	VerifyURL string
}

func LoadConfig() *Config {
	err := godotenv.Load(".env")
	if err != nil {
		log.Println("Error loading .env file")
	}
	return &Config{
		SendlerConfig: SendlerConfig{
			Email:    os.Getenv("EMAIL"),
			Password: os.Getenv("PASSWORD"),
			Address:  os.Getenv("ADDRESS"),
			PORT:     os.Getenv("PORT"),
		},
		VerifyURL: os.Getenv("VERIFY_URL"),
	}
}
