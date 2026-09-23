package app

import (
	"context"
	"errors"
	"testing"

	"github.com/lsdch/biome/config"
	"github.com/lsdch/biome/services"
	"github.com/stretchr/testify/require"
)

func TestBootstrapMailer(t *testing.T) {
	cause := errors.New("SMTP connection refused")
	for _, tc := range []struct {
		name string
		ok   bool
		err  error
	}{
		{"available", true, nil},
		{"connection failure", false, cause},
		{"negative status", false, nil},
	} {
		t.Run(tc.name, func(t *testing.T) {
			cfg := config.Config{SMTP: config.SMTPConfig{SMTPHost: "smtp.example.org", SMTPPort: 587}}
			appServices := &AppServices{SettingsService: services.NewSettingsService(cfg)}
			bootstrap := NewAppBootstrap(nil, config.BootstrapConfig{}, appServices)
			ctx := context.Background()
			calls := 0
			require.NotPanics(t, func() {
				bootstrap.bootstrapMailer(ctx, func(got context.Context) (bool, error) {
					require.Equal(t, ctx, got)
					calls++
					return tc.ok, tc.err
				})
			})
			require.Equal(t, 1, calls)
			require.NotNil(t, appServices.Mailer)
			if tc.ok {
				smtp, ok := appServices.Mailer.(*services.EmailService)
				require.True(t, ok)
				require.Equal(t, cfg.SMTP.SMTPHost, smtp.Dialer().Host)
				require.Equal(t, int(cfg.SMTP.SMTPPort), smtp.Dialer().Port)
			} else {
				err := appServices.Mailer.Send(ctx, "to@example.org", "from@example.org", "Test", nil)
				require.ErrorIs(t, err, services.ErrMailerUnavailable)
				if tc.err != nil {
					require.ErrorIs(t, err, tc.err)
				}
			}
		})
	}
}
