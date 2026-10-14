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

func TestImportHabitatRoutesAllowContributors(t *testing.T) {
	r := router.New(gin.New(), "", huma.DefaultConfig("test", "1"))
	NewImportController(nil, nil).RegisterRoutes(&r)
	spec := r.API.OpenAPI()
	path := spec.Paths["/imports/batches/{id}/habitats"]
	require.NotNil(t, path)
	require.NotNil(t, path.Get)
	require.NotNil(t, path.Patch)
	require.Equal(t, "GetHabitatsResolution", path.Get.OperationID)
	require.Equal(t, "ResolveHabitats", path.Patch.OperationID)
	for _, op := range []*huma.Operation{path.Get, path.Patch} {
		policy, ok := middleware.GetAccessPolicy(*op)
		require.True(t, ok)
		require.False(t, policy.Eval(nil))
		require.False(t, policy.Eval(&models.User{Role: biomedb.UserRoleVisitor}))
		for _, role := range []biomedb.UserRole{biomedb.UserRoleContributor, biomedb.UserRoleMaintainer, biomedb.UserRoleAdmin} {
			require.True(t, policy.Eval(&models.User{Role: role}))
		}
	}
	schema := spec.Components.Schemas.Map()["SamplingHabitatResolutionInput"]
	require.NotNil(t, schema)
	require.Contains(t, schema.Properties, "habitat_ids")
	require.Contains(t, schema.Properties, "create_habitats")
}
