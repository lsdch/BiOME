package settings

import (
	"context"
	"darco/proto/controllers"
	"darco/proto/models/settings"
	"darco/proto/resolvers"
	"darco/proto/router"
	_ "image/jpeg"
	_ "image/png"
	"net/http"

	"github.com/danielgtaylor/huma/v2"
)

func RegisterRoutes(r router.Router) {
	api := r.RouteGroup("/settings").WithTags([]string{"Settings"})

	router.Register(api, "InstanceSettings",
		huma.Operation{
			Path:    "/instance",
			Method:  http.MethodGet,
			Summary: "Instance settings",
		}, GetInstanceSettings)
	router.Register(api, "UpdateInstanceSettings",
		huma.Operation{
			Path:    "/instance",
			Method:  http.MethodPost,
			Summary: "Update instance settings",
		}, UpdateInstanceSettings)

	router.Register(api, "SecuritySettings",
		huma.Operation{
			Path:    "/security",
			Method:  http.MethodGet,
			Summary: "Security settings",
		}, GetSecuritySettings)
	router.Register(api, "UpdateSecuritySettings",
		huma.Operation{
			Path:    "/security",
			Method:  http.MethodPost,
			Summary: "Update security settings",
		}, UpdateSecuritySettings)

	router.Register(api, "EmailSettings",
		huma.Operation{
			Path:    "/emailing",
			Method:  http.MethodGet,
			Summary: "Email settings",
		}, GetEmailSettings)
	router.Register(api, "UpdateEmailSettings",
		huma.Operation{
			Path:    "/emailing",
			Method:  http.MethodPost,
			Summary: "Update email settings",
		}, UpdateEmailSettings)

	router.Register(api, "ServiceSettings",
		huma.Operation{
			Path:    "/services",
			Method:  http.MethodGet,
			Summary: "Service settings",
		},
		func(ctx context.Context,
			input *resolvers.AccessRestricted[resolvers.Admin],
		) (*controllers.ResponseBody[settings.ServiceSettings], error) {
			return &controllers.ResponseBody[settings.ServiceSettings]{Body: settings.Services()}, nil
		})

	router.Register(api, "UpdateServiceSettings",
		huma.Operation{
			Path:    "/services",
			Method:  http.MethodPatch,
			Summary: "Update service settings",
		}, controllers.CreateHandler[settings.ServiceSettingsUpdate])

	router.Register(api, "TestSMTP",
		huma.Operation{
			Path:    "/emailing/test-dial",
			Method:  http.MethodPost,
			Summary: "Test SMTP connection",
		}, TestSMTP)

	router.Register(api, "SetAppIcon",
		huma.Operation{
			Path:    "/icon",
			Method:  http.MethodPost,
			Summary: "Set app icon",
		}, SetAppIcon)
}

type InstanceSettings struct{ Body settings.InstanceSettings }

func GetInstanceSettings(ctx context.Context, input *struct{}) (*InstanceSettings, error) {
	return &InstanceSettings{Body: settings.Instance()}, nil
}
func UpdateInstanceSettings(ctx context.Context,
	input *struct {
		resolvers.AccessRestricted[resolvers.Admin]
		Body settings.InstanceSettingsInput
	},
) (*InstanceSettings, error) {
	updated, err := input.Body.Save(input.DB())
	return &InstanceSettings{Body: *updated}, err
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
	return &SecuritySettings{Body: *updated}, err
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
