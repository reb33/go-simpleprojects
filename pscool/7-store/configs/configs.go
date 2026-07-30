package configs

import (
	"log"
	"os"

	"github.com/joho/godotenv"
)

type DbConfig struct {
	Dsn string
}

type AuthConfig struct {
	Secret string
}

type Configs struct {
	Db   DbConfig
	Auth AuthConfig
}

func LoadConfigs() *Configs {
	err := godotenv.Load(".env")
	if err != nil {
		log.Fatal("error load env file")
	}
	return &Configs{
		Db: DbConfig{
			Dsn: os.Getenv("DSN"),
		},
		Auth: AuthConfig{
			Secret: os.Getenv("SECRET"),
		},
	}
}
