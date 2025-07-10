package presets

import (
	"context"
	"encoding/json"

	"github.com/geldata/gel-go/geltypes"
	"github.com/lsdch/biome/models"
	"github.com/lsdch/biome/models/people"
)

type PresetsInner struct {
	Name     string `json:"name" gel:"name"`
	Spec     string `json:"spec" gel:"spec" doc:"A JSON representation of settings."`
	IsGlobal bool   `json:"is_global" gel:"is_global" doc:"Global settings are considered as part of the application."`
	IsPublic bool   `json:"is_public" gel:"is_public" doc:"Public settings are available to all users, while private settings are only available to the user who created them."`
}

type PresetsInput struct {
	ID           models.OptionalInput[geltypes.UUID] `json:"id,omitempty" gel:"id"`
	PresetsInner `json:",inline" gel:"$inline"`
	Description  models.OptionalInput[string] `json:"description,omitempty" gel:"description"`
}

type DataFeedSpecInput PresetsInput
type MapToolPresetInput PresetsInput

type Presets struct {
	ID           geltypes.UUID        `edgedb:"id" json:"id" format:"uuid"`
	Description  geltypes.OptionalStr `json:"description,omitempty" gel:"description"`
	PresetsInput `json:",inline" gel:"$inline"`
	Meta         people.Meta `json:"meta" gel:"meta"`
}

type DataFeedSpec Presets
type MapToolPreset Presets

func ListDataFeedSpecs(db geltypes.Executor) (specs []DataFeedSpec, err error) {
	err = db.Query(context.Background(),
		`#edgeql
			select settings::DataFeedSpec { *, meta: { * } }
			filter .is_public or (.meta.created_by_user = global current_user) ?? false
		`, &specs)
	return
}

func ListMapPresets(db geltypes.Executor) (specs []MapToolPreset, err error) {
	err = db.Query(context.Background(),
		`#edgeql
			select settings::MapToolPreset { *, meta: { * } }
			filter .is_public or (.meta.created_by_user = global current_user) ?? false
		`, &specs)
	return
}

func GetDataFeedSpec(db geltypes.Executor, name string) (spec DataFeedSpec, err error) {
	err = db.QuerySingle(context.Background(),
		`#edgeql
			select settings::DataFeedSpec { *, meta: { * } }
			filter .name = <str>$name
		`, &spec, name)
	return
}

func GetMapPreset(db geltypes.Executor, name string) (spec MapToolPreset, err error) {
	err = db.QuerySingle(context.Background(),
		`#edgeql
			select settings::MapToolPreset { *, meta: { * } }
			filter .name = <str>$name
		`, &spec, name)
	return
}

func DeleteMapPreset(db geltypes.Executor, name string) (deleted MapToolPreset, err error) {
	err = db.QuerySingle(context.Background(),
		`#edgeql
			select (
				delete settings::MapToolPreset filter .name = <str>$0
		 	) {  *, meta: { * } };`,
		&deleted, name)
	return
}

func DeleteDataFeedSpec(db geltypes.Executor, name string) (deleted DataFeedSpec, err error) {
	err = db.QuerySingle(context.Background(),
		`#edgeql
			select (
				delete settings::DataFeedSpec filter .name = <str>$0
		 	) {  *, meta: { * } };`,
		&deleted, name)
	return
}

func (i DataFeedSpecInput) Save(e geltypes.Executor) (created DataFeedSpec, err error) {
	data, _ := json.Marshal(i)
	err = e.QuerySingle(context.Background(), `#edgeql
			with data := <json>$0,
			select (
				(
					update settings::DataFeedSpec
					filter .id = <uuid>json_get(data, 'id')
					set {
						name := <str>data['name'],
						description := <str>json_get(data, 'description'),
						spec := data['spec'],
						is_global := <bool>json_get(data, 'is_global'),
						is_public := <bool>json_get(data, 'is_public'),
					}
				) ?? (
					insert settings::DataFeedSpec {
						name := <str>data['name'],
						description := <str>json_get(data, 'description'),
						spec := data['spec'],
						is_global := <bool>json_get(data, 'is_global'),
						is_public := <bool>json_get(data, 'is_public'),
				}
			)
		) { *, meta: { * } }
		`, &created, data)
	return
}

func (i MapToolPresetInput) Save(e geltypes.Executor) (created MapToolPreset, err error) {
	data, _ := json.Marshal(i)
	err = e.QuerySingle(context.Background(),
		`#edgeql
			with data := <json>$0,
			select (
				(
					update settings::MapToolPreset
					filter .id = <uuid>json_get(data, 'id')
					set {
						name := <str>data['name'],
						description := <str>json_get(data, 'description'),
						spec := data['spec'],
						is_global := <bool>json_get(data, 'is_global'),
						is_public := <bool>json_get(data, 'is_public'),
					}
				) ?? (
					insert settings::MapToolPreset {
						name := <str>data['name'],
						description := <str>json_get(data, 'description'),
						spec := data['spec'],
						is_global := <bool>json_get(data, 'is_global'),
						is_public := <bool>json_get(data, 'is_public'),
				}
			)
		) { *, meta: { * } }
		`, &created, data)
	return
}
