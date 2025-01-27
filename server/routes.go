package main

import (
	"github.com/lsdch/biome/controllers/biomaterial"
	"github.com/lsdch/biome/controllers/datasets"
	"github.com/lsdch/biome/controllers/events"
	"github.com/lsdch/biome/controllers/habitats"
	"github.com/lsdch/biome/controllers/institution"
	"github.com/lsdch/biome/controllers/location"
	"github.com/lsdch/biome/controllers/occurrences"
	"github.com/lsdch/biome/controllers/person"
	"github.com/lsdch/biome/controllers/references"
	"github.com/lsdch/biome/controllers/sequences"
	"github.com/lsdch/biome/controllers/services"
	"github.com/lsdch/biome/controllers/settings"
	"github.com/lsdch/biome/controllers/sites"
	"github.com/lsdch/biome/controllers/taxonomy"
	users "github.com/lsdch/biome/controllers/users"
	"github.com/lsdch/biome/router"
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
	biomaterial.RegisterRoutes(router)
	references.RegisterRoutes(router)
	occurrences.RegisterRoutes(router)
	services.RegisterGeoapifyRoutes(router)
}
