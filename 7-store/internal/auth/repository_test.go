package auth_test

import (
	"demo-store/configs"
	"demo-store/internal/auth"
	"demo-store/pkg/db"
	"log"
	"os"
	"testing"

	"github.com/joho/godotenv"
)

var repo *auth.AuthRepository

func TestMain(m *testing.M) {
	err := godotenv.Load("../../.env")
	if err != nil {
		log.Fatal("error load env file")
	}
	configs := &configs.Configs{
		Db: configs.DbConfig{
			Dsn: os.Getenv("DSN"),
		},
		Auth: configs.AuthConfig{
			Secret: "",
		},
	}
	repo = auth.NewAuthRepository(db.NewDb(configs))
	m.Run()
}

func TestGetByPhoneNotExist(t *testing.T) {
	phone := "not_exist_phone"
	phoneDB, _ := repo.GetByPhone(phone)

	if phoneDB != nil {
		t.Errorf("expected nil, got: %v", phoneDB)
	}
}
