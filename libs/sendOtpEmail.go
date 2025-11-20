package libs


import (
    "fmt"
    "net/smtp"
)


func SendOTPEmail(toEmail, otp string) error {
    from := "youremail@gmail.com"        
    password := "your-app-password"      
    smtpHost := "smtp.gmail.com"
    smtpPort := "587"

    subject := "Your OTP Code"
    body := fmt.Sprintf("Your OTP code is: %s. It will expire in 2 minutes.", otp)

    msg := "From: " + from + "\n" +
        "To: " + toEmail + "\n" +
        "Subject: " + subject + "\n\n" +
        body

    auth := smtp.PlainAuth("", from, password, smtpHost)

    return smtp.SendMail(smtpHost+":"+smtpPort, auth, from, []string{toEmail}, []byte(msg))
}
