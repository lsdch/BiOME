package settings

import (
	"context"
	"darco/proto/db"

	"github.com/edgedb/edgedb-go"
	"github.com/sirupsen/logrus"
)

type Settings struct {
	ID       edgedb.UUID      `edgedb:"id" json:"-"`
	Instance InstanceSettings `edgedb:"instance" json:"instance"`
	Email    EmailSettings    `edgedb:"email" json:"email,omitempty"`
	Security SecuritySettings `edgedb:"security" json:"security"`
}

func init() {
	secretKey := generateSecretKeyJWT()
	if err := db.Client().Execute(context.Background(),
		`#edgeql
			with module admin,
			security := (select SecuritySettings limit 1) ?? (insert SecuritySettings {
				jwt_secret_key := <str>$0
			}),
			instance := (select InstanceSettings limit 1) ?? (insert InstanceSettings{}),
			select(<str>{});
		`,
		secretKey); err != nil {
		logrus.Fatalf("Failed to initialize settings or get existing ones: %v", err)
	}

	if err := db.Client().QuerySingle(context.Background(),
		`select ((select admin::Settings) ?? (insert admin::Settings {})) { ** } limit 1`,
		settings,
		secretKey,
	); err != nil {
		logrus.Fatalf("Failed to initialize settings or get existing ones: %v", err)
	}
}

var settings = new(Settings)

func Get() *Settings {

	if err := db.Client().QuerySingle(context.Background(),
		`select admin::Settings { ** } limit 1`, settings,
	); err != nil {
		logrus.Fatalf("Failed to get settings: %v", err)
	}
	return settings
}

func Security() SecuritySettings {
	return Get().Security
}

func Email() EmailSettings {
	return Get().Email
}

func Instance() InstanceSettings {
	return Get().Instance
}
