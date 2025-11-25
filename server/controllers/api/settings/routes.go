package settings

import (
	"context"
	"net/http"

	"github.com/danielgtaylor/huma/v2"
	"github.com/lsdch/biome/controllers"
	"github.com/lsdch/biome/models/settings"
	"github.com/lsdch/biome/resolvers"
	"github.com/lsdch/biome/router"
)

func init() {
	settingsAPI := func(r *router.Router) router.Group {
		return r.RouteGroup("/settings").
			WithTags([]string{"Settings"})
	}

	router.RegisterSpec(
		settingsAPI,
		"InstanceSettings",
		huma.Operation{
			Path:    "/instance",
			Method:  http.MethodGet,
			Summary: "Instance settings",
		},
		GetInstanceSettings)

	router.RegisterSpec(
		settingsAPI,
		"UpdateInstanceSettings",
		huma.Operation{
			Path:    "/instance",
			Method:  http.MethodPost,
			Summary: "Update instance settings",
		},
		UpdateInstanceSettings)

	router.RegisterSpec(
		settingsAPI,
		"SecuritySettings",
		huma.Operation{
			Path:    "/security",
			Method:  http.MethodGet,
			Summary: "Security settings",
		},
		GetSecuritySettings)

	router.RegisterSpec(
		settingsAPI,
		"UpdateSecuritySettings",
		huma.Operation{
			Path:    "/security",
			Method:  http.MethodPost,
			Summary: "Update security settings",
		},
		UpdateSecuritySettings)

	router.RegisterSpec(
		settingsAPI,
		"EmailSettings",
		huma.Operation{
			Path:    "/emailing",
			Method:  http.MethodGet,
			Summary: "Email settings",
		},
		GetEmailSettings)

	router.RegisterSpec(
		settingsAPI,
		"UpdateEmailSettings",
		huma.Operation{
			Path:    "/emailing",
			Method:  http.MethodPost,
			Summary: "Update email settings",
		},
		UpdateEmailSettings)

	router.RegisterSpec(
		settingsAPI,
		"ServiceSettings",
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

	router.RegisterSpec(
		settingsAPI,
		"UpdateServiceSettings",
		huma.Operation{
			Path:    "/services",
			Method:  http.MethodPatch,
			Summary: "Update service settings",
		},
		controllers.CreateHandler[settings.ServiceSettingsUpdate])

	router.RegisterSpec(
		settingsAPI,
		"TestSMTP",
		huma.Operation{
			Path:    "/emailing/test-dial",
			Method:  http.MethodPost,
			Summary: "Test SMTP connection",
		},
		TestSMTP)

	router.RegisterSpec(
		settingsAPI,
		"SetAppIcon",
		huma.Operation{
			Path:    "/icon",
			Method:  http.MethodPost,
			Summary: "Set app icon",
		},
		SetAppIcon)
}
