package main

import (
	"darco/proto/controllers/institution"
	country "darco/proto/controllers/location"
	"darco/proto/controllers/person"
	"darco/proto/controllers/taxonomy"
	accounts "darco/proto/controllers/users"
	"darco/proto/db"
	mw "darco/proto/middlewares"
	"darco/proto/models/validations"
	"darco/proto/router"
	"darco/proto/services/email"
	"reflect"
	"strings"

	"github.com/danielgtaylor/huma/v2"
	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	"github.com/go-playground/validator/v10"
	log "github.com/sirupsen/logrus"
)

//go:generate go run generators/generate_enums.go

func humaConfig(basePath string) huma.Config {
	var cfg = huma.DefaultConfig("Proto API", "1.0")
	cfg.Info.Description = "DarCo Prototype API"
	cfg.Components.SecuritySchemes = map[string]*huma.SecurityScheme{
		"bearer": {
			Type:         "http",
			Scheme:       "bearer",
			BearerFormat: "JWT",
		},
	}
	cfg.Info.Contact = &huma.Contact{
		Name:  "Louis Duchemin",
		Email: "louis.duchemin@univ-lyon1.fr",
	}
	cfg.OpenAPI.Servers = []*huma.Server{
		{URL: basePath},
	}

	return cfg
}

func setupRoutes(r *gin.Engine, basePath string) *gin.RouterGroup {
	router := router.New(r, basePath, humaConfig(basePath))
	accounts.RegisterRoutes(router)
	institution.RegisterRoutes(router)
	person.RegisterRoutes(router)
	country.RegisterRoutes(router)
	taxonomy.RegisterRoutes(router)
	taxonomy.RegisterImportRoutes(router)
	if err := router.WriteSpecJSON("../client/openapi.json"); err != nil {
		panic(err)
	}

	return router.BaseAPI
}

func setupRouter() *gin.Engine {
	r := gin.Default()
	r.Use(gin.Recovery())
	r.Use(mw.ErrorHandler)
	r.Use(mw.AuthenticationMiddleware)

	setupRoutes(r, "/api/v1")

	return r
}

func setupValidators() {
	if engine, ok := binding.Validator.Engine().(*validator.Validate); ok {
		validations.RegisterValidators(engine)
		// Use json names
		engine.RegisterTagNameFunc(func(fld reflect.StructField) string {
			name := strings.SplitN(fld.Tag.Get("json"), ",", 2)[0]
			// skip if tag key says it should be ignored
			if name == "-" {
				return ""
			}
			return name
		})
	}
}

func main() {

	gin.ForceConsoleColor()
	if gin.Mode() == gin.DebugMode {
		log.SetLevel(log.DebugLevel)
	}

	setupValidators()

	if err := email.LoadTemplates("templates/**"); err != nil {
		log.Fatalf("Failed to load email templates: %v", err)
	}
	r := setupRouter()
	if err := r.Run(":8080"); err != nil {
		log.Fatalf("Failed to start Gin router: %v", err)
	}
	defer db.Client().Close()
}
