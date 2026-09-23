package services

import (
	"context"
	"errors"
	"fmt"
	"testing"

	"github.com/lsdch/biome/lib/app_errors"
	"github.com/stretchr/testify/require"
)

func TestUnavailableMailer(t *testing.T) {
	cause := errors.New("SMTP connection refused")
	for _, cause := range []error{nil, cause} {
		mailer := NewUnavailableMailer(cause)
		// An unavailable transport must fail before trying to render a template.
		err := mailer.Send(context.Background(), "to@example.org", "from@example.org", "Test", nil)
		require.ErrorIs(t, err, ErrMailerUnavailable)
		if cause != nil {
			require.ErrorIs(t, err, cause)
			require.Contains(t, err.Error(), cause.Error())
		}
		appErr := app_errors.AsAppError(fmt.Errorf("send invitation: %w", err))
		require.Equal(t, 503, appErr.GetStatus())
		require.Equal(t, "Email service is unavailable", appErr.Detail)
	}
}
