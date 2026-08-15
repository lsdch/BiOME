package main

import (
	"github.com/danielgtaylor/huma/v2"
	"github.com/gin-gonic/gin"
	"github.com/lsdch/biome/app"
	"github.com/lsdch/biome/config"
	"github.com/sirupsen/logrus"
)

//go:generate go run generators/enums/generate_enums.go
//go:generate go run generators/mapstructure/generate_mapstructure.go models

func main() {
	huma.DefaultArrayNullable = false

	cfg, err := config.LoadConfig(".", "config")
	if err != nil {
		logrus.Fatalf("Failed to load config file: %v", err)
	}

	logrus.Infof("Loaded backend configuration")

	if gin.Mode() == gin.DebugMode {
		logrus.SetLevel(logrus.DebugLevel)
	}
	// Disable logging all routes
	gin.DebugPrintRouteFunc = func(httpMethod, absolutePath, handlerName string, nuHandlers int) {}

	biome := app.NewApp(cfg)
	biome.Bootstrap()
	biome.RegisterRoutes()
	if err := biome.WriteOpenAPISpec("../client/openapi.json"); err != nil {
		logrus.Fatalf("Failed to write OpenAPI spec: %v", err)
	}

	defer biome.Close()
	biome.Run()
}
