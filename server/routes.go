package main

import (
	"darco/proto/controllers/datasets"
	"darco/proto/controllers/events"
	"darco/proto/controllers/habitats"
	"darco/proto/controllers/institution"
	"darco/proto/controllers/location"
	"darco/proto/controllers/person"
	"darco/proto/controllers/sequences"
	"darco/proto/controllers/settings"
	"darco/proto/controllers/sites"
	"darco/proto/controllers/taxonomy"
	users "darco/proto/controllers/users"
	"darco/proto/router"
)

func registerRoutes(router router.Router) {
	users.RegisterRoutes(router)
	institution.RegisterRoutes(router)
	person.RegisterRoutes(router)
	location.RegisterRoutes(router)
	taxonomy.RegisterRoutes(router)
	taxonomy.RegisterImportRoutes(router)
	settings.RegisterRoutes(router)
	sites.RegisterRoutes(router)
	datasets.RegisterRoutes(router)
	events.RegisterRoutes(router)
	sequences.RegisterRoutes(router)
	habitats.RegisterRoutes(router)
}
