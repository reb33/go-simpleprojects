package api_test

import (
	"3-struct/api"
	"3-struct/config"
	"encoding/json"
	"log"
	"testing"

	"github.com/joho/godotenv"
)

var binBody = []byte(`{"sample":"Hello World"}`)
var appAPI *api.Api

type Result struct {
    Record json.RawMessage `json:"record"`
}

func TestMain(m *testing.M) {
	// Инициализация перед всеми тестами
	err := godotenv.Load("../.env")
	if err != nil {
		log.Fatal("не удалось загрузить .env")
	}
	conf := config.NewConfig("")
	appAPI = api.NewApi(*conf)

	// Запуск всех тестов
	m.Run()
}

func TestCreateBin(t *testing.T) {
	meta, err := appAPI.CreateBin(binBody)
	if err != nil {
		t.Fatalf("Ошибка создания бина: %v", err.Error())
	}
	if meta == nil {
		t.Fatalf("Ошибка создания бина: метаданные не получены")
	}
	if meta.Id == "" {
		t.Fatalf("Ошибка создания бина: ID не получен")
	}
	appAPI.DeleteBin(meta.Id)
}

func TestGetBin(t *testing.T) {
	meta, err := appAPI.CreateBin(binBody)
	if err != nil {
		t.Fatalf("Ошибка создания бина: %v", err.Error())
	}

	bin, err := appAPI.GetBin(meta.Id)
	if err != nil {
		t.Fatalf("Ошибка получения бина: %v", err.Error())
	}
	if bin == nil {
		t.Fatalf("Ошибка получения бина: бин не получен")
	}
	var result Result
	err = json.Unmarshal(bin, &result)
	if string(result.Record) != string(binBody) {
		t.Errorf("Ожидалось: %s, получено: %s", string(binBody), result.Record)
	}

	appAPI.DeleteBin(meta.Id)
}

func TestUnmarshal(t *testing.T) {
	bin := []byte(`{"record":{"sample":"Hello World"},"metadata":{"id":"6a08b16ac0954111d832cbe4","private":true,"createdAt":"2026-05-16T18:03:22.148Z"}}`)
	var result Result
	err := json.Unmarshal(bin, &result)
	if err != nil {
		t.Fatalf("Ошибка получения бина: %v", err.Error())
	}
	if string(result.Record) != string(binBody) {
		t.Errorf("Ожидалось: %s, получено: %s", string(binBody), result.Record)
	}
}

func TestUpdateBin(t *testing.T) {
	meta, err := appAPI.CreateBin(binBody)
	if err != nil {
		t.Fatalf("Ошибка создания бина: %v", err.Error())
	}

	err = appAPI.UpdateBin(meta.Id, binBody)
	if err != nil {
		t.Fatalf("Ошибка обновления бина: %v", err.Error())
	}

	appAPI.DeleteBin(meta.Id)
}

func TestDeleteBin(t *testing.T) {
	meta, err := appAPI.CreateBin(binBody)
	if err != nil {
		t.Fatalf("Ошибка создания бина: %v", err.Error())
	}

	err = appAPI.DeleteBin(meta.Id)
	if err != nil {
		t.Fatalf("Ошибка удаления бина: %v", err.Error())
	}
}

