package app

import (
	"github.com/lsdch/biome/imports"
	"github.com/lsdch/biome/services"
	"github.com/lsdch/biome/services/geoapify"
)

type AppServices struct {
	AbioticService     *services.AbioticService
	ImportBatchService *services.ImportBatchService
	AuthService        *services.AuthService
	SettingsService    *services.SettingsService
	SamplingsService   *services.SamplingService
	HabitatService     *services.HabitatService
	ArticleService     *services.PublicationService
	DatasetService     *services.DatasetsService
	GeoapifyService    *geoapify.GeoapifyService
	OccurrencesService *services.OccurrencesService
	AccountsService    *services.AccountService
	LocationService    *services.LocationService
	TaxonomyService    *services.TaxonomyService
	TaxonResolver      imports.TaxonResolver
}

// type ServiceFactory struct {
// 	db     *db.DB
// 	gbif   *gbif.GBIFClient
// 	config config.Config
// }

// func NewServiceFactory(db *db.DB, cfg config.Config) *ServiceFactory {
// 	return &ServiceFactory{
// 		db:     db,
// 		gbif:   gbif.NewClient(cfg),
// 		config: cfg,
// 	}
// }

// func (f *ServiceFactory) Services(q db.Querier) *AppServices {

// 	samplings := services.NewSamplingService(q)

// 	importService := imports.NewImportService(
// 		q,
// 		f.gbif,
// 		samplings,
// 		imports.NewTaxonResolutionService(q, f.gbif),
// 	)

// 	return &AppServices{
// 		ImportService:      importService,
// 		ImportBatchService: services.NewImportBatchService(q),
// 		AuthService:        services.NewAuthService(q, f.config.AuthTokens),
// 		SettingsService:    services.NewSettingsService(q, f.config),
// 		SamplingsService:   samplings,
// 		HabitatService:     services.NewHabitatService(q),
// 		ArticleService:     services.NewArticleService(q),
// 		DatasetService:     services.NewDatasetsService(q),
// 		TaxonomyService:    services.NewTaxonomyService(q),
// 		LocationService:    services.NewLocationService(q),
// 		GeoapifyService:    geoapify.NewGeoapifyService(q, http.DefaultClient, f.config.Geoapify),
// 	}
// }
