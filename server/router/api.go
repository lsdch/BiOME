package router

import (
	"bytes"
	"context"
	"darco/proto/models/location"
	"encoding/json"
	"net/http"
	"os"
	"path"
	"reflect"
	"slices"
	"time"

	"github.com/danielgtaylor/huma/v2"
	"github.com/danielgtaylor/huma/v2/adapters/humagin"
	"github.com/edgedb/edgedb-go"
	"github.com/gin-gonic/gin"
)

type Router struct {
	BasePath string
	API      huma.API
	Engine   *gin.Engine
	BaseAPI  *gin.RouterGroup
	config   huma.Config
}

func New(r *gin.Engine, basePath string, config huma.Config) Router {
	baseAPI := r.Group(basePath)
	API := humagin.NewWithGroup(r, baseAPI, config)

	registry := huma.NewMapRegistry("#/components/schemas/", huma.DefaultSchemaNamer)
	registry.RegisterTypeAlias(reflect.TypeFor[edgedb.OptionalStr](), reflect.TypeOf(""))
	registry.RegisterTypeAlias(reflect.TypeFor[edgedb.OptionalInt32](), reflect.TypeOf(0))
	registry.RegisterTypeAlias(reflect.TypeFor[edgedb.OptionalBool](), reflect.TypeOf(true))
	registry.RegisterTypeAlias(reflect.TypeFor[edgedb.OptionalDateTime](), reflect.TypeFor[time.Time]())
	registry.RegisterTypeAlias(reflect.TypeFor[location.OptionalHabitatRecord](), reflect.TypeFor[location.HabitatRecord]())
	API.OpenAPI().Components.Schemas = registry

	return Router{
		BasePath: basePath,
		API:      API,
		BaseAPI:  baseAPI,
		Engine:   r,
		config:   config,
	}
}

func (r *Router) WriteSpecJSON(outputPath string) error {
	file, err := os.Create(outputPath)
	if err != nil {
		return err
	}
	var indentedJSON bytes.Buffer
	bytes, err := r.API.OpenAPI().MarshalJSON()
	if err != nil {
		return err
	}
	if err = json.Indent(&indentedJSON, bytes, "", "\t"); err != nil {
		return err
	}
	_, err = file.Write(indentedJSON.Bytes())
	return err
}

func (r *Router) RouteGroup(prefix string) group {

	return group{r, r.API, prefix, []string{}}
}

type group struct {
	router *Router
	API    huma.API
	Prefix string
	Tags   []string
}

func (g group) WithTags(tags []string) group {
	g.Tags = tags
	return g
}

type Endpoint[I, O any] func(context.Context, *I) (*O, error)

func Register[I, O any](
	group group,
	operationID string,
	op huma.Operation,
	handler func(context.Context, *I) (*O, error),
) string {
	ingroupOp := op
	ingroupOp.OperationID = operationID
	ingroupOp.Path = path.Join(group.Prefix, op.Path)
	ingroupOp.Tags = slices.Concat(ingroupOp.Tags, group.Tags)
	ingroupOp.Errors = append(op.Errors, http.StatusInternalServerError)
	huma.Register(group.API, ingroupOp, handler)
	return path.Join(group.router.BasePath, ingroupOp.Path)
}
