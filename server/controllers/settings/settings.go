package settings

import (
	"context"
	"darco/proto/models/settings"
	"darco/proto/resolvers"
	"darco/proto/router"
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

	router.Register(api, "SecuritySettings",
		huma.Operation{
			Path:    "/security",
			Method:  http.MethodGet,
			Summary: "Security settings",
		}, GetSecuritySettings)

	router.Register(api, "EmailSettings",
		huma.Operation{
			Path:    "/emailing",
			Method:  http.MethodGet,
			Summary: "Email settings",
		}, GetEmailSettings)
}

type GetInstanceSettingsOutput struct{ Body settings.InstanceSettings }

func GetInstanceSettings(ctx context.Context, input *struct{}) (*GetInstanceSettingsOutput, error) {
	return &GetInstanceSettingsOutput{Body: settings.Instance()}, nil
}

type GetSecuritySettingsOutput struct{ Body settings.SecuritySettings }

func GetSecuritySettings(ctx context.Context, input *resolvers.AccessRestricted[resolvers.Admin]) (*GetSecuritySettingsOutput, error) {
	return &GetSecuritySettingsOutput{Body: settings.Security()}, nil
}

type GetEmailSettingsOutput struct{ Body settings.EmailSettings }

func GetEmailSettings(ctx context.Context, input *resolvers.AccessRestricted[resolvers.Admin]) (*GetEmailSettingsOutput, error) {
	return &GetEmailSettingsOutput{Body: settings.Email()}, nil
}
