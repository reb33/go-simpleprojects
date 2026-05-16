package config

import "os"

type Config struct {
	Key string
	ApiUrl string
}

func NewConfig(url string) *Config{
	key := os.Getenv("KEY")
	return &Config{
		Key : key,
		ApiUrl: url,
	}
}