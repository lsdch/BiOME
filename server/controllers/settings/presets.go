package settings

import (
	"net/http"

	"github.com/danielgtaylor/huma/v2"
	"github.com/lsdch/biome/controllers"
	"github.com/lsdch/biome/db"
	"github.com/lsdch/biome/models/people"
	"github.com/lsdch/biome/models/presets"
	"github.com/lsdch/biome/resolvers"
	"github.com/lsdch/biome/router"
)

func RegisterMapPresetsRoutes(r router.Group) {
	api := r.RouteGroup("/mapping").WithTags([]string{"Settings"})

	router.Register(api, "ListDataFeeds",
		huma.Operation{
			Path:    "/data-feeds",
			Method:  http.MethodGet,
			Summary: "List saved data feeds",
		}, controllers.ListHandler[*struct {
			resolvers.AuthResolver
		}](presets.ListDataFeedSpecs))

	router.Register(api, "CreateDataFeed",
		huma.Operation{
			Path:    "/data-feeds",
			Method:  http.MethodPost,
			Summary: "Save data feed",
		}, controllers.CreateHandler[presets.DataFeedSpecInput])

	router.Register(api, "ListMapPresets",
		huma.Operation{
			Path:    "/map-presets",
			Method:  http.MethodGet,
			Summary: "List saved map presets",
		}, controllers.ListHandler[*struct {
			resolvers.AuthResolver
		}](presets.ListMapPresets))

	router.Register(api, "CreateUpdateMapPreset",
		huma.Operation{
			Path:        "/map-presets",
			Method:      http.MethodPut,
			Summary:     "Save map preset",
			Description: "Creates a new map preset or updates an existing one. If the preset already exists and is owned by the current user, it will be updated. Admins can update global presets.",
		}, controllers.CreateHandlerWithInput[*CreateMapPresetInput, presets.MapToolPresetInput])

	router.Register(api, "DeleteMapPreset",
		huma.Operation{
			Path:        "/map-presets/{name}",
			Method:      http.MethodDelete,
			Summary:     "Delete map preset",
			Description: "Deletes a map preset by name. Only the owner of the preset or an admin can delete it.",
		}, controllers.DeleteHandler[*struct {
			controllers.NameInput
			resolvers.AuthRequired
		}](presets.DeleteMapPreset))
}

type CreateMapPresetInput struct {
	controllers.CreateHandlerInput[presets.MapToolPresetInput, presets.MapToolPreset]
}

func (i *CreateMapPresetInput) Resolve(ctx huma.Context) []error {

	maybePreset, err := presets.GetMapPreset(i.DB(), i.Body.Name)

	if user, ok := i.AuthUser(); ok {
		if user.IsGranted(people.Admin) {
			return nil // Contributors can create or update presets without restrictions
		}

		if i.Body.IsPublic && !user.IsGranted(people.Contributor) {
			return []error{
				huma.NewError(http.StatusForbidden, "Only contributors and above are allowed to create public presets"),
			}
		}

		if (maybePreset.IsGlobal || i.Body.IsGlobal) && !user.IsGranted(people.Maintainer) {
			return []error{
				huma.NewError(http.StatusForbidden, "Only maintainers and admins are allowed to create or update global presets"),
			}
		}

		if db.IsNoData(err) {
			return nil // No existing preset, so we can create a new one
		}

		if maybePreset.Meta.CreatedBy.Value.ID != user.ID {
			return []error{
				huma.NewError(http.StatusUnprocessableEntity, "You are not allowed to create or update this preset", &huma.ErrorDetail{
					Message:  "Name is already in use by another user.",
					Location: "body.name",
					Value:    i.Body.Name,
				}),
			}
		}

		if err != nil {
			return []error{err}
		}

	}

	return i.CreateHandlerInput.Resolve(ctx)

}
