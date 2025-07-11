package settings

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/geldata/gel-go/geltypes"
	"github.com/lsdch/biome/db"

	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
)

type SuperAdmin struct {
	Email string `gel:"email" json:"email"`
	Name  string `gel:"name" json:"name"`
}

type Settings struct {
	ID              geltypes.UUID    `gel:"id" json:"-"`
	Instance        InstanceSettings `gel:"instance" json:"instance"`
	Email           EmailSettings    `gel:"email" json:"email,omitempty"`
	Security        SecuritySettings `gel:"security" json:"security"`
	SuperAdmin      SuperAdmin       `gel:"superadmin" json:"superadmin"`
	ServiceSettings `gel:"$inline" json:"services"`
}

var settings = new(Settings)

type SettingsInput struct {
	Instance       InstanceSettingsInput `json:"instance"`
	SuperAdminID   geltypes.UUID         `json:"super_admin_id"`
	GeoapifyApiKey *string               `json:"geoapify_api_key,omitempty" map_structure:"GEOAPIFY_API_KEY"`
}

func (i *SettingsInput) LoadConfig(path string) error {
	viper.SetConfigFile(path)
	if err := viper.ReadInConfig(); err != nil {
		return err
	}
	err := viper.Unmarshal(i)
	return err
}

func (i SettingsInput) SaveTx(tx geltypes.Tx) error {
	// Init security settings with JWT secret key
	secretKey := generateSecretKeyJWT()
	if err := tx.Execute(context.Background(),
		`#edgeql
			insert admin::SecuritySettings { jwt_secret_key := <str>$0 }
		`, secretKey,
	); err != nil {
		return fmt.Errorf("Failed to initialize security settings: %v", err)
	}

	if _, err := i.Instance.Save(tx); err != nil {
		return fmt.Errorf("Failed to initialize instance settings: %v", err)
	}

	data, err := json.Marshal(i)
	if err != nil {
		return err
	}

	if err := tx.QuerySingle(context.Background(),
		`#edgeql
			with data := <json>$0
			select (insert admin::Settings {
				superadmin := (assert_exists(
					<people::User><uuid>data['super_admin_id'],
					message := 'Super admin not found'
				)),
				geoapify_api_key := <str>json_get(data, 'geoapify_api_key')
			}) {
				**,
				superadmin: { email, name := .identity.full_name }
			} limit 1
		`, settings, data,
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

func Services() ServiceSettings {
	return Get().ServiceSettings
}
