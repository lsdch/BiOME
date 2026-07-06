package main

// func setupRoutes(r *gin.Engine, basePath string) *gin.RouterGroup {
// 	apiConfig := apiConfig(basePath)
// 	router := router.New(r, basePath, apiConfig)

// 	router.CollectRoutes()

// 	if err := router.WriteSpecJSON("../client/openapi.json"); err != nil {
// 		panic(err)
// 	}

// 	return router.BaseAPI
// }

// func setupRouter() *gin.Engine {
// 	r := gin.Default()
// 	r.Use(gin.Recovery())

// 	ginAPI := setupRoutes(r, "/api/v1")
// 	ginAPI.Static("/assets/", "./assets")
// 	return r
// }

// func main() {

// 	huma.DefaultArrayNullable = false

// 	if config, err := config.LoadConfig(".", "config"); err != nil {
// 		log.Fatalf("Failed to load config file: %v", err)
// 	} else {
// 		logrus.Infof("Loaded backend configuration: %+v", config)
// 	}

// 	if gin.Mode() == gin.DebugMode {
// 		log.SetLevel(log.DebugLevel)
// 	}
// 	// Disable logging all routes
// 	gin.DebugPrintRouteFunc = func(httpMethod, absolutePath, handlerName string, nuHandlers int) {}

// 	// if err := location.SetupCountries(db.Client()); err != nil {
// 	// 	logrus.Fatalf("Failed to setup countries in database: %v", err)
// 	// }

// 	r := setupRouter()
// 	if err := r.Run(":8080"); err != nil {
// 		log.Fatalf("Failed to start Gin router: %v", err)
// 	}
// 	defer db.Client().Close()
// }
