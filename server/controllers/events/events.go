package events

import (
	"net/http"

	"github.com/lsdch/biome/controllers"
	"github.com/lsdch/biome/models/occurrence"
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

	router.Register(eventsAPI, "UpdateSpotting",
		huma.Operation{
			Path:    "/{id}/spottings",
			Method:  http.MethodPut,
			Summary: "Update spotting",
		},
		controllers.UpdateByIDHandler[occurrence.SpottingUpdate])

}
