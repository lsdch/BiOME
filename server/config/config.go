package config

import (
	"fmt"
	"net/url"
	"path/filepath"
	"time"

	"github.com/danielgtaylor/huma/v2"
	"github.com/spf13/viper"
)

type AppEnv string

const (
	EnvDev  AppEnv = "dev"
	EnvProd AppEnv = "prod"
)

type APIConfig struct {
	BasePath     string `mapstructure:"API_BASE_PATH"`
	Host         string `mapstructure:"API_HOST"`
	Port         string `mapstructure:"API_PORT"`
	Version      string `mapstructure:"VERSION"`
	Title        string `mapstructure:"API_TITLE"`
	Description  string `mapstructure:"API_DESCRIPTION"`
	ContactName  string `mapstructure:"API_CONTACT_NAME"`
	ContactEmail string `mapstructure:"API_CONTACT_EMAIL"`
}

func (c APIConfig) ToHumaConfig() huma.Config {
	var cfg = huma.DefaultConfig(c.Title, c.Version)
	cfg.Info.Description = fmt.Sprintf("%s: %s", c.Title, c.Description)
	cfg.Components.SecuritySchemes = map[string]*huma.SecurityScheme{
		"bearer": {
			Type:         "http",
			Scheme:       "bearer",
			BearerFormat: "JWT",
		},
		"cookieAuth": {
			Type: "apiKey",
			In:   "cookie",
			Name: "auth_token",
		},
	}
	cfg.Info.Contact = &huma.Contact{
		Name:  c.ContactName,
		Email: c.ContactEmail,
	}
	cfg.OpenAPI.Servers = []*huma.Server{
		{URL: c.BasePath},
	}
	cfg.Security = []map[string][]string{
		{"bearer": {}},
		{"cookieAuth": {}},
	}

	return cfg
}

type DBConfig struct {
	DSN             string        `mapstructure:"DSN"`
	MaxConns        int32         `mapstructure:"MAX_CONNS"`
	MinConns        int32         `mapstructure:"MIN_CONNS"`
	MaxConnLifetime time.Duration `mapstructure:"MAX_CONN_LIFETIME"`
	MaxConnIdleTime time.Duration `mapstructure:"MAX_CONN_IDLE_TIME"`
	HealthTimeout   time.Duration `mapstructure:"HEALTH_TIMEOUT"`
}

type BootstrapConfig struct {
	CountriesJSON_URL      string `mapstructure:"COUNTRIES_JSON_URL"`
	CountryNameKey         string `mapstructure:"COUNTRY_NAME_KEY"`
	CountryCodeKey         string `mapstructure:"COUNTRY_CODE_KEY"`
	CountryContinentKey    string `mapstructure:"COUNTRY_CONTINENT_KEY"`
	CountrySubcontinentKey string `mapstructure:"COUNTRY_SUBCONTINENT_KEY"`
}

type AuthTokensConfig struct {
	AuthTokenCookieName    string        `mapstructure:"AUTH_TOKEN_COOKIE_NAME"`
	AuthTokenLifetime      time.Duration `mapstructure:"AUTH_TOKEN_LIFETIME"`
	SecretKey              string        `mapstructure:"JWT_SECRET_KEY"`
	RefreshTokenCookieName string        `mapstructure:"REFRESH_TOKEN_COOKIE_NAME"`
	RefreshTokenLifetime   time.Duration `mapstructure:"REFRESH_TOKEN_LIFETIME"`
	RefreshTokenPepper     string        `mapstructure:"REFRESH_TOKEN_PEPPER"`
}

type SMTPConfig struct {
	SMTPHost     string `mapstructure:"SMTP_HOST"`
	SMTPPort     int    `mapstructure:"SMTP_PORT"`
	SMTPUser     string `mapstructure:"SMTP_USER"`
	SMTPPassword string `mapstructure:"SMTP_PASSWORD"`
}

type InstanceConfig struct {
	AppName        string  `mapstructure:"APP_NAME"`
	AppSubtitle    *string `mapstructure:"APP_SUBTITLE"`
	AppDescription *string `mapstructure:"APP_DESCRIPTION"`
	AdminEmail     string  `mapstructure:"ADMIN_EMAIL"`
	IsPublic       bool    `mapstructure:"IS_PUBLIC"`

	AccountRequestsEnabled bool   `mapstructure:"ACCOUNT_REQUESTS_ENABLED"`
	MailFromAddress        string `mapstructure:"MAIL_FROM_ADDRESS"`
	MailFromName           string `mapstructure:"MAIL_FROM_NAME"`
}

type GeoapifyConfig struct {
	GeoApifyApiKey  string `mapstructure:"GEOAPIFY_API_KEY"`
	DailyUsageLimit int32  `mapstructure:"GEOAPIFY_DAILY_USAGE_LIMIT"`
}

type Config struct {
	Instance               InstanceConfig   `mapstructure:"instance"`
	AppPublicBaseURL       url.URL          `mapstructure:"APP_PUBLIC_BASE_URL"`
	DB                     DBConfig         `mapstructure:"DB"`
	SMTP                   SMTPConfig       `mapstructure:"SMTP"`
	Env                    AppEnv           `mapstructure:"ENV"`
	API                    APIConfig        `mapstructure:"API"`
	AuthTokens             AuthTokensConfig `mapstructure:"auth_tokens"`
	GeneratedTokenLength   uint             `mapstructure:"TOKEN_LENGTH"`
	Geoapify               GeoapifyConfig   `mapstructure:"geoapify"`
	GBIFBackboneDatasetKey string           `mapstructure:"GBIF_BACKBONE_DATASET_KEY"`
	GBIFMaxConcurrent      int              `mapstructure:"GBIF_MAX_CONCURRENT"`
	CrossRefMaxConcurrent  int              `mapstructure:"CROSSREF_MAX_CONCURRENT"`
	CountriesJSON_URL      string           `mapstructure:"COUNTRIES_JSON_URL"`
	Bootstrap              BootstrapConfig  `mapstructure:"bootstrap"`
}

func (c Config) Validate() error {
	if c.AuthTokens.SecretKey == "" {
		return fmt.Errorf("missing JWT_SECRET_KEY")
	}
	return nil
}

func LoadConfig(dir string, name string) (Config, error) {
	v := viper.New()
	v.AddConfigPath(dir)
	v.SetConfigName(name)
	v.SetConfigType("yaml")

	v.AutomaticEnv()

	// 1. defaults (fallback)
	setDefaults(v)

	// 2. base config (required)
	if err := readFile(v, name); err != nil {
		return Config{}, fmt.Errorf("base config: %w", err)
	}

	// 3. env config (dev/prod)
	env := v.GetString("ENV")
	if err := readFile(v, fmt.Sprintf("%s.%s", name, env)); err != nil {
		return Config{}, fmt.Errorf("env config: %w", err)
	}

	// 4. local override (optional)
	_ = tryMerge(v, filepath.Join(dir, fmt.Sprintf("%s.%s.local.yaml", name, env)))

	var cfg Config
	if err := v.Unmarshal(&cfg); err != nil {
		return Config{}, err
	}

	return cfg, cfg.Validate()
}

func setDefaults(v *viper.Viper) {
	v.SetDefault("JWT_LIFETIME_MINUTES", 30)
	v.SetDefault("ACCOUNT_TOKEN_LIFETIME_HOURS", 24)
	v.SetDefault("TOKEN_LENGTH", 32)
}

func readFile(v *viper.Viper, name string) error {
	v.SetConfigName(name)
	return v.ReadInConfig()
}

func tryMerge(v *viper.Viper, file string) error {
	v.SetConfigFile(file)

	err := v.MergeInConfig()
	if err == nil {
		return nil
	}

	if _, ok := err.(viper.ConfigFileNotFoundError); ok {
		return nil
	}

	return err
}
