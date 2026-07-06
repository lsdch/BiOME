package router

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"os"
	"path"
	"path/filepath"
	"reflect"
	"slices"

	"github.com/lsdch/biome/lib/auth"
	"github.com/lsdch/biome/middleware"
	"github.com/lsdch/biome/services/crossref"

	"github.com/caltechlibrary/crossrefapi"
	"github.com/danielgtaylor/huma/v2"
	"github.com/danielgtaylor/huma/v2/adapters/humagin"
	"github.com/gin-gonic/gin"
)

type Router struct {
	BasePath string
	API      huma.API
	Engine   *gin.Engine
	BaseAPI  *gin.RouterGroup
	config   huma.Config
}

func SchemaNamer(t reflect.Type, hint string) string {
	return huma.DefaultSchemaNamer(t, hint)
}

func New(r *gin.Engine, basePath string, config huma.Config) Router {
	baseAPI := r.Group(basePath)
	API := humagin.NewWithGroup(r, baseAPI, config)

	registry := huma.NewMapRegistry("#/components/schemas/", SchemaNamer)

	registry.RegisterTypeAlias(reflect.TypeFor[crossrefapi.Person](), reflect.TypeFor[crossref.CrossRefPerson]())
	registry.RegisterTypeAlias(reflect.TypeFor[crossrefapi.DateRange](), reflect.TypeFor[crossref.CrossRefDateRange]())

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
	if err := os.MkdirAll(filepath.Dir(outputPath), os.ModeDir); err != nil {
		return err
	}
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

func (r *Router) RouteGroup(prefix string) Group {
	return Group{r, r.API, prefix, []string{}, auth.Public()}
}

type Group struct {
	router       *Router
	API          huma.API
	Prefix       string
	Tags         []string
	AccessPolicy auth.Policy
}

func (g Group) WithTags(tags []string) Group {
	g.Tags = tags
	return g
}

func (g Group) WithAccessPolicy(policy auth.Policy) Group {
	g.AccessPolicy = policy
	return g
}

func (g Group) RouteGroup(prefix string) Group {
	return Group{router: g.router, API: g.API, Prefix: path.Join(g.Prefix, prefix), Tags: g.Tags}
}

type Endpoint[I, O any] func(context.Context, *I) (*O, error)

func Register[I, O any](
	group Group,
	operationID string,
	op huma.Operation,
	handler func(context.Context, *I) (*O, error),
) string {
	ingroupOp := op
	ingroupOp.OperationID = operationID
	ingroupOp.Path = path.Join(group.Prefix, op.Path)
	ingroupOp.Tags = slices.Concat(ingroupOp.Tags, group.Tags)
	ingroupOp.Errors = append(op.Errors, http.StatusInternalServerError)
	// If the operation doesn't have an access policy, use the group's access policy
	if _, ok := middleware.GetAccessPolicy(ingroupOp); !ok && group.AccessPolicy != nil {
		ingroupOp = middleware.WithAccessPolicy(ingroupOp, group.AccessPolicy)
	}
	huma.Register(group.API, ingroupOp, handler)
	return path.Join(group.router.BasePath, ingroupOp.Path)
}

type RouteSpec interface {
	Register(r *Router) RouteSpec
	Path(r *Router) string
	WithAccessPolicy(policy auth.Policy) RouteSpec
}

type routeSpec[I, O any] struct {
	Group       Group
	OperationID string
	huma.Operation
	Handler func(context.Context, *I) (*O, error)
}

func (spec routeSpec[I, O]) Register(r *Router) RouteSpec {
	Register(spec.Group, spec.OperationID, spec.Operation, spec.Handler)
	return spec
}

func (spec routeSpec[I, O]) WithAccessPolicy(policy auth.Policy) RouteSpec {
	spec.Operation = middleware.WithAccessPolicy(spec.Operation, policy)
	return spec
}

func (spec routeSpec[I, O]) Path(r *Router) string {
	return path.Join(spec.Group.router.BasePath, spec.Group.Prefix, spec.Operation.Path)
}

func NewSpec[I, O any](
	group Group,
	operationID string,
	op huma.Operation,
	handler func(context.Context, *I) (*O, error),
) RouteSpec {
	return routeSpec[I, O]{
		Group:       group,
		OperationID: operationID,
		Operation:   op,
		Handler:     handler,
	}
}
