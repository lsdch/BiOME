package config

import (
	"encoding/json"
	"fmt"
	"net/url"
	"path/filepath"
	"strings"
	"time"

	"github.com/danielgtaylor/huma/v2"
	"github.com/go-playground/validator/v10"
	"github.com/lsdch/biome/models"
	"github.com/lsdch/biome/services/crossref"
	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
)

var validate = validator.
	New(validator.WithRequiredStructEnabled(), validator.WithTagNameFuncBlankOmit())

type Validable interface {
	Valid() bool
}

func validateEnum(fl validator.FieldLevel) bool {
	validable, ok := fl.Field().Interface().(Validable)
	if ok {
		return validable.Valid()
	}
	return true
}

type AppEnv string

const (
	EnvDev  AppEnv = "dev"
	EnvProd AppEnv = "prod"
)

func (e AppEnv) Valid() bool {
	switch e {
	case EnvDev, EnvProd:
		return true
	default:
		return false
	}
}

type APIConfig struct {
	BasePath     string `mapstructure:"API_BASE_PATH" validate:"required"`
	Host         string `mapstructure:"API_HOST" validate:"required"`
	Port         string `mapstructure:"API_PORT" validate:"required"`
	Version      string `mapstructure:"VERSION" validate:"required,semver"`
	Title        string `mapstructure:"API_TITLE" validate:"required,min=3"`
	Description  string `mapstructure:"API_DESCRIPTION"`
	ContactName  string `mapstructure:"API_CONTACT_NAME" validate:"required"`
	ContactEmail string `mapstructure:"API_CONTACT_EMAIL" validate:"required,email"`
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
	DSN             string        `mapstructure:"DSN" validate:"required"`
	MaxConns        int32         `mapstructure:"MAX_CONNS" validate:"required,gt=0,gtefield=MinConns"`
	MinConns        int32         `mapstructure:"MIN_CONNS" validate:"gt=0"`
	MaxConnLifetime time.Duration `mapstructure:"MAX_CONN_LIFETIME"`
	MaxConnIdleTime time.Duration `mapstructure:"MAX_CONN_IDLE_TIME"`
	HealthTimeout   time.Duration `mapstructure:"HEALTH_TIMEOUT" validate:"gt=0"`
}

type UserBootstrap struct {
	Login     string          `mapstructure:"login" validate:"required"`
	Email     string          `mapstructure:"email" validate:"required,email"`
	Password  string          `mapstructure:"password" validate:"required,min=6"`
	Role      models.UserRole `mapstructure:"role" validate:"required,enum"`
	FirstName string          `mapstructure:"first_name" validate:"required,min=2"`
	LastName  string          `mapstructure:"last_name" validate:"required,min=2"`
}

type CountriesBootstrap struct {
	CountriesJSON_URL      string `mapstructure:"COUNTRIES_JSON_URL" validate:"required,url"`
	CountryJSON_CachePath  string `mapstructure:"COUNTRIES_JSON_CACHE_PATH" validate:"required"`
	CountryNameKey         string `mapstructure:"COUNTRY_NAME_KEY" validate:"required"`
	CountryCodeKey         string `mapstructure:"COUNTRY_CODE_KEY" validate:"required"`
	CountryContinentKey    string `mapstructure:"COUNTRY_CONTINENT_KEY" validate:"required"`
	CountrySubcontinentKey string `mapstructure:"COUNTRY_SUBCONTINENT_KEY" validate:"required"`
}
type BootstrapConfig struct {
	Countries CountriesBootstrap `mapstructure:"countries"`
	Users     []UserBootstrap    `mapstructure:"users"`
}

type AuthTokensConfig struct {
	AuthTokenCookieName    string        `mapstructure:"AUTH_TOKEN_COOKIE_NAME" validate:"required,min=3"`
	AuthTokenLifetime      time.Duration `mapstructure:"AUTH_TOKEN_LIFETIME" validate:"gt=0"`
	SecretKey              string        `mapstructure:"JWT_SECRET_KEY" validate:"required,min=32"`
	RefreshTokenCookieName string        `mapstructure:"REFRESH_TOKEN_COOKIE_NAME" validate:"required,min=3"`
	RefreshTokenLifetime   time.Duration `mapstructure:"REFRESH_TOKEN_LIFETIME" validate:"gt=0"`
	RefreshTokenPepper     string        `mapstructure:"REFRESH_TOKEN_PEPPER" validate:"required,min=8"`
}

type SMTPConfig struct {
	SMTPHost     string `mapstructure:"SMTP_HOST" validate:"required,hostname"`
	SMTPPort     uint   `mapstructure:"SMTP_PORT" validate:"required,port"`
	SMTPUser     string `mapstructure:"SMTP_USER" validate:"required,min=3"`
	SMTPPassword string `mapstructure:"SMTP_PASSWORD" validate:"required,min=3"`
}

type InstanceConfig struct {
	AppName                string  `mapstructure:"APP_NAME" validate:"required,min=3"`
	AppSubtitle            *string `mapstructure:"APP_SUBTITLE" validate:"omitempty,min=3"`
	AppDescription         *string `mapstructure:"APP_DESCRIPTION" validate:"omitempty"`
	AdminEmail             string  `mapstructure:"ADMIN_EMAIL" validate:"required,email"`
	IsPublic               bool    `mapstructure:"IS_PUBLIC"`
	AccountRequestsEnabled bool    `mapstructure:"ACCOUNT_REQUESTS_ENABLED"`
	MailFromAddress        string  `mapstructure:"MAIL_FROM_ADDRESS" validate:"required,email"`
	MailFromName           string  `mapstructure:"MAIL_FROM_NAME" validate:"required,min=3"`
	MolecularDataEnabled   bool    `mapstructure:"MOLECULAR_DATA_ENABLED"`
}

type GeoapifyConfig struct {
	GeoApifyApiKey  string `mapstructure:"GEOAPIFY_API_KEY"`
	DailyUsageLimit int32  `mapstructure:"GEOAPIFY_DAILY_USAGE_LIMIT" validate:"min=0"`
}

type GBIFConfig struct {
	UserAgent          string `mapstructure:"GBIF_USER_AGENT" validate:"required"`
	BackboneDatasetKey string `mapstructure:"GBIF_BACKBONE_DATASET_KEY" validate:"required,uuid"`
	MaxConcurrent      int    `mapstructure:"GBIF_MAX_CONCURRENT" validate:"min=1"`
}

type Config struct {
	Instance             InstanceConfig   `mapstructure:"instance" validate:"required"`
	appPublicBaseURL     string           `mapstructure:"APP_PUBLIC_BASE_URL" validate:"required,url"`
	AppPublicBaseURL     url.URL          `json:"-"`
	DB                   DBConfig         `mapstructure:"DB" validate:"required"`
	RawFileStorageRoot   string           `mapstructure:"RAW_FILE_STORAGE_ROOT" validate:"required"`
	SMTP                 SMTPConfig       `mapstructure:"SMTP" validate:"required"`
	Env                  AppEnv           `mapstructure:"ENV" validate:"required,enum"`
	API                  APIConfig        `mapstructure:"API" validate:"required"`
	AuthTokens           AuthTokensConfig `mapstructure:"auth_tokens" validate:"required"`
	GeneratedTokenLength uint             `mapstructure:"TOKEN_LENGTH" validate:"min=16"`
	Geoapify             GeoapifyConfig   `mapstructure:"geoapify" validate:"required"`
	GBIF                 GBIFConfig       `mapstructure:"gbif" validate:"required"`
	Bootstrap            BootstrapConfig  `mapstructure:"bootstrap" validate:"required"`
	CrossRef             crossref.Config  `mapstructure:"crossref" validate:"required"`
}

func (c *Config) Validate() error {
	if err := validate.RegisterValidation("enum", validateEnum); err != nil {
		return err
	}
	return validate.Struct(c)
}

func LoadConfig(dir string, name string) (Config, error) {
	v := viper.New()
	v.AddConfigPath(dir)
	v.SetConfigName(name)
	v.SetConfigType("yaml")

	v.SetEnvPrefix("BIOME")
	v.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))

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

	if cfg.Env == EnvDev {
		cfgJSON, _ := json.MarshalIndent(cfg, "", "\t")
		logrus.Infof("Loaded config :\n%s", cfgJSON)
	}

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
