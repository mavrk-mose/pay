package service

import (
	"fmt"
	"net/smtp"
)

// EmailNotifier sends email notifications via SMTP
type EmailNotifier struct {
	from     string // Sender email address
	smtpHost string // SMTP server host
	smtpPort string // SMTP server port
	auth     smtp.Auth
}

func NewEmailNotifier(from, smtpHost, smtpPort, smtpUser, smtpPass string) *EmailNotifier {
	auth := smtp.PlainAuth("", smtpUser, smtpPass, smtpHost)
	return &EmailNotifier{
		from:     from,
		smtpHost: smtpHost,
		smtpPort: smtpPort,
		auth:     auth,
	}
}

func (n *EmailNotifier) Send(userID, title, message string) error {
	// Fetch the user's email from the database (mock implementation)
	to := "user@example.com" // Replace with the user's email

	// Compose the email
	subject := fmt.Sprintf("Subject: %s\n", title)
	body := fmt.Sprintf("%s\n\n%s", title, message)
	msg := []byte(subject + "\n" + body)

	// Send the email
	err := smtp.SendMail(n.smtpHost+":"+n.smtpPort, n.auth, n.from, []string{to}, msg)
	if err != nil {
		return fmt.Errorf("failed to send email: %v", err)
	}

	return nil
}
