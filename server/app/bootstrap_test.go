package app

import (
	"context"
	"testing"

	"github.com/lsdch/biome/config"
	"github.com/lsdch/biome/services"
	"github.com/stretchr/testify/require"
)

func TestEmailServiceRequiresPersistedSettings(t *testing.T) {
	// Environment settings are only a bootstrap seed, never the active transport.
	email := services.NewEmailService(config.MailTransportConfig{SMTPHost: "smtp.example.org", SMTPPort: 587})
	appServices := &AppServices{EmailService: email}
	err := appServices.EmailService.Send(context.Background(), "to@example.org", "Test", nil)
	require.ErrorIs(t, err, services.ErrMailerUnavailable)
	require.Empty(t, appServices.EmailService.GetSettings())
}
