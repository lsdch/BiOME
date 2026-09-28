package controllers

import (
	"testing"

	"github.com/danielgtaylor/huma/v2"
	"github.com/gin-gonic/gin"
	"github.com/lsdch/biome/router"
	"github.com/stretchr/testify/require"
)

func TestSettingsRoutesSeparateInstanceAndEmail(t *testing.T) {
	r := router.New(gin.New(), "", huma.DefaultConfig("test", "1"))
	controller := NewSettingsController(nil, nil, nil)
	controller.RegisterRoutes(&r)
	spec := r.API.OpenAPI()
	require.NotNil(t, spec.Paths["/settings/instance"].Get)
	require.NotNil(t, spec.Paths["/settings/instance"].Put)
	require.NotNil(t, spec.Paths["/settings/email"].Get)
	require.NotNil(t, spec.Paths["/settings/email"].Put)
	require.NotNil(t, spec.Paths["/settings/email/smtp/test"].Get)
	require.NotContains(t, spec.Paths, "/settings/smtp/test")
	schemas := spec.Components.Schemas.Map()
	for _, name := range []string{"InstanceSettings", "InstanceSettingsUpdate"} {
		require.NotContains(t, schemas[name].Properties, "mail_from_address")
		require.NotContains(t, schemas[name].Properties, "mail_from_name")
	}
	require.Contains(t, schemas["Mailing"].Properties, "smtp_host")
	require.Contains(t, schemas["UpsertMailingParams"].Properties, "smtp_host")
}
