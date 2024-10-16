package email

import (
	"darco/proto/models/settings"
	"errors"
	"io/fs"
	"os"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/edgedb/edgedb-go"
	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
)

func loadEmailConfig(path string) (settings.EmailSettingsInput, error) {
	var config settings.EmailSettingsInput
	viper.SetConfigFile(path)
	if err := viper.ReadInConfig(); err != nil {
		return config, err
	}
	err := viper.Unmarshal(&config)
	return config, err
}

type EmailSetupArgs struct {
	NoAuto bool
	Skip   bool
}

type EmailSetup struct {
	settings.EmailSettingsInput
	connectionOK bool
	EmailSetupArgs
}

const SMTP_CONFIG_PATH = "./config/email.yaml"

func SetupEmailConfig(db edgedb.Executor, args EmailSetupArgs) {
	var emailConfig settings.EmailSettingsInput
	_, err := os.Stat(SMTP_CONFIG_PATH)
	if err != nil && errors.Is(err, fs.ErrNotExist) {
		emailConfig = settings.EmailSettingsInput{}
		emailConfig.WriteYAML(SMTP_CONFIG_PATH)
		logrus.Infof(
			"No existing SMTP config file. Generated empty config @ %s",
			SMTP_CONFIG_PATH)
	} else {
		emailConfig, err = loadEmailConfig(SMTP_CONFIG_PATH)
		if err != nil {
			logrus.Fatalf("Failed to load SMTP configuration: %v", err)
		}
		logrus.Infof("Existing SMTP configuration loaded from %s",
			SMTP_CONFIG_PATH)
	}

	var setup = EmailSetup{
		EmailSettingsInput: emailConfig,
		connectionOK:       false,
		EmailSetupArgs:     args,
	}
	_, err = tea.NewProgram(initialModel(&setup)).Run()
	if err != nil {
		logrus.Fatalf("SMTP configuration error: %v", err)
	}

	if setup.Skip {
		logrus.Infof("SMTP configuration skipped.")
		return
	}

	if !setup.connectionOK {
		logrus.Fatalf("Connection failed")
	}
	logrus.Infof(successStyle.Render("🟢 Connection succeeded"))

	if _, err := emailConfig.Save(db); err != nil {
		logrus.Fatalf(
			"Failed to save SMTP configuration in DB settings: %v",
			err)
	}

	if setup.EmailSettingsInput.WriteYAML(SMTP_CONFIG_PATH) != nil {
		logrus.Fatalf("Failed to write SMTP configuration to file: %v", err)
	}

	logrus.Infof(successStyle.Render(
		"💾 SMTP configuration saved to database and updated in config file",
	))
}
