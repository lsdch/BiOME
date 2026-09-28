package services

import (
	"testing"

	"github.com/lsdch/biome/config"
	"github.com/lsdch/biome/models"
	"github.com/stretchr/testify/require"
)

func TestEmailServiceCachedSettings(t *testing.T) {
	service := NewEmailService(config.MailTransportConfig{})
	_, err := service.Dialer()
	require.ErrorIs(t, err, ErrMailerUnavailable)
	host, password := "db.example.org", "secret"
	port := int32(2525)
	settings := models.Mailing{SmtpHost: &host, SmtpPort: &port, SmtpPassword: &password}
	service.settings.Store(&settings)
	snapshot := service.GetSettings()
	*snapshot.SmtpHost = "mutated.example.org"
	*snapshot.SmtpPassword = "mutated"
	dialer, err := service.Dialer()
	require.NoError(t, err)
	require.Equal(t, host, dialer.Host)
	require.Equal(t, int(port), dialer.Port)
	require.Equal(t, password, dialer.Password)
	next := service.GetSettings()
	nextHost := "updated.example.org"
	next.SmtpHost = &nextHost
	next.SmtpPassword = nil
	service.settings.Store(&next)
	dialer, err = service.Dialer()
	require.NoError(t, err)
	require.Equal(t, nextHost, dialer.Host)
	require.Empty(t, dialer.Password)
}
