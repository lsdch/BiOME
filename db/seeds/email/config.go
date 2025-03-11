package email

import (
	"errors"
	"fmt"
	"io/fs"
	"os"

	"github.com/geldata/gel-go/geltypes"
	"github.com/lsdch/biome/models/settings"

	tea "github.com/charmbracelet/bubbletea"
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

func SetupEmailConfig(db geltypes.Executor, args EmailSetupArgs) error {
	var emailConfig settings.EmailSettingsInput
	_, err := os.Stat(SMTP_CONFIG_PATH)
	if err != nil && errors.Is(err, fs.ErrNotExist) {
		emailConfig = settings.EmailSettingsInput{}
		if err := emailConfig.WriteYAML(SMTP_CONFIG_PATH); err != nil {
			return fmt.Errorf("Failed to create email config file: %v", err)
		}
		logrus.Infof(
			"No existing SMTP config file. Generated empty config @ %s",
			SMTP_CONFIG_PATH)
	} else {
		emailConfig, err = loadEmailConfig(SMTP_CONFIG_PATH)
		if err != nil {
			return fmt.Errorf("Failed to load SMTP configuration: %v", err)
		}
		logrus.Infof("Existing SMTP configuration loaded from %s : %+v",
			SMTP_CONFIG_PATH, emailConfig)
	}

	var setup = EmailSetup{
		EmailSettingsInput: emailConfig,
		connectionOK:       false,
		EmailSetupArgs:     args,
	}
	_, err = tea.NewProgram(initialModel(&setup)).Run()
	if err != nil {
		return fmt.Errorf("SMTP configuration error: %v", err)
	}

	if setup.Skip {
		logrus.Infof("SMTP configuration skipped.")
		return nil
	}

	if !setup.connectionOK {
		logrus.Fatalf("Connection failed")
	}
	logrus.Info(successStyle.Render("🟢 Connection succeeded"))

	if _, err := emailConfig.Save(db); err != nil {
		return fmt.Errorf(
			"Failed to save SMTP configuration in DB settings: %v",
			err)
	}

	if setup.EmailSettingsInput.WriteYAML(SMTP_CONFIG_PATH) != nil {
		return fmt.Errorf("Failed to write SMTP configuration to file: %v", err)
	}

	logrus.Info(successStyle.Render(
		"💾 SMTP configuration saved to database and updated in config file",
	))
	return nil
}
