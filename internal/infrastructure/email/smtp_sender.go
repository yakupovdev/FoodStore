package email

import (
	"context"
	"net/smtp"
)

type SMTPSender struct {
	from     string
	password string
	host     string
	port     string
}

func NewSMTPSender(from, password, host, port string) *SMTPSender {
	return &SMTPSender{
		from:     from,
		password: password,
		host:     host,
		port:     port,
	}
}

func (s *SMTPSender) SendRecoveryCode(_ context.Context, emailTo, code string) error {
	auth := smtp.PlainAuth(
		"",
		s.from,
		s.password,
		s.host,
	)

	msg := []byte(
		"From: FoodStore <" + s.from + ">\r\n" +
			"To:" + emailTo + "\r\n" +
			"Subject: Test Email from FoodStore\r\n" +
			"MIME-Version: 1.0\r\n" +
			"Content-UserType: text/plain; charset=\"UTF-8\"\r\n" +
			"\r\n" +
			"Here your code: " + code + "\n",
	)

	return smtp.SendMail(
		s.host+":"+s.port,
		auth,
		s.from,
		[]string{emailTo},
		msg,
	)
}
