package files

import (
	"errors"
	"os"
	"strings"
)

type FileDB struct {
	filename string
}

func NewFileDB(name string) (*FileDB, error) {
	if !strings.HasSuffix(name, ".json") {
		return nil, errors.New("файл должен быть josn")
	}
	return &FileDB{
		filename: name,
	}, nil
}

func (filedb *FileDB) Read() ([]byte, error) {
	data, err := os.ReadFile(filedb.filename)
	if err != nil {
		return nil, err
	}
	return data, nil
}

func (filedb *FileDB) Write(data []byte) error {
	return os.WriteFile(filedb.filename, data, 0644)
}