package auth

import (
	"demo-store/pkg/jwt"
	"demo-store/pkg/random"
)

const CODE = "123456"
const SESSION_ID_LENGTH = 16

type Service struct {
	jwt   *jwt.JWT
	repo *AuthRepository
}

func NewService(jwt *jwt.JWT, repo *AuthRepository) *Service {
	return &Service{
		jwt:   jwt,
		repo: repo,
	}
}

func (s *Service) AddPhone(phone string) (string, error) {
	data, err := s.repo.Upsert(phone, random.GenSessionId(SESSION_ID_LENGTH), CODE)
	if err != nil {
		return "", err
	}
	return data.SessionId, nil
}

func (s *Service) Verify(sessionId, code string) (string, error) {
	data, _ := s.repo.GetBySessionId(sessionId)
	if data == nil {
		return "", ErrInvalidSessionId
	}
	if data.Code != code {
		return "", ErrInvalidCode
	}
	token, err := s.jwt.Create(data.Phone)
	if err != nil {
		return "", err
	}
	return token, nil
}