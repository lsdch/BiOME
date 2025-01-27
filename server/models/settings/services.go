package settings

import (
	"context"
	"encoding/json"

	"github.com/lsdch/biome/db"
	"github.com/lsdch/biome/models"

	"github.com/edgedb/edgedb-go"
)

type ServiceSettings struct {
	GeoapifyApiKey edgedb.OptionalStr `edgedb:"geoapify_api_key" json:"geoapify_api_key"`
}

type ServiceSettingsUpdate struct {
	GeoapifyApiKey models.OptionalNull[string] `edgedb:"geoapify_api_key" json:"geoapify_api_key"`
}

func (u ServiceSettingsUpdate) Save(e edgedb.Executor) (updated ServiceSettings, err error) {
	data, _ := json.Marshal(u)
	query := db.UpdateQuery{
		Frame: `#edgeql
			with item := <json>$0,
			select (update admin::Settings set {
				%s
			}) { geoapify_api_key }
			limit 1
		`,
		Mappings: map[string]string{
			"geoapify_api_key": "<str>item['geoapify_api_key']",
		},
	}
	err = e.QuerySingle(context.Background(), query.Query(u), &updated, data)
	return
}
