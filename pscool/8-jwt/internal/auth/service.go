package auth

import (
	"some-project/pkg/jwt"
	"some-project/pkg/random"
)

const CODE = "123456"
const SESSION_ID_LENGTH = 16

type Service struct {
	jwt   *jwt.JWT
	store *Store
}

func NewService(jwt *jwt.JWT) *Service {
	return &Service{
		jwt:   jwt,
		store: NewStore(),
	}
}

func (s *Service) AddPhone(phone string) string {
	data := s.store.Upsert(phone, random.GenSessionId(SESSION_ID_LENGTH), CODE)
	return data.sessionId
}

func (s *Service) Verify(sessionId, code string) (string, error) {
	data := s.store.GetBySessionId(sessionId)
	if data == nil {
		return "", ErrInvalidSessionId
	}
	if data.code != code {
		return "", ErrInvalidCode
	}
	token, err := s.jwt.Create(data.phone)
	if err != nil {
		return "", err
	}
	return token, nil
}