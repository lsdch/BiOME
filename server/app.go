package main

import (
	"context"
	"fmt"

	"github.com/danielgtaylor/huma/v2"
	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/lsdch/biome/config"
	"github.com/lsdch/biome/controllers"
	"github.com/lsdch/biome/db"
	"github.com/lsdch/biome/router"
	"github.com/lsdch/biome/services"
	"github.com/lsdch/biome/services/gbif"
	"github.com/lsdch/biome/services/imports"
	"github.com/sirupsen/logrus"
	log "github.com/sirupsen/logrus"
)

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
	DB       *db.DB
	Config   config.Config
	Router   *router.Router
	Services *AppServices
}

type AppServices struct {
	ImportService    *imports.ImportService
	AuthService      *services.AuthService
	SettingsService  *services.SettingsService
	SamplingsService *services.SamplingService
	BootstrapService *services.BootstrapService
	HabitatService   *services.HabitatService
}

func NewApp(config config.Config) *App {
	dbPool, err := NewPgxPool(context.Background(), config.DB)
	if err != nil {
		panic(fmt.Sprintf("failed to create database pool: %v", err))
	}

	database := db.NewDB(dbPool)
	router := makeRouter(config.API)

	appServices := &AppServices{
		ImportService:    imports.NewImportService(database, gbif.NewClient(config)),
		AuthService:      services.NewAuthService(database, config.AuthTokens),
		SettingsService:  services.NewSettingsService(database, config),
		SamplingsService: services.NewSamplingService(database),
		HabitatService:   services.NewHabitatService(database),
	}

	appServices.BootstrapService = services.NewBootstrapService(
		database,
		config.Bootstrap,
		appServices.SettingsService,
		appServices.SamplingsService,
		appServices.HabitatService,
	)

	return &App{
		DB:       database,
		Config:   config,
		Router:   router,
		Services: appServices,
	}
}

func (a *App) Bootstrap() {
	if err := a.Services.BootstrapService.Bootstrap(context.Background()); err != nil {
		panic(fmt.Sprintf("failed to bootstrap database: %v", err))
	}
}

func (a *App) RegisterRoutes() {
	r := a.Router
	importController := controllers.NewImportController(a.Services.ImportService)
	importController.RegisterRoutes(r)
}

func (a *App) Run() {
	if err := a.Router.Engine.Run(fmt.Sprintf("%s:%s", a.Config.API.Host, a.Config.API.Port)); err != nil {
		log.Fatalf("Failed to start Gin router: %v", err)
	}
}

func (a *App) Close() {
	a.DB.Close()
}

func apiConfig(c config.APIConfig) huma.Config {
	var cfg = huma.DefaultConfig(c.Title, c.Version)
	cfg.Info.Description = c.Description
	cfg.Components.SecuritySchemes = map[string]*huma.SecurityScheme{
		"bearer": {
			Type:         "http",
			Scheme:       "bearer",
			BearerFormat: "JWT",
		},
		"cookieAuth": {
			Type: "apiKey",
			In:   "cookie",
			Name: "auth_token",
		},
	}
	cfg.Info.Contact = &huma.Contact{
		Name:  c.ContactName,
		Email: c.ContactEmail,
	}
	cfg.OpenAPI.Servers = []*huma.Server{
		{URL: c.BasePath},
	}
	cfg.Security = []map[string][]string{
		{"bearer": {}},
		{"cookieAuth": {}},
	}

	return cfg
}

func setupRoutes(r *gin.Engine, cfg config.APIConfig) *gin.RouterGroup {
	apiConfig := apiConfig(cfg)
	router := router.New(r, cfg.BasePath, apiConfig)

	router.CollectRoutes()

	if err := router.WriteSpecJSON("../client/openapi.json"); err != nil {
		panic(err)
	}

	return router.BaseAPI
}

func makeRouter(cfg config.APIConfig) *router.Router {
	r := gin.Default()
	r.Use(gin.Recovery())

	ginAPI := setupRoutes(r, cfg)
	ginAPI.Static("/assets/", "./assets")
	apiRouter := router.New(r, cfg.BasePath, cfg.ToHumaConfig())
	return &apiRouter
}

func main() {
	huma.DefaultArrayNullable = false

	cfg, err := config.LoadConfig(".", "config")
	if err != nil {
		log.Fatalf("Failed to load config file: %v", err)
	}

	logrus.Infof("Loaded backend configuration: %+v", cfg)

	if gin.Mode() == gin.DebugMode {
		log.SetLevel(log.DebugLevel)
	}
	// Disable logging all routes
	gin.DebugPrintRouteFunc = func(httpMethod, absolutePath, handlerName string, nuHandlers int) {}

	app := NewApp(cfg)
	app.RegisterRoutes()

	defer app.Close()
	app.Run()
}
