package request

import (
	"encoding/json"
	"io"
	"net/http"

	"github.com/go-playground/validator/v10"
)

func HandleBody[T any](r *http.Request) (*T, error) {
	body, err := DecodeRequest[T](r.Body) //декодирование
	if err != nil {
		return nil, err
	}
	err = IsValid(body) // запуск валидации
	if err != nil {
		return nil, err
	}
	return &body, nil
}

func DecodeRequest[T any](body io.ReadCloser) (T, error) {
	var payload T
	err := json.NewDecoder(body).Decode(&payload)
	if err != nil {
		return payload, err
	}
	return payload, nil
}

func IsValid[T any](payload T) error {
	validate := validator.New()
	err := validate.Struct(payload)
	return err
}
