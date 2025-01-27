package email

import (
	"bytes"
	"context"

	"github.com/lsdch/biome/models/settings"

	"github.com/a-h/templ"
	"github.com/k3a/html2text"
	"gopkg.in/gomail.v2"
)

type EmailData struct {
	// Email address of the recipient
	To string
	// Email 'From' header
	From string
	// Email subject
	Subject string
	// Email template
	Template templ.Component
}

func (emailData *EmailData) Body() (*bytes.Buffer, error) {
	var body bytes.Buffer
	err := emailData.Template.Render(context.Background(), &body)
	return &body, err
}

func (email *EmailData) Send(from string) (err error) {
	body, err := email.Body()
	if err != nil {
		return err
	}

	m := gomail.NewMessage()
	m.SetHeader("From", from)
	m.SetHeader("To", email.To)
	m.SetHeader("Subject", email.Subject)
	m.SetBody("text/html", body.String())
	m.AddAlternative("text/plain", html2text.HTML2Text(body.String()))

	settings := settings.Email()
	dialer := settings.Dialer()
	// dialer.TLSConfig = &tls.Config{InsecureSkipVerify: true}

	return dialer.DialAndSend(m)
}
