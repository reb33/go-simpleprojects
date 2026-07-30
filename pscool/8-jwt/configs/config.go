package configs

import (
	"os"

	"github.com/joho/godotenv"
)

type Config struct {
	Secret string
}

func LoadConfig() *Config {
	err := godotenv.Load(".env")
	if err != nil {
		panic(err)
	}
	return &Config{
		Secret: os.Getenv("SECRET"),
	}
}
