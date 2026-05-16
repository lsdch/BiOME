package settings

import (
	"context"
	_ "image/jpeg"
	_ "image/png"

	"github.com/lsdch/biome/models/settings"
	"github.com/lsdch/biome/resolvers"

	"github.com/danielgtaylor/huma/v2"
)

type InstanceSettings struct{ Body settings.InstanceSettings }

func GetInstanceSettings(ctx context.Context, input *struct{}) (*InstanceSettings, error) {
	return &InstanceSettings{Body: settings.Instance()}, nil
}
func UpdateInstanceSettings(ctx context.Context,
	input *struct {
		resolvers.AccessRestricted[resolvers.Admin]
		Body settings.InstanceSettingsUpdate
	},
) (*InstanceSettings, error) {
	updated, err := input.Body.Save(input.DB())
	if err != nil {
		return nil, err
	}
	return &InstanceSettings{Body: updated}, nil
}

type SecuritySettings struct{ Body settings.SecuritySettings }

func GetSecuritySettings(ctx context.Context, input *resolvers.AccessRestricted[resolvers.Admin]) (*SecuritySettings, error) {
	return &SecuritySettings{Body: settings.Security()}, nil
}

func UpdateSecuritySettings(ctx context.Context,
	input *struct {
		resolvers.AccessRestricted[resolvers.Admin]
		Body settings.SecuritySettingsInput
	},
) (*SecuritySettings, error) {
	updated, err := input.Body.Save(input.DB())
	if err != nil {
		return nil, err
	}
	return &SecuritySettings{Body: *updated}, nil
}

type EmailSettings struct{ Body settings.EmailSettings }

func GetEmailSettings(ctx context.Context, input *resolvers.AccessRestricted[resolvers.Admin]) (*EmailSettings, error) {
	return &EmailSettings{Body: settings.Email()}, nil
}

type EmailSettingsInput struct {
	resolvers.AccessRestricted[resolvers.Admin]
	Body settings.EmailSettingsInput
}

type SMTPConnectionStatus struct{ Body bool }

func UpdateEmailSettings(ctx context.Context, input *EmailSettingsInput) (*EmailSettings, error) {
	if status, err := TestSMTP(ctx, input); !status.Body {
		return nil, err
	}
	updated, err := input.Body.Save(input.DB())
	return &EmailSettings{Body: *updated}, err
}

func TestSMTP(ctx context.Context, input *EmailSettingsInput) (*SMTPConnectionStatus, error) {
	if err := input.Body.TestConnection(); err != nil {
		return &SMTPConnectionStatus{false}, huma.Error422UnprocessableEntity("SMTP connection failed", err)
	}
	return &SMTPConnectionStatus{true}, nil
}
