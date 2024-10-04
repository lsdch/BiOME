package email

import (
	"bytes"
	"darco/proto/models/settings"
	"html/template"
	"sync"

	"github.com/k3a/html2text"
	"github.com/sirupsen/logrus"
	"gopkg.in/gomail.v2"
)

type EmailData struct {
	// Email address of the recipient
	To      string
	Subject string
	// Template file name for the body of the email
	Template string
	// Template variables
	Data map[string]any
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

func AdminEmailAddress() string {
	// DarCo Admin <admin@darco.instance.edu>
	logrus.Fatalf("Not implemented")
	return ""
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
