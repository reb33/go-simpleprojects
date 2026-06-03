package service

import (
	"email-verify/configs"
	"email-verify/internal/repository"
	"email-verify/pkg/generate"
	"errors"

	"fmt"
	"net/smtp"

	"github.com/jordan-wright/email"
)

type EmailService struct {
	*configs.Config
	*repository.EmailRepository
}

func NewEmailService(conf *configs.Config, repository *repository.EmailRepository) *EmailService {
	return &EmailService{
		Config: conf,
		EmailRepository: repository,
	}
}

func (service *EmailService) SendEmail(verifyEmail string) error {
	e := email.NewEmail()
	hash := generate.Hash()
	verifyLink := fmt.Sprintf("<a>%v/%v</a>", service.Config.VerifyURL, hash)
	e.HTML = []byte(verifyLink)
	e.From = "asbc@gmail.com"
	e.To = []string{verifyEmail}
	e.Subject = "Verify"

	smtp_server := fmt.Sprintf("%v:%v", service.Config.SendlerConfig.Address, service.Config.SendlerConfig.PORT)
	err := e.Send(
		smtp_server,
		smtp.PlainAuth(
			"",
			service.Config.SendlerConfig.Email,
			service.Config.SendlerConfig.Password,
			service.Config.SendlerConfig.Address,
		),
	)
	if err != nil {
		return err
	}
	err = service.EmailRepository.Add(verifyEmail, hash)
	return err
}

func (service *EmailService) VerifyEmail(hash string) error {
	if hash == "" {
		return errors.New("hash is empty")
	}
	_, err := service.EmailRepository.Get(hash)
	return err
}
