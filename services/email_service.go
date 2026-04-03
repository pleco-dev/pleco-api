package services

import (
	"fmt"
	"os"

	"github.com/sendgrid/sendgrid-go"
	"github.com/sendgrid/sendgrid-go/helpers/mail"
)

// EmailService defines methods for sending emails (interface)
type EmailService interface {
	SendVerificationEmail(toEmail, token string) error
	SendPasswordReset(toEmail, token string) error
}

// emailService implements EmailService interface
type emailService struct {
	apiKey string
	from   string
}

var _ EmailService = (*emailService)(nil)

// NewEmailService returns a new emailService as EmailService
func NewEmailService() EmailService {
	return &emailService{
		apiKey: os.Getenv("SENDGRID_API_KEY"),
		from:   os.Getenv("SENDGRID_EMAIL"),
	}
}

func (s *emailService) SendVerificationEmail(toEmail, token string) error {
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

func (s *emailService) SendPasswordReset(toEmail string, token string) error {
	link := fmt.Sprintf("http://localhost:8080/reset?token=%s", token)

	from := mail.NewEmail("Go App", s.from)
	subject := "Password Reset Request"
	to := mail.NewEmail("", toEmail)

	plainText := fmt.Sprintf(
		"We received a password reset request for your account. If you did not make this request, you can ignore this email.\n\nTo reset your password, please visit the following link:\n%s\n\nThis link will expire in 15 minutes.",
		link,
	)
	htmlContent := fmt.Sprintf(`
		<strong>Password reset request</strong><br>
		We received a password reset request for your account.<br>
		If you did not make this request, you can ignore this email.<br><br>
		To reset your password, please click the following link:<br>
		<a href="%s">Reset Password</a><br><br>
		This link will expire in 15 minutes.
	`, link)

	message := mail.NewSingleEmail(from, subject, to, plainText, htmlContent)
	client := sendgrid.NewSendClient(s.apiKey)
	response, err := client.Send(message)

	if err != nil {
		return err
	}

	// SendGrid marks success by status 202, errors are 4xx/5xx
	if response != nil && response.StatusCode >= 400 {
		return fmt.Errorf("failed to send email: status=%d, body=%s", response.StatusCode, response.Body)
	}

	return nil
}

// containsIgnoreCase checks if substr is in s (case-insensitive)
func containsIgnoreCase(s, substr string) bool {
	return len(s) >= len(substr) && (len(substr) == 0 || (s != "" && substr != "" &&
		(len(s) == len(substr) && (s == substr || (len(s) > 0 && len(substr) > 0 &&
			containsIgnoreCase(s[1:], substr))))) ||
		(len(s) > 0 && (s[0]|32) == (substr[0]|32) && containsIgnoreCase(s[1:], substr[1:])))
}
