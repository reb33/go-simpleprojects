package repository_test

import (
	"email-verify/internal/repository"
	"encoding/json"
	"os"
	"testing"
)

var FILE_NAME = "test.json"

func deleteFile() {
	if err := os.Remove(FILE_NAME); err != nil {
	}
}

func readAndDecodeFile() (*repository.Store, error) {
	data, err := os.ReadFile(FILE_NAME)
	if err != nil {
		return nil, err
	}
	var store repository.Store
	err = json.Unmarshal(data, &store)
	if err != nil {
		return nil, err
	}
	return &store, nil
}

func createFileWithData() {
	store := repository.Store{
		Emails: []repository.Email{
			{
				Email: "test@test.com",
				Hash:  "testhash",
			},
			{
				Email: "test2@test.com",
				Hash:  "testhash2",
			},
		},
	}
	data, err := json.Marshal(store)
	if err != nil {
		panic(err)
	}
	if err := os.WriteFile(FILE_NAME, data, 0644); err != nil {
		panic(err)
	}
}

func TestLoadNotExistFile(t *testing.T) {
	deleteFile()
	defer deleteFile()
	repository.NewEmailRepository(FILE_NAME)

	_, err := os.Stat(FILE_NAME)
	if err != nil {
		t.Fatal("Файл должен был создаться")
	}
}

func TestAdd(t *testing.T) {
	defer deleteFile()
	repo := repository.NewEmailRepository(FILE_NAME)

	repo.Add("test@test.com", "testhash")

	if len(repo.Store.Emails) != 1 {
		t.Fatal("Ошибка добавления 1 email")
	}

	repo.Add("test@test.com", "testhash")

	if len(repo.Store.Emails) != 2 {
		t.Fatal("Ошибка добавления 2 email")
	}

	store, err := readAndDecodeFile()
	if err != nil {
		t.Fatal(err)
	}
	if len(store.Emails) != 2 {
		t.Fatal("Ошибка сохранения данных в файл")
	}
}

func TestLoadAndAddFileExist(t *testing.T) {
	deleteFile()
	defer deleteFile()
	createFileWithData()
	repo := repository.NewEmailRepository(FILE_NAME)

	if len(repo.Store.Emails) != 2 {
		t.Fatal("Ошибка загрузки данных из файла")
	}

	repo.Add("test3@test.com", "testhash3")
	
	if len(repo.Store.Emails) != 3 {
		t.Fatal("Ошибка добавления данных")
	}
}

func TestGet(t *testing.T) {
	deleteFile()
	defer deleteFile()
	createFileWithData()
	repo := repository.NewEmailRepository(FILE_NAME)

	if len(repo.Store.Emails) != 2 {
		t.Fatal("Ошибка загрузки данных из файла")
	}

	email, err := repo.Get("testhash")
	if err != nil {
		t.Fatal(err)
	}
	if email != "test@test.com" {
		t.Fatal("Ошибка получения email")
	}
}

func TestDelete(t *testing.T) {
	deleteFile()
	defer deleteFile()
	createFileWithData()
	repo := repository.NewEmailRepository(FILE_NAME)

	if len(repo.Store.Emails) != 2 {
		t.Fatal("Ошибка загрузки данных из файла")
	}

	repo.Delete("test@test.com")
	if len(repo.Store.Emails) != 1 {
		t.Fatal("Ошибка удаления email")
	}

	store, err := readAndDecodeFile()
	if err != nil {
		t.Fatal(err)
	}
	if len(store.Emails) != 1 {
		t.Fatal("Ошибка сохранения данных в файл")
	}
}

func TestDeleteAndAdd(t *testing.T) {
	deleteFile()
	defer deleteFile()
	createFileWithData()
	repo := repository.NewEmailRepository(FILE_NAME)

	if len(repo.Store.Emails) != 2 {
		t.Fatal("Ошибка загрузки данных из файла")
	}

	repo.Delete("test@test.com")
	repo.Delete("test2@test.com")
	if len(repo.Store.Emails) != 0 {
		t.Fatal("Ошибка удаления email")
	}

	repo.Add("test@test.com", "testhash")
	if len(repo.Store.Emails) != 1 {
		t.Fatal("Ошибка добавления email")
	}
}
