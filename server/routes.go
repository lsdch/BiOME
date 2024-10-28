package main

import (
	"darco/proto/controllers/events"
	"darco/proto/controllers/institution"
	"darco/proto/controllers/location"
	"darco/proto/controllers/person"
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
	sites.RegisterDatasetRoutes(router)
	events.RegisterRoutes(router)
}
