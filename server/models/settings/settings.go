package settings

import (
	"context"
	"darco/proto/db"
	"fmt"

	"github.com/edgedb/edgedb-go"
	"github.com/sirupsen/logrus"
)

type SuperAdmin struct {
	Email string `edgedb:"email" json:"email"`
	Name  string `edgedb:"name" json:"name"`
}

type Settings struct {
	ID         edgedb.UUID      `edgedb:"id" json:"-"`
	Instance   InstanceSettings `edgedb:"instance" json:"instance"`
	Email      EmailSettings    `edgedb:"email" json:"email,omitempty"`
	Security   SecuritySettings `edgedb:"security" json:"security"`
	SuperAdmin SuperAdmin       `edgedb:"superadmin" json:"superadmin"`
}

var settings = new(Settings)

func Setup(db edgedb.Executor, superAdminID edgedb.UUID) error {
	secretKey := generateSecretKeyJWT()
	if err := db.Execute(context.Background(),
		`#edgeql
			with module admin,
			security := (select SecuritySettings limit 1) ?? (insert SecuritySettings {
				jwt_secret_key := <str>$0
			}),
			instance := (select InstanceSettings limit 1) ?? (insert InstanceSettings{}),
			select(<str>{});
		`, secretKey,
	); err != nil {
		return fmt.Errorf("Failed to initialize settings: %v", err)
	}

	if err := db.QuerySingle(context.Background(),
		`#edgeql
			select (insert admin::Settings {
				superadmin := (assert_exists(<people::User><uuid>$0))
			}) { **, superadmin: { email, name := .identity.full_name } } limit 1
		`, settings, superAdminID,
	); err != nil {
		return fmt.Errorf("Failed to initialize settings: %v", err)
	}
	return nil
}

func Get() *Settings {

	if err := db.Client().QuerySingle(context.Background(),
		`#edgeql
			select admin::Settings {
				**,
				superadmin: { email, name := .identity.full_name }
			} limit 1
		`, settings,
	); err != nil {
		if db.IsNoData(err) {
			logrus.Fatalf("Settings are not initialized")
		}
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
