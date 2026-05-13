package config

import "os"

type Config struct {
	Key string
}

func NewConfig() *Config{
	key := os.Getenv("KEY")
	return &Config{
		Key : key,
	}
}