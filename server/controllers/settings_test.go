package controllers

import (
	"testing"

	"github.com/danielgtaylor/huma/v2"
	"github.com/gin-gonic/gin"
	"github.com/lsdch/biome/db/biomedb"
	"github.com/lsdch/biome/middleware"
	"github.com/lsdch/biome/models"
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

func TestServiceSettingsRoutes(t *testing.T) {
	r := router.New(gin.New(), "", huma.DefaultConfig("test", "1"))
	controller := NewSettingsController(nil, nil, nil)
	controller.RegisterRoutes(&r)
	spec := r.API.OpenAPI()
	path := spec.Paths["/settings/services"]
	require.NotNil(t, path)
	for _, op := range []*huma.Operation{path.Get, path.Put} {
		require.NotNil(t, op)
		policy, ok := middleware.GetAccessPolicy(*op)
		require.True(t, ok)
		require.False(t, policy.Eval(nil))
		require.False(t, policy.Eval(&models.User{Role: biomedb.UserRoleVisitor}))
		require.True(t, policy.Eval(&models.User{Role: biomedb.UserRoleAdmin}))
	}
	require.Equal(t, "GetServiceSettings", path.Get.OperationID)
	require.Equal(t, "UpdateServiceSettings", path.Put.OperationID)
	schema := spec.Components.Schemas.Map()["ServiceSettings"]
	require.NotNil(t, schema)
	require.Contains(t, schema.Properties, "geoapify_api_key")
	require.True(t, schema.Properties["geoapify_api_key"].Nullable)
	require.Equal(t, "#/components/schemas/ServiceSettings", path.Put.RequestBody.Content["application/json"].Schema.Ref)
}
