package email

import (
	"context"
	"net/smtp"
)

type SMTPSender struct {
	From     string `envconfig:"SENDER" required:"true"`
	Password string `envconfig:"PASSWORD" required:"true"`
	Host     string `envconfig:"HOST" required:"true"`
	Port     string `envconfig:"PORT" required:"true"`
}

func NewSMTPSender(from, password, host, port string) *SMTPSender {
	return &SMTPSender{
		From:     from,
		Password: password,
		Host:     host,
		Port:     port,
	}
}

func (s *SMTPSender) SendRecoveryCode(_ context.Context, emailTo, code string) error {
	auth := smtp.PlainAuth(
		"",
		s.From,
		s.Password,
		s.Host,
	)

	msg := []byte(
		"From: FoodStore <" + s.From + ">\r\n" +
			"To:" + emailTo + "\r\n" +
			"Subject: Test Email From FoodStore\r\n" +
			"MIME-Version: 1.0\r\n" +
			"Content-UserType: text/plain; charset=\"UTF-8\"\r\n" +
			"\r\n" +
			"Here your code: " + code + "\n",
	)

	return smtp.SendMail(
		s.Host+":"+s.Port,
		auth,
		s.From,
		[]string{emailTo},
		msg,
	)
}

func (s *SMTPSender) SendMessage(emailTo, message string) error {
	auth := smtp.PlainAuth(
		"",
		s.From,
		s.Password,
		s.Host,
	)

	msg := []byte(
		"From: FoodStore <" + s.From + ">\r\n" +
			"To:" + emailTo + "\r\n" +
			"Subject: Test Email From FoodStore\r\n" +
			"MIME-Version: 1.0\r\n" +
			"Content-UserType: text/plain; charset=\"UTF-8\"\r\n" +
			"\r\n" + message + "\n",
	)

	return smtp.SendMail(
		s.Host+":"+s.Port,
		auth,
		s.From,
		[]string{emailTo},
		msg,
	)
}
