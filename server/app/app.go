package main

import (
	"context"
	"fmt"
	"net/http"

	"github.com/danielgtaylor/huma/v2"
	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/lsdch/biome/config"
	"github.com/lsdch/biome/controllers"
	"github.com/lsdch/biome/db"
	"github.com/lsdch/biome/middleware"
	"github.com/lsdch/biome/router"
	"github.com/lsdch/biome/services"
	"github.com/lsdch/biome/services/gbif"
	"github.com/lsdch/biome/services/geoapify"
	"github.com/lsdch/biome/services/imports"
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

type App struct {
	DB          *db.DB
	Config      config.Config
	Router      *router.Router
	Services    *AppServices
	Controllers []controllers.Controller
	bootstrap   *AppBootstrap
}

func NewApp(config config.Config) *App {
	dbPool, err := NewPgxPool(context.Background(), config.DB)
	if err != nil {
		panic(fmt.Sprintf("failed to create database pool: %v", err))
	}

	database := db.NewDB(dbPool)
	router := makeRouter(config.API)

	gbifClient := gbif.NewClient(config)

	appServices := &AppServices{
		ImportBatchService: services.NewImportBatchService(),
		AuthService:        services.NewAuthService(config.AuthTokens),
		SettingsService:    services.NewSettingsService(config),
		SamplingsService:   services.NewSamplingService(),
		HabitatService:     services.NewHabitatService(),
		ArticleService:     services.NewArticleService(),
		DatasetService:     services.NewDatasetsService(),
		TaxonomyService:    services.NewTaxonomyService(),
		LocationService:    services.NewLocationService(),
		GeoapifyService:    geoapify.NewGeoapifyService(http.DefaultClient, config.Geoapify),
	}

	appServices.ImportService = imports.NewImportService(gbifClient, appServices.SamplingsService, imports.NewTaxonResolutionService(gbifClient))
	appServices.AccountsService = services.NewAccountService(appServices.AuthService, config.Bootstrap)

	appServices.OccurrencesService = services.NewOccurrencesService(
		appServices.SamplingsService,
		appServices.DatasetService,
		appServices.ImportBatchService,
		appServices.TaxonomyService,
	)

	authMiddleware := middleware.NewAuthMiddleware(router.API, database, config.AuthTokens, appServices.AuthService)
	router.API.UseMiddleware(authMiddleware.AuthN, authMiddleware.AuthZ)

	return &App{
		DB:        database,
		Config:    config,
		Router:    router,
		Services:  appServices,
		bootstrap: NewAppBootstrap(database, config.Bootstrap, appServices),
	}
}

func (a *App) Bootstrap() {
	logrus.Infof("Bootstrapping database...")
	if err := a.bootstrap.Bootstrap(context.Background()); err != nil {
		panic(fmt.Sprintf("failed to bootstrap database: %v", err))
	}
}

func (a *App) RegisterRoutes() {

	a.Controllers = []controllers.Controller{
		controllers.NewImportController(a.DB, a.Services.ImportService),
		controllers.NewArticlesController(a.DB, a.Services.ArticleService),
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
	}
	for _, controller := range a.Controllers {
		controller.RegisterRoutes(a.Router)
	}
	if err := a.Router.WriteSpecJSON("../client/openapi.json"); err != nil {
		panic(err)
	}
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

func main() {
	huma.DefaultArrayNullable = false

	cfg, err := config.LoadConfig(".", "config")
	if err != nil {
		log.Fatalf("Failed to load config file: %v", err)
	}

	logrus.Infof("Loaded backend configuration")

	if gin.Mode() == gin.DebugMode {
		log.SetLevel(log.DebugLevel)
	}
	// Disable logging all routes
	gin.DebugPrintRouteFunc = func(httpMethod, absolutePath, handlerName string, nuHandlers int) {}

	app := NewApp(cfg)
	app.Bootstrap()
	app.RegisterRoutes()

	defer app.Close()
	app.Run()
}
