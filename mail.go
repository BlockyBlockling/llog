package llog

import (
	"gopkg.in/gomail.v2"
)

var mail_from string
var mail_to []string
var mail_dial gomail.Dialer

func NotifyMail(msg any, a ...any) error {
	mailMessage := gomail.NewMessage()
	message := formatMessage(msg, a...)

	mailMessage.SetHeader("From", mail_from)
	mailMessage.SetHeader("To", mail_to...)
	mailMessage.SetHeader("Subject", "LLog Notification!")
	mailMessage.SetBody("text/plain", message)
	return mail_dial.DialAndSend(mailMessage)
}

func InitMail(from string, to []string, host string, port int, user string, password string, subject string) {
	if subject == "" {
		subject = "LLog Notification"
	}

	mail_dial = *gomail.NewDialer(host, port, user, password)
	mail_from = from
	mail_to = to
}
