package services

import (
	"fmt"
	"net/smtp"
	"os"
)

type Email struct {
	Subject string
	Body    string
}

type EmailService interface {
	Send(email Email) error
}

type emailService struct {
}

func NewEmailService() EmailService {
	return &emailService{}
}

func (s *emailService) Send(email Email) error {

	// Sender data.
	from := os.Getenv("EMAIL_LOGIN")
	password := os.Getenv("EMAIL_PASSWORD")

	// Receiver email address.
	to := os.Getenv("EMAIL_LOGIN")

	// smtp server configuration.
	smtpHost := os.Getenv("SMTP_HOST")
	smtpPort := os.Getenv("SMTP_PORT")

	// Email subject and body
	subject := email.Subject
	body := email.Body

	// Message construction
	message := fmt.Sprintf("To: %s\r\nSubject: %s\r\n\r\n%s", to, subject, body)

	// Authentication.
	auth := smtp.PlainAuth("", from, password, smtpHost)

	// Sending email.
	err := smtp.SendMail(smtpHost+":"+smtpPort, auth, from, []string{to}, []byte(message))
	if err != nil {
		logger.Err(err).Send()
		return err
	}

	return nil
}
