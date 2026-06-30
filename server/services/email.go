package services

import (
	"bytes"
	"context"

	"github.com/a-h/templ"
	"github.com/k3a/html2text"
	"github.com/lsdch/biome/config"
	"gopkg.in/gomail.v2"
)

type EmailService struct {
	config config.SMTPConfig
}

func (s *EmailService) Send(ctx context.Context, to string, from string, subject string, template templ.Component) (err error) {
	var body bytes.Buffer
	err = template.Render(ctx, &body)
	if err != nil {
		return err
	}

	m := gomail.NewMessage()
	m.SetHeader("From", from)
	m.SetHeader("To", to)
	m.SetHeader("Subject", subject)
	m.SetBody("text/html", body.String())
	m.AddAlternative("text/plain", html2text.HTML2Text(body.String()))

	dialer := s.Dialer()

	return dialer.DialAndSend(m)
}

func (s *EmailService) Dialer() *gomail.Dialer {
	return gomail.NewDialer(
		s.config.SMTPHost,
		int(s.config.SMTPPort),
		s.config.SMTPUser,
		s.config.SMTPPassword,
	)
}
