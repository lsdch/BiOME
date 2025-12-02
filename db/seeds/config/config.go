package config

import (
	"fmt"

	"github.com/geldata/gel-go/geltypes"
	"github.com/go-viper/mapstructure/v2"
	"github.com/lsdch/biome/db"
	"github.com/lsdch/biome/models/people"
	"github.com/lsdch/biome/models/settings"
	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
)

type InstanceConfig struct {
	Instance     settings.InstanceSettingsInput `json:"instance" mapstructure:"instance"`
	SuperAdmin   people.SuperAdminInput         `json:"superadmin" mapstructure:"superadmin"`
	Email        settings.EmailSettingsInput    `json:"email" mapstructure:"email"`
	Organisation people.OrganisationInput       `json:"organisation" mapstructure:"organisation"`
	Services     settings.ServicesSettingsInput `json:"services" mapstructure:"services"`
}

func (c *InstanceConfig) saveOrganisation(tx geltypes.Tx) error {
	maybeOrg, err := people.FindOrganisation(tx, c.Organisation.Code)
	if err != nil {
		if !db.IsNoData(err) {
			return err
		}
		if _, err := c.Organisation.Save(tx); err != nil {
			return fmt.Errorf("Organisation: %v", err)
		}
	}

	logrus.Infof("✅ Found existing organisation: [%s] %s", maybeOrg.Code, maybeOrg.Name)
	return nil
}

func (c *InstanceConfig) SaveTx(tx geltypes.Tx) error {
	if err := c.saveOrganisation(tx); err != nil {
		return err
	}

	if settings.CheckSettingsInitialized(tx) {
		logrus.Infof("✅ Instance settings already initialized, skipping.")
	} else {
		superAdmin, err := c.SuperAdmin.Save(tx)
		if err != nil {
			return fmt.Errorf("SuperAdmin: %v", err)
		}
		settings := settings.SettingsInput{
			Instance:     c.Instance,
			SuperAdminID: superAdmin.ID,
			Services:     c.Services,
		}

		if err := settings.SaveTx(tx); err != nil {
			return fmt.Errorf("Settings: %v", err)
		}
	}

	emailSettings := settings.Email()
	if emailSettings.Missing() {
		emailCfg, err := SetupEmailConfig(tx, EmailSetupArgs{Config: c.Email})
		if err != nil {
			return fmt.Errorf("Email setup: %v", err)
		}

		c.Email = emailCfg
	} else {
		logrus.Infof("Checking existing email settings")
		if err := emailSettings.TestConnection(); err != nil {
			logrus.Warnf("⚠️  Existing email settings are invalid: %v", err)
			emailCfg, err := SetupEmailConfig(tx, EmailSetupArgs{Config: c.Email})
			if err != nil {
				return fmt.Errorf("Email setup: %v", err)
			}
			c.Email = emailCfg
		} else {
			logrus.Infof("✅ Existing email settings are valid, skipping.")
		}
	}

	return nil
}

func LoadConfig[T any](v *viper.Viper, name string) (T, error) {
	var config T
	v.SetConfigType("yaml")
	v.SetConfigName(name)
	if err := v.ReadInConfig(); err != nil {
		return config, err
	}
	err := v.Unmarshal(&config, func(dc *mapstructure.DecoderConfig) {
		dc.TagName = "mapstructure"
		dc.DecodeHook = mapstructure.ComposeDecodeHookFunc(OptionalInputDecodeHook)
	})
	logrus.Debugf("Loaded config %+v", config)
	return config, err
}
