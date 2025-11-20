package libs

import (
	"fmt"
	"log"
	"net/smtp"
)

func SendOTPEmail(toEmail, otp string) error {
	from := "youremail@gmail.com"
	password := "your-app-password"

	smtpHost := "smtp.gmail.com"
	smtpPort := "587"

	subject := "Your OTP Code"
	body := fmt.Sprintf("Your OTP code is: %s\nIt will expire in 2 minutes.", otp)

	msg := "From: " + from + "\r\n" +
		"To: " + toEmail + "\r\n" +
		"Subject: " + subject + "\r\n\r\n" +
		body

	auth := smtp.PlainAuth("", from, password, smtpHost)

	err := smtp.SendMail(smtpHost+":"+smtpPort, auth, from, []string{toEmail}, []byte(msg))
	if err != nil {
		log.Println("Failed to send OTP email:", err)
		return fmt.Errorf("failed to send OTP email: %w", err)
	}

	log.Println("OTP email sent successfully to", toEmail)
	return nil
}
