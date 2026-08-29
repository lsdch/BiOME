package app

import (
	"context"
	"fmt"
	"net/http"
	"path/filepath"

	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/lsdch/biome/config"
	"github.com/lsdch/biome/controllers"
	"github.com/lsdch/biome/db"
	"github.com/lsdch/biome/imports"
	"github.com/lsdch/biome/middleware"
	"github.com/lsdch/biome/router"
	"github.com/lsdch/biome/services"
	"github.com/lsdch/biome/services/crossref"
	"github.com/lsdch/biome/services/gbif"
	"github.com/lsdch/biome/services/geoapify"
	"github.com/lsdch/biome/services/storage"
	"github.com/lsdch/biome/stores"
	"github.com/sirupsen/logrus"
	log "github.com/sirupsen/logrus"
)

//go:generate go run ../generators/enums/generate_enums.go
//go:generate go run ../generators/mapstructure/generate_mapstructure.go ../models

func NewPgxPool(ctx context.Context, cfg config.DBConfig) (*pgxpool.Pool, error) {

	poolCfg, err := pgxpool.ParseConfig(cfg.DSN)
	if err != nil {
		return nil, fmt.Errorf("parse dsn: %w", err)
	}
	poolCfg.MaxConns = cfg.MaxConns
	poolCfg.MinConns = cfg.MinConns
	poolCfg.MaxConnLifetime = cfg.MaxConnLifetime
	poolCfg.MaxConnIdleTime = cfg.MaxConnIdleTime

	pool, err := pgxpool.NewWithConfig(ctx, poolCfg)
	if err != nil {
		return nil, fmt.Errorf("create pool: %w", err)
	}

	// health check
	ctx, cancel := context.WithTimeout(ctx, cfg.HealthTimeout)
	defer cancel()

	if err := pool.Ping(ctx); err != nil {
		pool.Close()
		return nil, fmt.Errorf("ping db: %w", err)
	}

	return pool, nil
}

type AppServices struct {
	AbioticService     *services.AbioticService
	ImportBatchService *services.ImportBatchService
	AuthService        *services.AuthService
	SettingsService    *services.SettingsService
	SamplingsService   *services.SamplingService
	HabitatService     *services.HabitatService
	PublicationService *services.PublicationService
	DatasetService     *services.DatasetsService
	GeoapifyService    *geoapify.GeoapifyService
	OccurrencesService *services.OccurrencesService
	AccountsService    *services.AccountService
	LocationService    *services.LocationService
	TaxonomyService    *services.TaxonomyService
	TaxonResolver      imports.TaxonResolver
	FileStorage        storage.RawFileStorage
}

type App struct {
	DB             *db.DB
	Config         config.Config
	Router         *router.Router
	Services       *AppServices
	Controllers    []controllers.Controller
	bootstrap      *AppBootstrap
	importsManager *imports.ImportManager
}

func NewApp(config config.Config) *App {
	dbPool, err := NewPgxPool(context.Background(), config.DB)
	if err != nil {
		panic(fmt.Sprintf("failed to create database pool: %v", err))
	}

	database := db.NewDB(dbPool)
	router := makeRouter(config.API)

	appServices := &AppServices{
		AbioticService:     services.NewAbioticService(),
		AuthService:        services.NewAuthService(config.AuthTokens),
		SettingsService:    services.NewSettingsService(config),
		SamplingsService:   services.NewSamplingService(),
		HabitatService:     services.NewHabitatService(),
		PublicationService: services.NewPublicationService(),
		DatasetService:     services.NewDatasetsService(),
		LocationService:    services.NewLocationService(),
		GeoapifyService:    geoapify.NewGeoapifyService(http.DefaultClient, config.Geoapify),
	}

	fileStorage, err := storage.NewFilesystemRawFileStorage(config.RawFileStorageRoot)
	if err != nil {
		panic(fmt.Sprintf("failed to create file storage: %v", err))
	}
	appServices.FileStorage = fileStorage
	appServices.ImportBatchService = services.NewImportBatchService(fileStorage)

	appServices.AccountsService = services.NewAccountService(appServices.AuthService, config.Bootstrap)

	authMiddleware := middleware.NewAuthMiddleware(router.API, database, config.AuthTokens, appServices.AuthService)
	router.API.UseMiddleware(authMiddleware.AuthN, authMiddleware.AuthZ)

	gbifClient := gbif.NewClient(config.GBIF)
	appServices.TaxonomyService = services.NewTaxonomyService(gbifClient)
	appServices.TaxonResolver = imports.NewTaxonResolutionService(gbifClient)

	appServices.OccurrencesService = services.NewOccurrencesService(
		appServices.SamplingsService,
		appServices.DatasetService,
		appServices.ImportBatchService,
		appServices.TaxonomyService,
		stores.NewOccurrenceStore(),
	)

	crossrefClient := crossref.NewClient(config.CrossRef)
	bibliographyResolver := imports.NewBibliographyResolver(crossrefClient)

	return &App{
		DB:        database,
		Config:    config,
		Router:    router,
		Services:  appServices,
		bootstrap: NewAppBootstrap(database, config.Bootstrap, appServices),
		importsManager: imports.NewImportManager(
			database,
			stores.NewBatchesStore(),
			appServices.TaxonResolver,
			bibliographyResolver,
			appServices.SamplingsService,
			appServices.OccurrencesService,
			fileStorage,
		),
	}
}

func (a *App) Bootstrap() {
	logrus.Infof("Bootstrapping database...")
	if err := a.bootstrap.Bootstrap(context.Background()); err != nil {
		panic(fmt.Sprintf("failed to bootstrap database: %v", err))
	}

	if err := a.importsManager.Restore(context.Background()); err != nil {
		panic(fmt.Sprintf("failed to restore imports: %v", err))
	}
}

func (a *App) RegisterRoutes() {
	logrus.Infof("Registering routes...")
	a.Controllers = []controllers.Controller{
		controllers.NewAbioticsController(a.DB, a.Services.AbioticService),
		controllers.NewImportController(a.DB, a.importsManager),
		controllers.NewArticlesController(a.DB, a.Services.PublicationService),
		controllers.NewGeoapifyController(a.DB, a.Services.GeoapifyService),
		controllers.NewHabitatsController(a.DB, a.Services.HabitatService),
		controllers.NewOccurrenceController(a.DB, a.Services.OccurrencesService),
		controllers.NewAccountsController(a.DB, a.Services.AccountsService),
		controllers.NewAuthController(a.DB, a.Services.AuthService),
		controllers.NewSettingsController(a.DB, a.Services.SettingsService),
		controllers.NewLocationController(a.DB, a.Services.LocationService),
		controllers.NewSamplingController(a.DB, a.Services.SamplingsService),
		controllers.NewTaxonomyController(a.DB, a.Services.TaxonomyService),
		controllers.NewDatasetController(a.DB, a.Services.DatasetService),
		controllers.NewImportBatchController(a.DB, a.Services.ImportBatchService),
	}
	for _, controller := range a.Controllers {
		controller.RegisterRoutes(a.Router)
	}
}

func (a *App) WriteOpenAPISpec(outputPath string) error {
	absPath, err := filepath.Abs(outputPath)
	if err != nil {
		return fmt.Errorf("failed to get absolute path for OpenAPI spec: %v", err)
	}
	logrus.Infof("Writing OpenAPI spec to %s", absPath)
	// registry := a.Router.API.OpenAPI().Components.Schemas
	// registry.Map()["ListOccurrencesParams"] = registry.Schema(reflect.TypeFor[stores.ListOccurrencesParams](), false, "ListOccurrencesParams")
	return a.Router.WriteSpecJSON(outputPath)
}

func (a *App) Run() {
	if err := a.Router.Engine.Run(fmt.Sprintf("%s:%s", a.Config.API.Host, a.Config.API.Port)); err != nil {
		log.Fatalf("Failed to start Gin router: %v", err)
	}
}

func (a *App) Close() {
	a.DB.Close()
}

func makeRouter(cfg config.APIConfig) *router.Router {
	r := gin.Default()
	r.Use(gin.Recovery())

	// ginAPI := setupRoutes(r, cfg)
	// ginAPI.Static("/assets/", "./assets")
	apiRouter := router.New(r, cfg.BasePath, cfg.ToHumaConfig())
	apiRouter.BaseAPI.Static("/assets/", "./assets")
	return &apiRouter
}
