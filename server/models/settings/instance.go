package settings

import (
	"context"
	"encoding/json"

	"github.com/geldata/gel-go/geltypes"
	"github.com/lsdch/biome/db"
	"github.com/lsdch/biome/models"
	"github.com/sirupsen/logrus"
)

// @mapstructure
type InstanceSettingsInner struct {
	Name               string `gel:"name" json:"name" maxLength:"20" doc:"The name of this database platform" fake:"{word}" mapstructure:"name"`
	IsPublic           bool   `gel:"public" json:"public" doc:"Whether the platform is accessible to unauthenticated users" mapstructure:"public"`
	AllowContribSignup bool   `gel:"allow_contributor_signup" json:"allow_contributor_signup" doc:"Whether requests to contribute to the database can be made." mapstructure:"allow_contributor_signup"`
}

// @mapstructure
type InstanceSettingsInput struct {
	InstanceSettingsInner `gel:"$inline" mapstructure:",squash"`
	Description           models.OptionalNull[string] `gel:"description" json:"description,omitempty" mapstructure:"description"`
}

type InstanceSettings struct {
	ID                    geltypes.UUID `gel:"id" json:"-"`
	InstanceSettingsInner `gel:"$inline"`
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

type InstanceSettingsUpdate struct {
	Name               models.OptionalInput[string] `gel:"name" json:"name,omitempty" maxLength:"20" fake:"{word}" mapstructure:"name"`
	IsPublic           models.OptionalInput[bool]   `gel:"public" json:"public,omitempty" doc:"Whether the platform is accessible to unauthenticated users" mapstructure:"public"`
	AllowContribSignup models.OptionalInput[bool]   `gel:"allow_contributor_signup" json:"allow_contributor_signup,omitempty" doc:"Whether requests to contribute to the database can be made."`
	Description        models.OptionalNull[string]  `gel:"description" json:"description,omitempty" mapstructure:"description"`
}

func (u InstanceSettingsUpdate) Save(e geltypes.Executor) (updated InstanceSettings, err error) {
	data, _ := json.Marshal(u)
	logrus.Debugf("Saving instance settings with data: %s", string(data))
	query := db.UpdateQuery{
		Frame: `#edgeql
			with item := <json>$0,
			select (update admin::InstanceSettings set {
				%s
			}) { * } limit 1;
			`,
		Mappings: map[string]string{
			"name":                     "<str>item['name']",
			"description":              "<str>item['description']",
			"public":                   "<bool>item['public']",
			"allow_contributor_signup": "<bool>item['allow_contributor_signup']",
		},
	}
	logrus.Debugf("Generated query: %s", query.Query(u))
	err = e.QuerySingle(context.Background(), query.Query(u), &updated, data)
	logrus.Debugf("Updated instance settings: %+v", updated)
	return
}
