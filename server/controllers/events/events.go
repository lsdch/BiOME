package events

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

func registerEventsRoutes(r router.Router) {
	eventsAPI := r.RouteGroup("/events").
		WithTags([]string{"Events"})

	router.Register(eventsAPI, "UpdateEvent",
		huma.Operation{
			Path:    "/{id}",
			Method:  http.MethodPatch,
			Summary: "Update event",
		},
		controllers.UpdateByIDHandler[occurrence.EventUpdate])

	router.Register(eventsAPI, "DeleteEvent",
		huma.Operation{
			Path:    "/{id}",
			Method:  http.MethodDelete,
			Summary: "Delete event",
		},
		controllers.DeleteByIDHandler(occurrence.DeleteEvent))

	router.Register(eventsAPI, "CreateSamplingAtEvent",
		huma.Operation{
			Path:    "/{id}/samplings",
			Method:  http.MethodPost,
			Summary: "Create sampling at event",
		},
		controllers.UpdateByIDHandler[occurrence.SamplingInput])

	router.Register(eventsAPI, "UpdateSpotting",
		huma.Operation{
			Path:    "/{id}/spottings",
			Method:  http.MethodPut,
			Summary: "Update spotting",
		},
		controllers.UpdateByIDHandler[occurrence.SpottingUpdate])

	router.Register(eventsAPI, "EventAddExternalOccurrence",
		huma.Operation{
			Tags:        []string{"Occurrences"},
			Path:        "/{id}/occurrences/external",
			Method:      http.MethodPost,
			Summary:     "Add occurrence from event",
			Description: "Register new occurrence resulting from the event, including sampling specification and biomaterial identification",
		},
		EventAddExternalOccurrence)

}

type EventAddExternalOccurrenceInput struct {
	resolvers.AccessRestricted[resolvers.Contributor]
	ID   geltypes.UUID `path:"id" format:"uuid" doc:"Event ID"`
	Body struct {
		Sampling    occurrence.SamplingInput       `json:"sampling" doc:"New sampling action during referenced event"`
		BioMaterial occurrence.ExternalBioMatInput `json:"biomaterial" doc:"New occurrence resulting from the sampling action"`
	} `nameHint:"ExternalOccurrenceAtEventInput"`
}

func EventAddExternalOccurrence(ctx context.Context, input *EventAddExternalOccurrenceInput) (*occurrences.RegisterOccurrenceOutput, error) {
	var occurrence occurrence.BioMaterialWithDetails
	err := input.DB().Tx(context.Background(), func(ctx context.Context, tx geltypes.Tx) error {
		event := input.ID
		sampling, err := input.Body.Sampling.Save(tx, event)
		if err != nil {
			return err
		}
		created, err := input.Body.BioMaterial.Save(tx, sampling.ID)
		if err != nil {
			return err
		}
		occurrence = created
		return nil
	})
	if err != nil {
		return nil, controllers.StatusError(err)
	}
	return &occurrences.RegisterOccurrenceOutput{Body: occurrence}, nil
}
