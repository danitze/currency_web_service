package service

import (
	"crypto/tls"
	"gopkg.in/gomail.v2"
	"os"
	"strconv"
)

func SendToEmails(emails []string, subject string, content string) error {
	message := gomail.NewMessage()
	message.SetHeader("From", os.Getenv("EMAIL_FROM"))
	message.SetHeader("To", emails...)
	message.SetHeader("Subject", subject)
	message.SetBody("text/plain", content)
	port, err := strconv.Atoi(os.Getenv("EMAIL_PORT"))
	if err != nil {
		return err
	}
	dialer := gomail.NewDialer(os.Getenv("EMAIL_HOST"), port, os.Getenv("EMAIL_FROM"), os.Getenv("EMAIL_PASSWORD"))
	dialer.TLSConfig = &tls.Config{InsecureSkipVerify: true}
	return dialer.DialAndSend(message)
}
