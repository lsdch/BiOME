package config

import (
	"github.com/stretchr/testify/require"
	"testing"
)

func TestMailTransportSeedValidation(t *testing.T) {
	seed := MailTransportConfig{MailFromAddress: "sender@example.org", MailFromName: "BiOME", SMTPHost: "smtp.example.org", SMTPPort: 587, SMTPUser: "user", SMTPPassword: "secret"}
	require.NoError(t, seed.Validate())
	seed.SMTPPort = 65536
	require.Error(t, seed.Validate())
	require.Error(t, (MailTransportConfig{}).Validate())
}
