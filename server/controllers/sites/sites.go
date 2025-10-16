package sites

import (
	"context"
	"net/http"

	"github.com/geldata/gel-go/geltypes"
	"github.com/lsdch/biome/controllers"
	"github.com/lsdch/biome/controllers/occurrences"
	"github.com/lsdch/biome/models/occurrence"
	"github.com/lsdch/biome/resolvers"
	"github.com/lsdch/biome/router"

	"github.com/danielgtaylor/huma/v2"
)

func SitesAPI(r router.Router) router.Group {
	return r.RouteGroup("/sites").WithTags([]string{"Location"})
}

func RegisterRoutes(r router.Router) {

	sites_API := SitesAPI(r)

	router.Register(sites_API, "ListSites",
		huma.Operation{
			Path:        "/",
			Method:      http.MethodGet,
			Summary:     "List sites",
			Description: "List all registered sites",
		}, controllers.ListHandlerWithOpts[*struct {
			resolvers.AuthResolver
			occurrence.ListSitesOptions
		}](occurrence.ListSites))

	router.Register(sites_API, "GetSite",
		huma.Operation{
			Path:        "/{code}",
			Method:      http.MethodGet,
			Summary:     "Get site",
			Description: "Get site infos using its code",
		}, controllers.GetByCodeHandler(occurrence.GetSite))

	router.Register(sites_API, "CreateSite",
		huma.Operation{
			Path:        "/",
			Method:      http.MethodPost,
			Summary:     "Create site",
			Description: "Create site infos using its code",
		}, controllers.CreateHandler[occurrence.SiteInput])

	router.Register(sites_API, "UpdateSite",
		huma.Operation{
			Path:        "/{code}",
			Method:      http.MethodPatch,
			Summary:     "Update site",
			Description: "Update site infos using its code",
		},
		controllers.UpdateByCodeHandler[occurrence.SiteUpdate],
	)

	router.Register(sites_API, "ListSiteSamplings",
		huma.Operation{
			Path:    "/{code}/samplings",
			Method:  http.MethodGet,
			Summary: "List samplings at site",
		},
		controllers.GetByCodeHandler(occurrence.ListSamplingsAtSite),
	)

	router.Register(sites_API, "CreateSamplingAtSite",
		huma.Operation{
			Path:        "/{code}/samplings",
			Method:      http.MethodPost,
			Summary:     "Create sampling at site",
			Description: "Register sampling event on a site identified by its code",
		},
		controllers.UpdateByCodeHandler[occurrence.SamplingInput],
	)

	router.Register(sites_API, "SiteAddExternalOccurrence",
		huma.Operation{
			Tags:        []string{"Occurrences"},
			Path:        "/{code}/occurrences/external",
			Method:      http.MethodPost,
			Summary:     "Add occurrence at site",
			Description: "Register new occurrence at site, including event + sampling specification and biomaterial identification",
		},
		SiteAddExternalOccurrence,
	)
}

type SiteAddExternalOccurrenceInput struct {
	resolvers.AccessRestricted[resolvers.Contributor]
	controllers.CodeInput
	Body struct {
		Sampling    occurrence.SamplingInput           `json:"sampling"`
		BioMaterial occurrence.ExternalOccurrenceInput `json:"biomaterial"`
	} `nameHint:"ExternalOccurrenceAtSiteInput"`
}

func SiteAddExternalOccurrence(ctx context.Context, input *SiteAddExternalOccurrenceInput) (*occurrences.RegisterOccurrenceOutput, error) {
	siteCode := input.Identifier()
	var created occurrence.GenericOccurrence[occurrence.SamplingOutline]
	err := input.DB().Tx(context.Background(), func(ctx context.Context, tx geltypes.Tx) error {

		sampling, err := input.Body.Sampling.Save(tx, siteCode)
		if err != nil {
			return err
		}
		biomat, err := input.Body.BioMaterial.Save(tx, sampling.Number)
		if err != nil {
			return err
		}
		created = biomat
		return nil
	})
	if err != nil {
		return nil, controllers.StatusError(err)
	}
	return &occurrences.RegisterOccurrenceOutput{Body: created}, nil
}
