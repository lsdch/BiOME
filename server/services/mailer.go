package services

import (
	"bytes"
	"context"
	"fmt"
	"net/http"

	"github.com/a-h/templ"
	"github.com/danielgtaylor/huma/v2"
	"github.com/k3a/html2text"
	"github.com/lsdch/biome/config"
	"github.com/lsdch/biome/lib/app_errors"
	"gopkg.in/gomail.v2"
)

type Mailer interface {
	Send(ctx context.Context, to, from, subject string, template templ.Component) error
}

type EmailService struct {
	config config.SMTPConfig
}

func NewEmailService(config config.SMTPConfig) *EmailService {
	return &EmailService{config: config}
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
	m.SetBody("text/plain", html2text.HTML2Text(body.String()))
	m.AddAlternative("text/html", body.String())

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

// ErrMailerUnavailable is recognized by errors.Is and the common HTTP error mapper.
var ErrMailerUnavailable error = mailerUnavailableError{}

type mailerUnavailableError struct{}

func (mailerUnavailableError) Error() string { return "mailer unavailable" }

func (mailerUnavailableError) AppError() *app_errors.AppError {
	return &app_errors.AppError{ErrorModel: huma.ErrorModel{
		Status: http.StatusServiceUnavailable,
		Title:  "Service Unavailable",
		Detail: "Email service is unavailable",
	}}
}

type UnavailableMailer struct {
	cause error
}

func NewUnavailableMailer(cause error) *UnavailableMailer {
	return &UnavailableMailer{cause: cause}
}

func (m *UnavailableMailer) Send(ctx context.Context, to, from, subject string, template templ.Component) error {
	if m.cause != nil {
		return fmt.Errorf("%w: %w", ErrMailerUnavailable, m.cause)
	}
	return ErrMailerUnavailable
}

var (
	_ Mailer = (*EmailService)(nil)
	_ Mailer = (*UnavailableMailer)(nil)
)
