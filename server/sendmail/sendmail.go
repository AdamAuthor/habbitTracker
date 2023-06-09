package sendmail

import (
	"log"
	"math/rand"
	"net/smtp"
	"os"
	"strconv"
)

func SendConfirmationEmail(email, token string) error {

	address := os.Getenv("EMAIL_SENDER_ADDRESS")
	password := os.Getenv("EMAIL_SENDER_PASSWORD")
	content := "Thank you for registering with our service! Please click the link below to confirm your registration:\n\n" +
		"https://habit-makers.herokuapp.com/confirm?token=" + token + "&random=" + strconv.Itoa(rand.Intn(1000))

	// Создаем сообщение для отправки
	msg := []byte("To: " + email + "\r\n" +
		"Subject: Confirm registration\r\n" +
		"\r\n" +
		content)

	// Отправляем сообщение на SMTP сервер
	err := smtp.SendMail("smtp.gmail.com:587",
		smtp.PlainAuth("", address, password, "smtp.gmail.com"),
		address,         // От кого
		[]string{email}, // Кому
		msg)
	if err != nil {
		log.Println("Error sending email:", err)
		return err
	}

	return nil
}

func SendPasswordResetEmail(email, resetLink string) error {
	address := os.Getenv("EMAIL_SENDER_ADDRESS")
	password := os.Getenv("EMAIL_SENDER_PASSWORD")
	content := "We received a password reset request for your account. Please click the link below to reset your password:\n\n" +
		resetLink

	// Создаем сообщение для отправки
	msg := []byte("To: " + email + "\r\n" +
		"Subject: Reset Password\r\n" +
		"\r\n" +
		content)

	// Отправляем сообщение на SMTP сервер
	err := smtp.SendMail("smtp.gmail.com:587",
		smtp.PlainAuth("", address, password, "smtp.gmail.com"),
		address,         // От кого
		[]string{email}, // Кому
		msg)
	if err != nil {
		log.Println("Error sending email:", err)
		return err
	}

	return nil
}
