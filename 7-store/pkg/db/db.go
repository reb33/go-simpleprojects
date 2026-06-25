package db

import (
	"demo-store/configs"
	"log"
	"os"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

type Db struct {
	*gorm.DB
}

func NewDb(conf *configs.Configs) *Db {
	newLogger := logger.New(
		log.New(os.Stdout, "\r\n", log.LstdFlags),
		logger.Config{
			LogLevel:                  logger.Info,
			IgnoreRecordNotFoundError: true, // Ignore ErrRecordNotFound error for logger
			Colorful:                  true, 
		},
	)

	db, err := gorm.Open(postgres.Open(conf.Db.Dsn), &gorm.Config{
		Logger: newLogger,
	})
	if err != nil {
		panic(err)
	}
	return &Db{db}
}
