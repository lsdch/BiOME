package email

import (
	"bytes"
	"crypto/tls"
	"darco/proto/config"
	"html/template"
	"sync"

	"github.com/k3a/html2text"
	"gopkg.in/gomail.v2"
)

type EmailData struct {
	To       string
	Subject  string
	Template string
	Data     interface{}
}

func (emailData *EmailData) Body() (*bytes.Buffer, error) {
	var body bytes.Buffer
	err := templates.ExecuteTemplate(&body, emailData.Template, &emailData.Data)
	if err != nil {
		return nil, err
	}
	return &body, nil
}

var (
	once      sync.Once
	templates template.Template
)

func LoadTemplates(pattern string) (err error) {
	once.Do(func() {
		parsedTemplates, tplErr := template.ParseGlob(pattern)
		if tplErr != nil {
			err = tplErr
			return
		}
		templates = *parsedTemplates
	})
	return
}

func Send(config *config.EmailConfig, email *EmailData) (err error) {
	body, err := email.Body()
	if err != nil {
		return err
	}
	m := gomail.NewMessage()
	m.SetHeader("From", config.From)
	m.SetHeader("To", email.To)
	m.SetHeader("Subject", email.Subject)
	m.SetBody("text/html", body.String())
	m.AddAlternative("text/plain", html2text.HTML2Text(body.String()))

	dialer := gomail.NewDialer(config.Host, config.Port, config.User, config.Pass)
	dialer.TLSConfig = &tls.Config{InsecureSkipVerify: true}

	return dialer.DialAndSend(m)
}
