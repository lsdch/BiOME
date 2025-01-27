package taxonomy

import (
	"context"
	"net/http"

	"github.com/lsdch/biome/controllers"
	"github.com/lsdch/biome/models/taxonomy"
	GBIF "github.com/lsdch/biome/models/taxonomy/GBIF"
	"github.com/lsdch/biome/resolvers"
	"github.com/lsdch/biome/router"

	"github.com/danielgtaylor/huma/v2"
	"github.com/danielgtaylor/huma/v2/sse"
	"github.com/edgedb/edgedb-go"
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

func RegisterImportRoutes(r router.Router) {
	var stream = NewServer()
	var APItag = "Taxonomy GBIF"

	huma.Register(r.API,
		huma.Operation{
			Path:        "/anchors/",
			OperationID: "ListAnchors",
			Method:      http.MethodGet,
			Summary:     "List GBIF anchor clades",
			Tags:        []string{APItag},
			Errors:      []int{http.StatusInternalServerError},
		}, controllers.ListHandler[*struct {
			resolvers.AuthResolver
		}](func(db edgedb.Executor) ([]taxonomy.TaxonWithParentRef, error) {
			return taxonomy.ListTaxa(db, taxonomy.ListFilters{IsAnchor: edgedb.NewOptionalBool(true)})
		}))

	huma.Register(r.API, huma.Operation{
		Path:        "/import/taxonomy",
		OperationID: "ImportGBIF",
		Method:      http.MethodPut,
		Summary:     "Import GBIF clade",
		Tags:        []string{APItag},
	}, ImportGBIFClade(stream))

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
}
