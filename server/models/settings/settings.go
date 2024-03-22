package settings

import (
	"context"
	"darco/proto/db"

	"github.com/edgedb/edgedb-go"
	"github.com/sirupsen/logrus"
)

type Settings struct {
	ID                  edgedb.UUID      `edgedb:"id" json:"id"`
	RegistrationEnabled bool             `edgedb:"registration_enabled" json:"registration_enabled"`
	Instance            InstanceSettings `edgedb:"instance" json:"instance"`
	Email               EmailSettings    `edgedb:"email" json:"email,omitempty"`
	Security            SecuritySettings `edgedb:"security" json:"security"`
}

var settings = new(Settings)

func init() {
	secretKey := generateSecretKeyJWT()
	if err := db.Client().QuerySingle(context.Background(),
		`with module admin
		select (
			(select Settings) ?? (insert Settings {
				instance := (select InstanceSettings limit 1) ?? (insert InstanceSettings {}),
				email := (select EmailSettings limit 1),
				security := (select SecuritySettings limit 1) ?? (insert SecuritySettings {
					jwt_secret_key := <str>$0
				})
			})
		) { ** } limit 1`,
		settings,
		secretKey,
	); err != nil {
		logrus.Fatalf("Failed to initialize settings or get existing ones: %v", err)
	}
}

func Get() Settings {
	return *settings
}

func Security() SecuritySettings {
	return (*settings).Security
}

func Email() EmailSettings {
	return (*settings).Email
}

func Instance() InstanceSettings {
	return (*settings).Instance
}
