package repository

import (
	"encoding/json"
	"errors"
	"os"
)

type EmailRepository struct {
	FileName string
	Store     *Store
}

type Email struct {
	Email    string `json:"email"`
	Hash string `json:"hash"`
}

type Store struct {
	Emails []Email `json:"emails"`
}

func NewEmailRepository(fileName string) *EmailRepository {
	repo := EmailRepository{
		FileName: fileName,
	}
	repo.Load()

	return &repo
}

func (repo *EmailRepository) Load() error {
	if _, err := os.Stat(repo.FileName); errors.Is(err, os.ErrNotExist) {
		return repo.createEmptyStore()
	}
	data, err := os.ReadFile(repo.FileName)
	if err != nil {
		return err
	}
	var store Store
	err = json.Unmarshal(data, &store)
	if err != nil {
		return repo.createEmptyStore()
	}
	repo.Store = &store
	return nil
}

func (repo *EmailRepository) createEmptyStore() error {
	repo.Store = &Store{
		Emails: []Email{},
	}
	return repo.Save()
}

func (repo *EmailRepository) Save() error {
	data, err := json.Marshal(repo.Store)
	if err != nil {
		return err
	}
	os.WriteFile(repo.FileName, data, 0644)
	return nil
}

func (repo *EmailRepository) Add(email string, hash string) error {
	repo.Store.Emails = append(repo.Store.Emails, Email{
		Email: email,
		Hash: hash,
	})
	return repo.Save()
}

func (repo *EmailRepository) Delete(email string) error {
	var emails []Email
	for _, row := range repo.Store.Emails {
		if row.Email == email {
			continue
		}
		emails = append(emails, row)
	}
	repo.Store.Emails = emails
	return repo.Save()
}

func (repo *EmailRepository) Get(hash string) (string, error) {
	for _, email := range repo.Store.Emails {
		if email.Hash == hash {
			return email.Email, nil
		}
	}
	return "", errors.New("email not found")
}


