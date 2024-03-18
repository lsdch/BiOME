package config

// type Config struct {
// 	Client ClientConfig `mapstructure:",squash"`
// }

// type ClientConfig struct {
// 	AppName    string `mapstructure:"VITE_APP_NAME"`
// 	DomainName string `mapstructure:"VITE_DOMAIN" validate:"required"`
// 	Port       int    `mapstructure:"VITE_PORT" validate:"required"`
// }

// var (
// 	once sync.Once

// 	// App configuration as a struct with some default values
// 	config = Config{}
// )

// Generates a URL to an endpoint on the server.
//
// URL path must omit the API prefix, e.g. '/api/v1/users' must be '/users/'
// func (config *Config) MakeURL(url_path string) url.URL {
// 	return url.URL{
// 		Scheme: "https",
// 		Host:   fmt.Sprintf("%s:%d", config.DomainName, config.Port),
// 		Path:   path.Join(config.BasePath, url_path),
// 	}
// }

// // Generates a URL to a client page.
// func (config *Config) MakeClientURL(url_path string) url.URL {
// 	return url.URL{
// 		Scheme: "https",
// 		Host:   fmt.Sprintf("%s:%d", config.Client.DomainName, config.Client.Port),
// 		Path:   url_path,
// 	}
// }

// func LoadConfig(path string) (*Config, error) {
// 	var err error = nil
// 	once.Do(func() {
// 		viper.AddConfigPath(path)
// 		viper.SetConfigType("env")

// 		var clientConfigFile string
// 		switch gin.Mode() {
// 		case gin.ReleaseMode:
// 			viper.SetConfigName("app")
// 			clientConfigFile = filepath.Join("../client", ".env.production")
// 		case gin.DebugMode, gin.TestMode:
// 			viper.SetConfigName("dev")
// 			clientConfigFile = filepath.Join("../client", ".env.development")
// 		}

// 		viper.AutomaticEnv()

// 		if err = viper.ReadInConfig(); err != nil {
// 			return
// 		}

// 		logrus.Debugf("Loading client config file %s", clientConfigFile)
// 		viper.SetConfigName(clientConfigFile)
// 		viper.AutomaticEnv()
// 		err = viper.MergeInConfig()

// 		err = viper.Unmarshal(&config)
// 		validate := validator.New()
// 		if err = validate.Struct(&config); err != nil {
// 			log.Fatalf("Invalid configuration: %v", err)
// 		}
// 	})
// 	return &config, err
// }

// func Get() *Config {
// 	config, err := LoadConfig(".")
// 	if err != nil {
// 		log.Fatalf("Failed to load configuration file : %v", err)
// 	}
// 	return config
// }
