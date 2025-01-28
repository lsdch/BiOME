package main

import (
	"fmt"

	"github.com/lsdch/biome/config"
	"github.com/lsdch/biome/db"
	"github.com/lsdch/biome/models/location"
	"github.com/lsdch/biome/router"

	"github.com/danielgtaylor/huma/v2"
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
	log "github.com/sirupsen/logrus"
)

//go:generate go run generators/generate_enums.go

func apiConfig(basePath string) huma.Config {
	title := "BiOME API"
	description := "Biodiversity and Occurrences for Molecular Ecology"

	var cfg = huma.DefaultConfig(title, "0.0")
	cfg.Info.Description = fmt.Sprintf("%s: %s", title, description)
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
		Name:  "Louis Duchemin",
		Email: "louis.duchemin@univ-lyon1.fr",
	}
	cfg.OpenAPI.Servers = []*huma.Server{
		{URL: basePath},
	}
	cfg.Security = []map[string][]string{
		{"bearer": {}},
		{"cookieAuth": {}},
	}

	return cfg
}

func setupRoutes(r *gin.Engine, basePath string) *gin.RouterGroup {
	apiConfig := apiConfig(basePath)
	router := router.New(r, basePath, apiConfig)
	registerRoutes(router)
	if err := router.WriteSpecJSON("../client/openapi.json"); err != nil {
		panic(err)
	}

	return router.BaseAPI
}

func setupRouter() *gin.Engine {
	r := gin.Default()
	r.Use(gin.Recovery())

	ginAPI := setupRoutes(r, "/api/v1")
	ginAPI.Static("/assets/", "./assets")
	return r
}

func main() {

	huma.DefaultArrayNullable = false

	if config, err := config.LoadConfig(".", "config"); err != nil {
		log.Fatalf("Failed to load config file: %v", err)
	} else {
		logrus.Infof("Loaded backend configuration: %+v", config)
	}

	if gin.Mode() == gin.DebugMode {
		log.SetLevel(log.DebugLevel)
	}
	// Disable logging all routes
	gin.DebugPrintRouteFunc = func(httpMethod, absolutePath, handlerName string, nuHandlers int) {}

	if err := location.SetupCountries(db.Client()); err != nil {
		logrus.Fatalf("Failed to setup countries in database: %v", err)
	}

	r := setupRouter()
	if err := r.Run(":8080"); err != nil {
		log.Fatalf("Failed to start Gin router: %v", err)
	}
	defer db.Client().Close()
}
