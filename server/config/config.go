package config

import (
	"log"
	"sync"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
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
	SecretKey     string        `mapstructure:"SECRET_KEY" validate:"required"`
	Emailer       EmailConfig   `mapstructure:",squash" validate:"required"`
	TokenLifetime time.Duration `mapstructure:"TOKEN_LIFETIME" validate:"required"`
	DomainName    string        `mapstructure:"DOMAIN_NAME" validate:"required"`
	Port          int           `mapstructure:"PORT" validate:"required"`
}

var (
	once   sync.Once
	config Config
)

func LoadConfig(path string) (err error) {
	once.Do(func() {
		viper.AddConfigPath(path)
		viper.SetConfigType("env")

		switch gin.Mode() {
		case gin.ReleaseMode:
			viper.SetConfigName("app")
		case gin.DebugMode, gin.TestMode:
			viper.SetConfigName("dev")
		}

		viper.AutomaticEnv()

		err = viper.ReadInConfig()
		if err != nil {
			return
		}

		err = viper.Unmarshal(&config)
		validate := validator.New()
		if err := validate.Struct(&config); err != nil {
			log.Fatalf("Missing required config variable: %v", err)
		}
	})
	return
}

func Get() *Config {
	return &config
}
