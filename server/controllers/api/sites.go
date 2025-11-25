package api

import (
	"context"
	"net/http"

	"github.com/geldata/gel-go/geltypes"
	"github.com/lsdch/biome/controllers"
	"github.com/lsdch/biome/models/occurrence"
	"github.com/lsdch/biome/resolvers"
	"github.com/lsdch/biome/router"

	"github.com/danielgtaylor/huma/v2"
)

func init() {
	sitesAPI := func(r *router.Router) router.Group {
		return r.RouteGroup("/sites").
			WithTags([]string{"Location"})
	}

	router.RegisterSpec(
		sitesAPI,
		"ListSites",
		huma.Operation{
			Path:        "/",
			Method:      http.MethodGet,
			Summary:     "List sites",
			Description: "List all registered sites",
		},
		controllers.ListHandlerWithOpts[*struct {
			resolvers.AuthResolver
			occurrence.ListSitesOptions
		}](occurrence.ListSites),
	)

	router.RegisterSpec(
		sitesAPI,
		"GetSite",
		huma.Operation{
			Path:        "/{code}",
			Method:      http.MethodGet,
			Summary:     "Get site",
			Description: "Get site infos using its code",
		},
		controllers.GetByCodeHandler(occurrence.GetSite))

	router.RegisterSpec(
		sitesAPI,
		"CreateSite",
		huma.Operation{
			Path:        "/",
			Method:      http.MethodPost,
			Summary:     "Create site",
			Description: "Create site infos using its code",
		},
		controllers.CreateHandler[occurrence.SiteInput])

	router.RegisterSpec(
		sitesAPI,
		"UpdateSite",
		huma.Operation{
			Path:        "/{code}",
			Method:      http.MethodPatch,
			Summary:     "Update site",
			Description: "Update site infos using its code",
		},
		controllers.UpdateByCodeHandler[occurrence.SiteUpdate])

	router.RegisterSpec(
		sitesAPI,
		"ListSiteSamplings",
		huma.Operation{
			Path:    "/{code}/samplings",
			Method:  http.MethodGet,
			Summary: "List samplings at site",
		},
		controllers.GetByCodeHandler(occurrence.ListSamplingsAtSite))

	router.RegisterSpec(
		sitesAPI,
		"CreateSamplingAtSite",
		huma.Operation{
			Path:        "/{code}/samplings",
			Method:      http.MethodPost,
			Summary:     "Create sampling at site",
			Description: "Register sampling event on a site identified by its code",
		},
		controllers.UpdateByCodeHandler[*occurrence.SamplingInput])

	router.RegisterSpec(
		sitesAPI,
		"SiteAddOccurrence",
		huma.Operation{
			Tags:        []string{"Occurrences"},
			Path:        "/{code}/occurrences/",
			Method:      http.MethodPost,
			Summary:     "Add occurrence at site",
			Description: "Register new occurrence at site, including event + sampling specification and biomaterial identification",
		},
		SiteAddOccurrence)
}

type SiteAddOccurrenceInput struct {
	resolvers.AccessRestricted[resolvers.Contributor]
	controllers.CodeInput
	Body struct {
		Sampling    occurrence.SamplingInput   `json:"sampling"`
		BioMaterial occurrence.OccurrenceInput `json:"biomaterial"`
	} `nameHint:"OccurrenceAtSiteInput"`
}

func SiteAddOccurrence(ctx context.Context, input *SiteAddOccurrenceInput) (*RegisterOccurrenceOutput, error) {
	siteCode := input.Identifier()
	var created occurrence.BaseOccurrence[occurrence.SamplingOutline]
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
	return &RegisterOccurrenceOutput{Body: created}, nil
}
