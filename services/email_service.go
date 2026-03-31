package services

import (
	"fmt"
	"os"

	"github.com/sendgrid/sendgrid-go"
	"github.com/sendgrid/sendgrid-go/helpers/mail"
)

type EmailService struct {
	apiKey string
	from   string
}

func NewEmailService() *EmailService {
	return &EmailService{
		apiKey: os.Getenv("SENDGRID_API_KEY"),
		from:   os.Getenv("SENDGRID_EMAIL"), // wajib verified di SendGrid
	}
}

func (s *EmailService) SendVerificationEmail(toEmail, token string) error {
	link := fmt.Sprintf("http://localhost:8080/verify?token=%s", token)

	from := mail.NewEmail("Go App", s.from)
	subject := "Verify Your Email"
	to := mail.NewEmail("", toEmail)

	plainText := fmt.Sprintf("Click this link to verify: %s", link)
	htmlContent := fmt.Sprintf(`
		<strong>Verify your account</strong><br>
		Click here: <a href="%s">Verify Email</a>
	`, link)

	message := mail.NewSingleEmail(from, subject, to, plainText, htmlContent)

	client := sendgrid.NewSendClient(s.apiKey)
	_, err := client.Send(message)

	return err
}
