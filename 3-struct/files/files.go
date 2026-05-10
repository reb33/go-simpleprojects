package files

import (
	"errors"
	"os"
	"strings"
)

func ReadFile(name string) ([]byte, error) {
	if !strings.HasSuffix(name, ".json") {
		return nil, errors.New("файл должен быть josn")
	}
	data, err := os.ReadFile(name)
	if err != nil {
		return nil, err
	}
	return data, nil
}

func WriteFile(name string, data []byte) error {
	if !strings.HasSuffix(name, ".json") {
		return errors.New("файл должен быть josn")
	}
	return os.WriteFile(name, data, 0644)
}