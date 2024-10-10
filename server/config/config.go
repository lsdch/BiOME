package config

import (
	"fmt"
	"time"

	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
)

type Config struct {
	AuthTokenLifetimeMinutes  uint `mapstructure:"JWT_LIFETIME_MINUTES"`
	AccountTokenLifetimeHours uint `mapstructure:"ACCOUNT_TOKEN_LIFETIME_HOURS"`
}

var cfg = Config{
	AuthTokenLifetimeMinutes:  30,
	AccountTokenLifetimeHours: 24,
}

func LoadConfig(dir string, name string) (Config, error) {
	viper.AddConfigPath(dir)
	viper.SetConfigName(name)
	viper.SetConfigType("yaml")

	viper.AutomaticEnv()

	if err := viper.ReadInConfig(); err != nil {
		return cfg, err
	}

	err := viper.Unmarshal(&cfg)
	return cfg, err
}

func Get() Config {
	return cfg
}

func (c Config) AuthTokenDuration() time.Duration {
	d, err := time.ParseDuration(
		fmt.Sprintf("%dm", c.AuthTokenLifetimeMinutes),
	)
	if err != nil {
		logrus.Fatalf("Failed to parse auth token duration: %v", err)
	}
	return d
}

func (c Config) AccountTokenDuration() time.Duration {
	d, err := time.ParseDuration(
		fmt.Sprintf("%dh", c.AccountTokenLifetimeHours),
	)
	if err != nil {
		logrus.Fatalf("Failed to parse account token duration: %v", err)
	}
	return d
}
