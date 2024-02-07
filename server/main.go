package main

import (
	"darco/proto/config"
	"darco/proto/controllers/institution"
	country "darco/proto/controllers/location"
	"darco/proto/controllers/person"
	"darco/proto/controllers/taxonomy"
	accounts "darco/proto/controllers/users"
	"darco/proto/middlewares"
	"darco/proto/models"
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
	r.Use(middlewares.ErrorHandler)
	r.Use(middlewares.AuthenticationMiddleware)

	// Swagger docs
	api := r.Group(router.Config.BasePath)
	api.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerfiles.Handler))

	country_api := api.Group("/countries")
	country_api.GET("/", country.List)
	country_api.GET("/setup", country.Setup)

	taxonomy_api := api.Group("/taxonomy")
	taxonomy_api.GET("/", taxonomy.ListTaxa)
	taxonomy_api.GET("/:code", models.WithDB(taxonomy.GetTaxon))
	taxonomy_api.DELETE("/:code", models.WithDB(taxonomy.DeleteTaxon))
	taxonomy_api.PATCH("/:code", models.WithDB(taxonomy.UpdateTaxon))
	importGBIF := taxonomy.ImportCladeGBIF()
	taxonomy_api.PUT("/import", importGBIF.Endpoint)
	taxonomy_api.GET("/import", importGBIF.ProgressTracker)
	taxonomy_api.GET("/anchors", taxonomy.ListAnchors)

	api.POST("/login", models.WithDB(accounts.Login))
	api.POST("/logout", accounts.Logout)

	api.GET("/account", models.WithDB(accounts.Current))
	users_api := api.Group("/users")
	users_api.POST("/register", accounts.Register)
	users_api.GET("/confirm", models.WithDB(accounts.ConfirmEmail))
	users_api.POST("/confirm/resend", models.WithDB(accounts.ResendConfirmation))
	users_api.POST("/forgotten-password", models.WithDB(accounts.RequestPasswordReset))
	users_api.GET("/password-reset/:token", accounts.ValidatePasswordToken)

	people_api := api.Group("/people")
	people_api.GET("/", models.WithDB(person.List))

	institution_api := people_api.Group("/institutions")
	institution_api.GET("/", models.WithDB(institution.List))
	institution_api.POST("/", models.WithDB(institution.Create))
	institution_api.DELETE("/:code", models.WithDB(institution.Delete))
	institution_api.PATCH("/:code", models.WithDB(institution.Update))

	person_api := people_api.Group("/persons")
	person_api.DELETE("/:id", models.WithDB(person.Delete))

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
func loadConfig() {
	config, err := config.LoadConfig(".")
	if err != nil {
		log.Fatalf("Failed to load configuration file : %v", err)
	} else {
		log.Debugf("Config loaded : %+v", config)
	}
}

func setupValidators() {
	if engine, ok := binding.Validator.Engine().(*validator.Validate); ok {
		validations.RegisterValidators(engine)
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
	loadConfig()

	if err := email.LoadTemplates("templates/**"); err != nil {
		log.Fatalf("Failed to load email templates: %v", err)
	}
	r := setupRouter()
	r.Run(":8080")
	defer models.DB().Close()
}
