package main

import (
	"darco/proto/controllers/taxonomy"
	"net/http"

	"github.com/gin-gonic/gin"
)

var db = make(map[string]string)

func handleErrors(c *gin.Context) {
	c.Next() // execute all the handlers

	// at this point, all the handlers finished. Let's read the errors!
	// in this example we only will use the **last error typed as public**
	// but you could iterate over all them since c.Errors is a slice!
	errorToPrint := c.Errors.Last()
	if errorToPrint != nil && errorToPrint.Meta != nil {
		c.JSON(500, errorToPrint.Meta)
	}
}

func setupRouter() *gin.Engine {
	// Disable Console Color
	// gin.DisableConsoleColor()
	r := gin.Default()
	r.Use(handleErrors)

	// r.Use(cors.New(cors.Config{
	// 	AllowOrigins:     []string{"http://localhost:5173"},
	// 	AllowHeaders:     []string{"Origin"},
	// 	ExposeHeaders:    []string{"Content-Length"},
	// 	AllowCredentials: true,
	// 	MaxAge:           12 * time.Hour,
	// }))

	// Ping test
	r.GET("/api/", func(c *gin.Context) {
		c.String(http.StatusOK, "Hello world!")
	})

	api := r.Group("/api")

	taxonomy_api := api.Group("/taxonomy")
	taxonomyUpdate := taxonomy.UpdateTaxonomyDB()
	taxonomy_api.POST("/update", taxonomyUpdate.Endpoint)
	taxonomy_api.GET("/update/progress", taxonomyUpdate.ProgressTracker)
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
