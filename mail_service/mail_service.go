package mail_service

import (
	"log"
	"taskel/config"

	"gopkg.in/gomail.v2"
)

const CONFIG_SMTP_HOST = "smtp.gmail.com"
const CONFIG_SMTP_PORT = 587
const CONFIG_SENDER_NAME = "Taskel <anggonohanggi@gmail.com>"
const CONFIG_AUTH_EMAIL = "anggonohanggi@gmail.com"

var CONFIG_AUTH_PASSWORD = config.Config.EmailPassword

func SendMail(subject string, body string, recipients ...string) {
	mailer := gomail.NewMessage()
	mailer.SetHeader("From", CONFIG_SENDER_NAME)
	mailer.SetHeader("To", recipients...)
	mailer.SetHeader("Subject", subject)
	mailer.SetBody("text/html", body)

	dialer := gomail.NewDialer(
		CONFIG_SMTP_HOST,
		CONFIG_SMTP_PORT,
		CONFIG_AUTH_EMAIL,
		CONFIG_AUTH_PASSWORD,
	)

	err := dialer.DialAndSend(mailer)
	if err != nil {
		log.Fatalf("EMAIL: failed to send %v \n", err)
	}

	log.Println("EMAIL: Email sent")
}
