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
	"net/http"
	"reflect"
	"strings"

	_ "darco/proto/docs" // import swagger docs

	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	"github.com/go-playground/validator/v10"
	log "github.com/sirupsen/logrus"
	swaggerfiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

// Generate OpenAPI docs from swaggo doc comments
//go:generate swag init --parseDependency --parseInternal -g main.go

// @title Proto API
// @version 1.0
// @description Testing Swagger APIs.
// @BasePath /api/v1
// @termsOfService http://swagger.io/terms/
// @contact.name Louis Duchemin
// @contact.email louis.duchemin@univ-lyon1.fr
// @contact.url http://www.swagger.io/support
// @securityDefinitions.apiKey JWT
// @in header
// @name token
// @accept json
// @produce json
// @license.name Apache 2.0
// @license.url http://www.apache.org/licenses/LICENSE-2.0.html
// @host localhost:8080
// @schemes http
func setupRouter() *gin.Engine {
	r := gin.Default()
	r.Use(gin.Recovery())
	r.Use(mw.ErrorHandler)
	r.Use(mw.AuthenticationMiddleware)

	// Swagger docs
	api := r.Group("/api/v1")

	api.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerfiles.Handler))

	country_api := router.Wrap(api.Group("/countries"), router.WithDB)
	country_api.GET("", country.List)
	country_api.GET("/setup", country.Setup)

	taxonomy_api := router.Wrap(api.Group("/taxonomy"), router.WithDB)
	taxonomy_api.GET("", taxonomy.ListTaxa)
	taxonomy_api.POST("", taxonomy.CreateTaxon)
	taxonomy_api.GET("/:code", taxonomy.GetTaxon)
	taxonomy_api.DELETE("/:code", taxonomy.DeleteTaxon)
	taxonomy_api.PATCH("/:code", taxonomy.UpdateTaxon)

	importGBIF := taxonomy.ImportCladeGBIF()
	taxonomy_api.PUT("/import", importGBIF.Endpoint)
	taxonomy_api.RouterGroup.GET("/import", importGBIF.ProgressTracker)
	taxonomy_api.GET("/anchors", taxonomy.ListAnchors)

	api.POST("/login", router.WithDB(accounts.Login))
	api.POST("/logout", accounts.Logout)

	account_api := api.Group("/account")
	account_api.GET("", router.WithUser(accounts.Current))
	account_api.POST("/password", router.WithUser(accounts.SetPassword))

	users_api := router.Wrap(api.Group("/users"), router.WithDB)
	users_api.POST("/register", accounts.Register)
	users_api.GET("/confirm", accounts.ConfirmEmail)
	users_api.POST("/confirm/resend", accounts.ResendConfirmation)
	users_api.POST("/forgotten-password", accounts.RequestPasswordReset)
	users_api.GET("/password-reset/:token", accounts.ValidatePasswordToken)

	people_api := api.Group("/people")

	institution_api := router.Wrap(people_api.Group("/institutions"), router.WithDB)
	institution_api.GET("", institution.List)
	institution_api.POST("", institution.Create)
	institution_api.DELETE("/:code", institution.Delete)
	institution_api.PATCH("/:code", institution.Update)

	person_api := router.Wrap(people_api.Group("/persons"), router.WithDB)
	person_api.GET("", person.List)
	person_api.POST("", person.Create)
	person_api.PATCH("/:id", person.Update)
	person_api.DELETE("/:id", person.Delete)

	// Authorized group (uses gin.BasicAuth() middleware)
	// Same than:
	// authorized := r.Group("/")
	// authorized.Use(gin.BasicAuth(gin.Credentials{
	//	  "foo":  "bar",
	//	  "manu": "123",
	//}))
	authorized := r.Group("/", gin.BasicAuth(gin.Accounts{
		"foo":  "bar", // user:foo password:bar
		"manu": "123", // user:manu password:123
	}))

	authorized.GET("test", func(c *gin.Context) {
		c.String(http.StatusOK, "Test OK")
	})

	return r
}

// Loads config from .env file
// func loadConfig() *config.Config {
// 	config, err := config.LoadConfig(".")
// 	if err != nil {
// 		log.Fatalf("Failed to load configuration file : %v", err)
// 	} else {
// 		log.Infof("Config loaded : %+v", config)
// 	}
// 	return config
// }

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
	// config := loadConfig()

	if err := email.LoadTemplates("templates/**"); err != nil {
		log.Fatalf("Failed to load email templates: %v", err)
	}
	r := setupRouter()
	r.Run(":8080")
	defer db.Client().Close()
}
