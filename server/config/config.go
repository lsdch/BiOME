package config

import (
	"encoding/json"
	"fmt"
	"net/url"
	"path/filepath"
	"time"

	"github.com/danielgtaylor/huma/v2"
	"github.com/lsdch/biome/models"
	"github.com/sirupsen/logrus"
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

	serverURL := url.URL{
		Host: fmt.Sprintf("%s:%d", c.Host, 5173),
		Path: c.BasePath,
	}
	cfg.OpenAPI.Servers = []*huma.Server{
		{URL: serverURL.String()},
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

type UserBootstrap struct {
	Login     string          `mapstructure:"login"`
	Email     string          `mapstructure:"email"`
	Password  string          `mapstructure:"password"`
	Role      models.UserRole `mapstructure:"role"`
	FirstName string          `mapstructure:"first_name"`
	LastName  string          `mapstructure:"last_name"`
}

type CountriesBootstrap struct {
	CountriesJSON_URL      string `mapstructure:"COUNTRIES_JSON_URL"`
	CountryNameKey         string `mapstructure:"COUNTRY_NAME_KEY"`
	CountryCodeKey         string `mapstructure:"COUNTRY_CODE_KEY"`
	CountryContinentKey    string `mapstructure:"COUNTRY_CONTINENT_KEY"`
	CountrySubcontinentKey string `mapstructure:"COUNTRY_SUBCONTINENT_KEY"`
}
type BootstrapConfig struct {
	Countries CountriesBootstrap `mapstructure:"countries"`
	Users     []UserBootstrap    `mapstructure:"users"`
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
	AppName                string  `mapstructure:"APP_NAME"`
	AppSubtitle            *string `mapstructure:"APP_SUBTITLE"`
	AppDescription         *string `mapstructure:"APP_DESCRIPTION"`
	AdminEmail             string  `mapstructure:"ADMIN_EMAIL"`
	IsPublic               bool    `mapstructure:"IS_PUBLIC"`
	AccountRequestsEnabled bool    `mapstructure:"ACCOUNT_REQUESTS_ENABLED"`
	MailFromAddress        string  `mapstructure:"MAIL_FROM_ADDRESS"`
	MailFromName           string  `mapstructure:"MAIL_FROM_NAME"`
	MolecularDataEnabled   bool    `mapstructure:"MOLECULAR_DATA_ENABLED"`
}

type GeoapifyConfig struct {
	GeoApifyApiKey  string `mapstructure:"GEOAPIFY_API_KEY"`
	DailyUsageLimit int32  `mapstructure:"GEOAPIFY_DAILY_USAGE_LIMIT"`
}

type Config struct {
	Instance               InstanceConfig   `mapstructure:"instance"`
	appPublicBaseURL       string           `mapstructure:"APP_PUBLIC_BASE_URL"`
	AppPublicBaseURL       url.URL          `json:"-"`
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

func (c *Config) Validate() error {
	if c.AuthTokens.SecretKey == "" {
		return fmt.Errorf("missing JWT_SECRET_KEY")
	}
	if u, err := url.Parse(c.appPublicBaseURL); err != nil {
		return fmt.Errorf("invalid APP_PUBLIC_BASE_URL: %w", err)
	} else {
		c.AppPublicBaseURL = *u
	}
	if c.GeneratedTokenLength < 16 {
		return fmt.Errorf("TOKEN_LENGTH must be at least 16")
	}
	if c.Geoapify.DailyUsageLimit <= 0 {
		return fmt.Errorf("GEOAPIFY_DAILY_USAGE_LIMIT must be greater than 0")
	}
	if c.GBIFMaxConcurrent <= 0 {
		return fmt.Errorf("GBIF_MAX_CONCURRENT must be greater than 0")
	}
	if c.CrossRefMaxConcurrent <= 0 {
		return fmt.Errorf("CROSSREF_MAX_CONCURRENT must be greater than 0")
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
	if env == "" {
		env = string(EnvDev)
	}
	if err := tryMerge(v, fmt.Sprintf("%s.%s.yml", name, env)); err != nil {
		return Config{}, fmt.Errorf("env config: %w", err)
	}

	// 4. local override (optional)
	_ = tryMerge(v, filepath.Join(dir, fmt.Sprintf("%s.local.yml", name)))

	// fmt.Printf("%#v\n", v.AllSettings())

	var cfg Config
	if err := v.Unmarshal(&cfg); err != nil {
		return Config{}, err
	}

	cfgJSON, _ := json.MarshalIndent(cfg, "", "\t")
	logrus.Infof("Loaded config :\n%s", cfgJSON)

	return cfg, cfg.Validate()
}

func setDefaults(v *viper.Viper) {
	v.SetDefault("JWT_LIFETIME_MINUTES", 30)
	v.SetDefault("ACCOUNT_TOKEN_LIFETIME_HOURS", 24)
	v.SetDefault("TOKEN_LENGTH", 32)
}

func readFile(v *viper.Viper, name string) error {
	logrus.Infof("Loading config file: %s", name)
	v.SetConfigName(name)
	return v.ReadInConfig()
}

func tryMerge(v *viper.Viper, file string) error {
	logrus.Infof("Attempting to merge config file: %s", file)
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
