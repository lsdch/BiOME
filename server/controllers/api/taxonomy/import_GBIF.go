package taxonomy

import (
	"context"
	"net/http"

	"github.com/geldata/gel-go/geltypes"
	"github.com/lsdch/biome/controllers"
	"github.com/lsdch/biome/models/taxonomy"
	GBIF "github.com/lsdch/biome/models/taxonomy/GBIF"
	"github.com/lsdch/biome/resolvers"
	"github.com/lsdch/biome/router"

	"github.com/danielgtaylor/huma/v2"
	"github.com/danielgtaylor/huma/v2/sse"
	"github.com/sirupsen/logrus"
)

type ImportGBIFCladeInput struct {
	resolvers.AuthRequired
	Body GBIF.ImportRequestGBIF
}

func ImportGBIFClade(stream *EventServer) router.Endpoint[ImportGBIFCladeInput, struct{}] {
	if !stream.Running {
		go stream.listen()
	}
	return func(ctx context.Context, input *ImportGBIFCladeInput) (*struct{}, error) {
		logrus.Infof("Received GBIF import request : %+v", input.Body)
		go GBIF.ImportTaxon(input.DB(), input.Body, stream.monitor)
		logrus.Infof(
			"Started import of taxon : [GBIF %d] with children: %v",
			input.Body.Key, input.Body.Children,
		)
		return nil, nil
	}
}

var stream = NewServer()

func init() {

	var APItag = "Taxonomy GBIF"

	gbifAPI := func(r *router.Router) router.Group {
		return r.RouteGroup("/taxonomy/gbif").
			WithTags([]string{APItag})
	}

	router.RegisterSpec(
		gbifAPI,
		"ListAnchors",
		huma.Operation{
			Path:    "/anchors",
			Method:  http.MethodGet,
			Summary: "List GBIF anchor clades",
			Errors:  []int{http.StatusInternalServerError},
		},
		controllers.ListHandler[*struct {
			resolvers.AuthResolver
		}](func(db geltypes.Executor) ([]taxonomy.TaxonWithParentRef, error) {
			return taxonomy.ListTaxa(db, taxonomy.ListFilters{IsAnchor: geltypes.NewOptionalBool(true)})
		}),
	)

	router.RegisterSpec(
		gbifAPI,
		"ImportGBIF",
		huma.Operation{
			Path:    "/import",
			Method:  http.MethodPut,
			Summary: "Import GBIF clade",
		},
		ImportGBIFClade(stream),
	)

	router.RegisterCustom(func(r *router.Router) {
		sse.Register(r.API,
			huma.Operation{
				Path:        "/import/taxonomy/monitor",
				OperationID: "MonitorGBIF",
				Method:      http.MethodGet,
				Summary:     "Monitor GBIF taxonomy imports",
				Tags:        []string{APItag},
			},
			map[string]any{
				"state": State{},
			},
			func(ctx context.Context, input *struct{}, send sse.Sender) {
				clientChan := stream.AddClient()
				go func() {
					<-ctx.Done()
					stream.ClosedClients <- clientChan
				}()
				msg, ok := <-clientChan
				if ok {
					if err := send.Data(msg); err != nil {
						logrus.Errorf("GBIF import monitoring error: %v", err)
					}
				}
			})
	})

}
