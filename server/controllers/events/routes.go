package events

import "darco/proto/router"

func RegisterRoutes(r router.Router) {

	registerProgramRoutes(r)
	registerSamplingRoutes(r)
	// eventsAPI := r.RouteGroup("/events").
	// 	WithTags([]string{"Account"})
}
