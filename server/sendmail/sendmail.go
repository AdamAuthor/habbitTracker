package sendmail

import (
	"log"
	"net/smtp"
)

func SendConfirmationEmail(email, token string) error {
	content := "Thank you for registering with our service! Please click the link below to confirm your registration:\n\n" +
		"http://localhost:8080/confirm?token=" + token
	log.Println(email)

	// Создаем сообщение для отправки
	msg := []byte("To: " + email + "\r\n" +
		"Subject: Confirm registration\r\n" +
		"\r\n" +
		content)

	// pszzdoosimhclkcv
	// Отправляем сообщение на SMTP серверВ
	err := smtp.SendMail("smtp.gmail.com:587",
		smtp.PlainAuth("", "ahmediarolzasov@gmail.com", "pszzdoosimhclkcv", "smtp.gmail.com"),
		"ahmediarolzasov@gmail.com", // От кого
		[]string{email},             // Кому
		msg)
	if err != nil {
		log.Println("Error sending email:", err)
		return err
	}

	log.Println("Email sent to:", email)
	return nil
}
