package main

import (
	"darco/proto/config"
	country "darco/proto/controllers/location"
	"darco/proto/controllers/taxonomy"
	accounts "darco/proto/controllers/users"
	"darco/proto/models/validations"
	"darco/proto/router"
	"darco/proto/services/email"
	"errors"
	"net/http"

	_ "darco/proto/docs" // import swagger docs

	"github.com/edgedb/edgedb-go"
	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	"github.com/go-playground/validator/v10"
	log "github.com/sirupsen/logrus"
	swaggerfiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

// Generate OpenAPI docs from swaggo doc comments
//go:generate swag init --parseDependency --parseInternal -g main.go

// Error handling middleware
func handleErrors(c *gin.Context) {
	c.Next() // execute all the handlers

	// at this point, all the handlers finished. Let's read the errors!
	// in this example we only will use the **last error typed as public**
	// but you could iterate over all them since c.Errors is a slice!
	err := c.Errors.Last()
	if err == nil {
		return
	}

	var dbErr edgedb.Error
	if errors.As(err, &dbErr) && dbErr.Category(edgedb.NoDataError) {
		c.AbortWithStatus(http.StatusNotFound)
		return
	}

	if err.Meta != nil {
		c.JSON(http.StatusInternalServerError, err.Meta)
		return
	}

	c.JSON(int(err.Type), err.Err.Error())
}

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
	r.Use(handleErrors)

	// Ping test
	r.GET("/api/", func(c *gin.Context) {
		c.String(http.StatusOK, "Hello world!")
	})

	// Swagger docs
	api := r.Group(router.Config.BasePath)
	api.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerfiles.Handler))

	country_api := api.Group("/countries")
	country_api.GET("/", country.List)
	country_api.GET("/setup", country.Setup)

	taxonomy_api := api.Group("/taxonomy")
	taxonomy_api.GET("/", taxonomy.ListTaxa)
	taxonomy_api.GET("/:code", taxonomy.GetTaxon)
	taxonomy_api.DELETE("/:code", taxonomy.DeleteTaxon)
	taxonomy_api.PATCH("/:code", taxonomy.UpdateTaxon)
	importGBIF := taxonomy.ImportCladeGBIF()
	taxonomy_api.PUT("/import", importGBIF.Endpoint)
	taxonomy_api.GET("/import", importGBIF.ProgressTracker)
	taxonomy_api.GET("/anchors", taxonomy.GetAnchors)

	users_api := api.Group("/users")
	users_api.POST("/register", accounts.Register)
	users_api.GET("/confirm", accounts.ConfirmEmail)
	users_api.POST("/confirm/resend", accounts.ResendConfirmation)
	users_api.POST("/forgotten-password", accounts.RequestPasswordReset)

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
	err := config.LoadConfig(".")
	if err != nil {
		log.Fatalf("Failed to load configuration file : %v", err)
	}
	log.Debugf("Config loaded : %+v", config.Get())

}

func main() {
	if v, ok := binding.Validator.Engine().(*validator.Validate); ok {
		v.RegisterValidation("unique_email", validations.EmailUniqueValidator)
	}
	gin.ForceConsoleColor()

	if gin.Mode() == gin.DebugMode {
		log.SetLevel(log.DebugLevel)
	}

	loadConfig()
	if err := email.LoadTemplates("templates/**"); err != nil {
		log.Fatalf("Failed to load email templates: %v", err)
	}
	r := setupRouter()
	r.Run(":8080")
}
