package settings

import (
	"context"
	"encoding/json"

	"github.com/edgedb/edgedb-go"
)

type InstanceSettingsInner struct {
	Name               string `edgedb:"name" json:"name" faker:"word"`
	IsPublic           bool   `edgedb:"public" json:"public"`
	AllowContribSignup bool   `edgedb:"allow_contributor_signup" json:"allow_contributor_signup"`
}

type InstanceSettingsInput struct {
	InstanceSettingsInner `edgedb:"$inline" json:",inline"`
	Description           *string `json:"description,omitnil"`
}

type InstanceSettings struct {
	ID                    edgedb.UUID `edgedb:"id" json:"id"`
	InstanceSettingsInner `edgedb:"$inline" json:",inline"`
	Description           edgedb.OptionalStr `edgedb:"description" json:"description"`
}

func (input *InstanceSettingsInner) Save(db edgedb.Executor) (*InstanceSettings, error) {
	jsonData, _ := json.Marshal(input)
	var settings InstanceSettings
	if err := db.QuerySingle(context.Background(),
		`with data := <json>$0
			select (update admin::InstanceSettings set {
				name := <str>data['name'],
				description := <str>json_get(data, 'description') ?? <str>{},
				public := <bool>data['public'],
				allow_contributor_signup := <bool>data['allow_contributor_signup']
			}) { * } limit 1`,
		&settings,
		jsonData,
	); err != nil {
		return nil, err
	}
	return &settings, nil
}
