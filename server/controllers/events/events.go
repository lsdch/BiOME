package events

import (
	"darco/proto/controllers"
	"darco/proto/models/occurrence"
	"darco/proto/router"
	"net/http"

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

	router.Register(eventsAPI, "UpdateSpotting",
		huma.Operation{
			Path:    "/{id}/spottings",
			Method:  http.MethodPut,
			Summary: "Update spotting",
		},
		controllers.UpdateByIDHandler[occurrence.SpottingUpdate])
}
