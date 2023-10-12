package main

import (
	country "darco/proto/controllers/location"
	"darco/proto/controllers/taxonomy"
	"errors"
	"net/http"

	_ "darco/proto/docs" // import swagger docs

	"github.com/edgedb/edgedb-go"
	"github.com/gin-gonic/gin"
	log "github.com/sirupsen/logrus"
	swaggerfiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

var db = make(map[string]string)

func handleErrors(c *gin.Context) {
	c.Next() // execute all the handlers

	// at this point, all the handlers finished. Let's read the errors!
	// in this example we only will use the **last error typed as public**
	// but you could iterate over all them since c.Errors is a slice!
	errorToPrint := c.Errors.Last()
	if errorToPrint == nil {
		return
	}

	var dbErr edgedb.Error
	if errors.As(errorToPrint, &dbErr) && dbErr.Category(edgedb.NoDataError) {
		c.AbortWithStatus(http.StatusNotFound)
		return
	}

	if errorToPrint.Meta != nil {
		c.JSON(http.StatusInternalServerError, errorToPrint.Meta)
		return
	}
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
// @name token// @license.name Apache 2.0
// @license.url http://www.apache.org/licenses/LICENSE-2.0.html
// @host localhost:8080
// @schemes http
func setupRouter() *gin.Engine {

	r := gin.Default()
	r.Use(handleErrors)

	if gin.Mode() == gin.DebugMode {
		log.SetLevel(log.DebugLevel)
	}

	// Ping test
	r.GET("/api/", func(c *gin.Context) {
		c.String(http.StatusOK, "Hello world!")
	})

	// Swagger docs
	api := r.Group("/api/v1")
	api.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerfiles.Handler))

	country_api := api.Group("/countries")
	country_api.GET("/", country.List)
	country_api.GET("/setup", country.Setup)

	taxa_api := api.Group("/taxa")
	taxa_api.GET("/", taxonomy.ListTaxa)
	taxa_api.GET("/:code", taxonomy.GetTaxon)
	taxa_api.DELETE("/:code", taxonomy.DeleteTaxon)
	taxa_api.PATCH("/:code", taxonomy.UpdateTaxon)
	taxonomy_api := api.Group("/taxonomy")
	importGBIF := taxonomy.ImportCladeGBIF()
	taxonomy_api.PUT("/import", importGBIF.Endpoint)
	taxonomy_api.GET("/import", importGBIF.ProgressTracker)
	taxonomy_api.GET("/anchors", taxonomy.GetAnchors)

	// Get user value
	r.GET("/user/:name", func(c *gin.Context) {
		user := c.Params.ByName("name")
		value, ok := db[user]
		if ok {
			c.JSON(http.StatusOK, gin.H{"user": user, "value": value})
		} else {
			c.JSON(http.StatusOK, gin.H{"user": user, "status": "no value"})
		}
	})

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

	/* example curl for /admin with basicauth header
	   Zm9vOmJhcg== is base64("foo:bar")

		curl -X POST \
	  	http://localhost:8080/admin \
	  	-H 'authorization: Basic Zm9vOmJhcg==' \
	  	-H 'content-type: application/json' \
	  	-d '{"value":"bar"}'
	*/
	authorized.POST("admin", func(c *gin.Context) {
		user := c.MustGet(gin.AuthUserKey).(string)

		// Parse JSON
		var json struct {
			Value string `json:"value" binding:"required"`
		}

		if c.Bind(&json) == nil {
			db[user] = json.Value
			c.JSON(http.StatusOK, gin.H{"status": "ok"})
		}
	})

	return r
}

func main() {
	r := setupRouter()
	r.Run(":8080")
}
