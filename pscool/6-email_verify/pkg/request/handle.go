package request

import (
	"net/http"
)

func HandleBody[T any](r *http.Request) (*T, error) {
	body, err := DecodeRequest[T](r.Body)
	if err != nil {
		return nil, err
	}
	err = IsValid(body)
	if err != nil {
		return nil, err
	}
	return &body, nil
}
