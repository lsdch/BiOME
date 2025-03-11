package settings

import (
	"context"
	"encoding/json"

	"github.com/geldata/gel-go/geltypes"
	"github.com/lsdch/biome/models"
)

type InstanceSettingsInner struct {
	Name               string `gel:"name" json:"name" doc:"The name of this database platform" fake:"{word}"`
	IsPublic           bool   `gel:"public" json:"public" doc:"Whether the platform is accessible to unauthenticated users"`
	AllowContribSignup bool   `gel:"allow_contributor_signup" json:"allow_contributor_signup" doc:"Whether requests to contribute to the database can be made."`
}

type InstanceSettingsInput struct {
	InstanceSettingsInner `gel:"$inline" json:",inline"`
	Description           models.OptionalNull[string] `json:"description,omitempty"`
}

type InstanceSettings struct {
	ID                    geltypes.UUID `gel:"id" json:"-"`
	InstanceSettingsInner `gel:"$inline" json:",inline"`
	Description           geltypes.OptionalStr `gel:"description" json:"description"`
}

func (input *InstanceSettingsInput) Save(db geltypes.Executor) (*InstanceSettings, error) {
	jsonData, _ := json.Marshal(input)
	var s InstanceSettings
	if err := db.QuerySingle(context.Background(),
		`#edgeql
			with data := <json>$0
			select (
				(
					update admin::InstanceSettings set {
						name := <str>data['name'],
						description := <str>json_get(data, 'description') ?? <str>{},
						public := <bool>data['public'],
						allow_contributor_signup := <bool>data['allow_contributor_signup']
					}
				) ?? (
					insert admin::InstanceSettings {
						name := <str>data['name'],
						description := <str>json_get(data, 'description') ?? <str>{},
						public := <bool>data['public'],
						allow_contributor_signup := <bool>data['allow_contributor_signup']
					}
				)
			) { * } limit 1;
		`,
		&s, jsonData,
	); err != nil {
		return nil, err
	}
	return &s, nil
}
