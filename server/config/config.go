package config

import (
	"darco/proto/router"
	"fmt"
	"log"
	"net/url"
	"path"
	"path/filepath"
	"sync"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
)

type EmailConfig struct {
	From          string        `mapstructure:"EMAIL_FROM" validate:"required"`
	Host          string        `mapstructure:"SMTP_HOST" validate:"required"`
	User          string        `mapstructure:"SMTP_USER" validate:"required"`
	Pass          string        `mapstructure:"SMTP_PASS" validate:"required"`
	Port          int           `mapstructure:"SMTP_PORT" validate:"required"`
	TokenLifetime time.Duration `mapstructure:"EMAIL_TOKEN_LIFETIME" validate:"required"`
}

type Config struct {
	SecretKey     string         `mapstructure:"SECRET_KEY" validate:"required"`
	Emailer       EmailConfig    `mapstructure:",squash" validate:"required"`
	TokenLifetime time.Duration  `mapstructure:"TOKEN_LIFETIME" validate:"required"`
	DomainName    string         `mapstructure:"DOMAIN_NAME" validate:"required"`
	Port          int            `mapstructure:"PORT" validate:"required"`
	Client        ClientConfig   `mapstructure:",squash"`
	Accounts      AccountsConfig `mapstructure:",squash"`
}

type AccountsConfig struct {
	PasswordStrength int `mapstructure:"PASSWORD_MIN_STRENGTH" validate:"min=1,max=5"`
}

type ClientConfig struct {
	AppName    string `mapstructure:"VITE_APP_NAME"`
	DomainName string `mapstructure:"VITE_DOMAIN" validate:"required"`
	Port       int    `mapstructure:"VITE_PORT" validate:"required"`
}

var (
	once sync.Once

	// App configuration as a struct with some default values
	config = Config{
		Accounts: AccountsConfig{
			PasswordStrength: 3,
		},
	}
)

// Generates a URL to an endpoint on the server.
//
// URL path must omit the API prefix, e.g. '/api/v1/users' must be '/users/'
func (config *Config) MakeURL(url_path string) url.URL {
	return url.URL{
		Scheme: "https",
		Host:   fmt.Sprintf("%s:%d", config.DomainName, config.Port),
		Path:   path.Join(router.Config.BasePath, url_path),
	}
}

// Generates a URL to a client page.
func (config *Config) MakeClientURL(url_path string) url.URL {
	return url.URL{
		Scheme: "https",
		Host:   fmt.Sprintf("%s:%d", config.Client.DomainName, config.Client.Port),
		Path:   url_path,
	}
}

func LoadConfig(path string) (*Config, error) {
	var err error = nil
	once.Do(func() {
		viper.AddConfigPath(path)
		viper.SetConfigType("env")
		var clientConfigFile string
		switch gin.Mode() {
		case gin.ReleaseMode:
			viper.SetConfigName("app")
			clientConfigFile = filepath.Join("../client", ".env.production")
		case gin.DebugMode, gin.TestMode:
			viper.SetConfigName("dev")
			clientConfigFile = filepath.Join("../client", ".env.development")
		}

		viper.AutomaticEnv()

		if err = viper.ReadInConfig(); err != nil {
			return
		}

		logrus.Debugf("Loading client config file %s", clientConfigFile)
		viper.SetConfigName(clientConfigFile)
		viper.AutomaticEnv()
		err = viper.MergeInConfig()

		err = viper.Unmarshal(&config)
		validate := validator.New()
		if err = validate.Struct(&config); err != nil {
			log.Fatalf("Invalid configuration: %v", err)
		}
	})
	return &config, err
}

func Get() *Config {
	return &config
}
