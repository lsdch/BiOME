package services

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"net/http"
	"net/mail"
	"sync/atomic"

	"github.com/a-h/templ"
	"github.com/danielgtaylor/huma/v2"
	"github.com/jackc/pgx/v5"
	"github.com/k3a/html2text"
	"github.com/lsdch/biome/config"
	"github.com/lsdch/biome/db"
	"github.com/lsdch/biome/lib/app_errors"
	"github.com/lsdch/biome/models"
	"gopkg.in/gomail.v2"
)

type EmailService struct {
	seed        config.MailTransportConfig
	settings    atomic.Pointer[models.Mailing]
	dialSuccess atomic.Bool
}

func NewEmailService(seed config.MailTransportConfig) *EmailService {
	return &EmailService{seed: seed}
}

func (s *EmailService) Available() bool {
	return s.GetSettings().Enabled && s.dialSuccess.Load()
}

func (s *EmailService) Toggle(ctx context.Context, q db.Querier, enabled bool) error {
	err := q.Queries().ToggleMailing(ctx, enabled)
	if err != nil {
		return err
	}
	return s.Reload(ctx, q)
}

// Bootstrap initializes mailing once; persisted settings always take precedence.
func (s *EmailService) Bootstrap(ctx context.Context, q db.Querier) error {
	err := q.WithTx(ctx, func(tx *db.Tx) error {
		if _, err := tx.Queries().GetMailing(ctx); err == nil {
			// already exists, nothing to do
			return nil
		} else if !errors.Is(err, pgx.ErrNoRows) {
			return err
		}
		smtp := s.seed
		if err := smtp.Validate(); err != nil {
			return fmt.Errorf("invalid initial mailing configuration: %w", err)
		}
		port := int32(smtp.SMTPPort)
		settings := models.UpsertMailingParams{
			MailFromAddress: smtp.MailFromAddress,
			MailFromName:    smtp.MailFromName,
			SmtpHost:        &smtp.SMTPHost,
			SmtpPort:        &port,
			SmtpUser:        &smtp.SMTPUser,
			SmtpPassword:    &smtp.SMTPPassword,
		}
		return s.SaveSettings(ctx, q, settings)
	})
	if err != nil {
		return fmt.Errorf("failed to bootstrap mailing: %w", err)
	}
	return s.Reload(ctx, q)
}

// GetSettings returns an independent snapshot, including nullable SMTP values.
func (s *EmailService) GetSettings() models.Mailing {
	settings := s.settings.Load()
	if settings == nil {
		return models.Mailing{}
	}
	return *settings
}

func (s *EmailService) Reload(ctx context.Context, q db.Querier) error {
	row, err := q.Queries().GetMailing(ctx)
	if err != nil {
		return fmt.Errorf("reload mailing settings: %w", err)
	}
	settings := models.MailingFromDB(row)
	s.settings.Store(&settings)
	dialSuccess, err := s.TestConnection(ctx)
	if err != nil {
		s.dialSuccess.Store(false)
		return fmt.Errorf("test mailer connection: %w", err)
	}
	s.dialSuccess.Store(dialSuccess)
	return nil
}

func (s *EmailService) SaveSettings(ctx context.Context, q db.Querier, input models.UpsertMailingParams) error {
	if _, err := q.Queries().UpsertMailing(ctx, input.ToDBParams()); err != nil {
		return err
	}
	return s.Reload(ctx, q)
}

func (s *EmailService) Send(ctx context.Context, to string, subject string, template templ.Component) (err error) {
	mailing := s.GetSettings()
	dialer, err := mailingDialer(mailing)
	if err != nil {
		return err
	}
	if available := s.dialSuccess.Load(); !available {
		return ErrMailerUnavailable
	}
	var body bytes.Buffer
	err = template.Render(ctx, &body)
	if err != nil {
		return err
	}

	m := gomail.NewMessage()
	from := (&mail.Address{Name: mailing.MailFromName, Address: mailing.MailFromAddress}).String()
	m.SetHeader("From", from)
	m.SetHeader("To", to)
	m.SetHeader("Subject", subject)
	m.SetBody("text/plain", html2text.HTML2Text(body.String()))
	m.AddAlternative("text/html", body.String())

	if err := ctx.Err(); err != nil {
		return err
	}
	return dialer.DialAndSend(m)
}

func (s *EmailService) Dialer() (*gomail.Dialer, error) {
	return mailingDialer(s.GetSettings())
}

func mailingDialer(settings models.Mailing) (*gomail.Dialer, error) {
	if !settings.IsConfigured() {
		return nil, fmt.Errorf("%w: SMTP host and port must be configured", ErrMailerUnavailable)
	}
	var user, password string
	if settings.SmtpUser != nil {
		user = *settings.SmtpUser
	}
	if settings.SmtpPassword != nil {
		password = *settings.SmtpPassword
	}
	return gomail.NewDialer(*settings.SmtpHost, int(*settings.SmtpPort), user, password), nil
}

func (s *EmailService) TestConnection(ctx context.Context) (bool, error) {
	dialer, err := s.Dialer()
	if err != nil {
		return false, err
	}
	if err := ctx.Err(); err != nil {
		return false, err
	}
	closer, err := dialer.Dial()
	if err != nil {
		return false, err
	}
	if err := closer.Close(); err != nil {
		return false, err
	}
	s.dialSuccess.Store(true)
	return true, nil
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
