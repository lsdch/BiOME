package main

import (
	"darco/proto/db"
	mw "darco/proto/middlewares"
	"darco/proto/models/location"
	"darco/proto/models/settings"
	"darco/proto/router"
	"darco/proto/services/email"
	"fmt"

	"github.com/danielgtaylor/huma/v2"
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
	log "github.com/sirupsen/logrus"
)

//go:generate go run generators/generate_enums.go

func apiConfig(basePath string) huma.Config {
	instance := settings.Instance()
	title := fmt.Sprintf("%s API", instance.Name)
	description, ok := instance.Description.Get()
	if !ok {
		description = instance.Name
	}

	var cfg = huma.DefaultConfig(title, "1.0")
	cfg.Info.Description = fmt.Sprintf("%s API", description)
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
	r.Use(mw.AuthenticationMiddleware)

	ginAPI := setupRoutes(r, "/api/v1")
	ginAPI.Static("/assets/", "./assets")
	return r
}

func main() {

	gin.ForceConsoleColor()
	if gin.Mode() == gin.DebugMode {
		log.SetLevel(log.DebugLevel)
	}

	if err := email.LoadTemplates("templates/**"); err != nil {
		log.Fatalf("Failed to load email templates: %v", err)
	}

	if err := location.SetupCountries(db.Client()); err != nil {
		logrus.Fatalf("Failed to setup countries in database: %v", err)
	}

	r := setupRouter()
	if err := r.Run(":8080"); err != nil {
		log.Fatalf("Failed to start Gin router: %v", err)
	}
	defer db.Client().Close()
}
