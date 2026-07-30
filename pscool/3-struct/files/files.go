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
	return os.ReadFile(filedb.filename)
}

func (filedb *FileDB) Write(data []byte) error {
	return os.WriteFile(filedb.filename, data, 0644)
}

func Read(filename string) ([]byte, error) {
	return os.ReadFile(filename)
}

func Write(data []byte, filename string) error {
	return os.WriteFile(filename, data, 0644)
}