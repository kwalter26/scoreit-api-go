package util

import (
	"errors"
	"fmt"
	"github.com/rs/zerolog/log"
	"github.com/spf13/viper"
	"os"
	"reflect"
	"time"
)

type Config struct {
	Environment          string        `mapstructure:"ENVIRONMENT"`
	DBDriver             string        `mapstructure:"DB_DRIVER"`
	DBSource             string        `mapstructure:"DB_SOURCE"`
	MigrationUrl         string        `mapstructure:"MIGRATION_URL"`
	HttpServerAddress    string        `mapstructure:"HTTP_SERVER_ADDRESS"`
	AccessTokenDuration  time.Duration `mapstructure:"ACCESS_TOKEN_DURATION"`
	RefreshTokenDuration time.Duration `mapstructure:"REFRESH_TOKEN_DURATION"`

	TokenSymmetricKey string `mapstructure:"TOKEN_SYMMETRIC_KEY"`
	CasbinModelPath   string `mapstructure:"CASBIN_MODEL_PATH"`
	CasbinPolicyPath  string `mapstructure:"CASBIN_POLICY_PATH"`

	AuthEnabled   bool   `mapstructure:"AUTH_ENABLED"`
	Auth0Domain   string `mapstructure:"AUTH0_DOMAIN"`
	Auth0Audience string `mapstructure:"AUTH0_AUDIENCE"`
}

const (
	Production  string = "production"
	Staging     string = "staging"
	Local       string = "local"
	Testing     string = "testing"
	Development string = "development"
)

const ConfigPath = "."
const ConfigType = "yaml"

// LoadConfig loads the configuration from the environment variables and the config file.
func LoadConfig(path string) (config Config, err error) {

	// Load environment using viper
	configLoader := viper.New()

	// Load environment variables
	configLoader.AutomaticEnv()

	// Set the config file name and path.
	configLoader.AddConfigPath(path)
	configLoader.SetConfigType(ConfigType)

	// Set the config file name based on the environment.
	environment := os.Getenv("ENVIRONMENT")
	if environment == "" {
		log.Info().Msg("Environment not set. Defaulting to development")
		environment = "development"
	} else if !IsValidEnvironment(environment) {
		return Config{}, fmt.Errorf("invalid environment: %s", environment)
	}
	log.Info().Str("environment", environment).Msg("environment loaded")

	log.Info().Str("configPath", path).Str("configType", ConfigType).Msg("finding config")
	configName, err := findConfig(path, environment)
	if err != nil {
		return Config{}, err
	}
	configLoader.SetConfigName(configName)

	err = BindEnv(configLoader)
	err = configLoader.ReadInConfig()
	if err != nil {
		var configFileNotFoundError viper.ConfigFileNotFoundError
		if errors.As(err, &configFileNotFoundError) {
			log.Info().Msg("Config file not found. Using environment variables")
		} else {
			return Config{}, fmt.Errorf("error reading config file: %w", err)
		}
	}

	err = configLoader.Unmarshal(&config)
	if err != nil {
		return Config{}, fmt.Errorf("unable to decode into struct: %w", err)
	}
	return config, nil
}

func IsValidEnvironment(environment string) bool {
	switch environment {
	case Production, Staging, Local, Testing, Development:
		return true
	}
	return false
}

func findConfig(path, environment string) (string, error) {
	if !IsValidEnvironment(environment) {
		return "", fmt.Errorf("invalid environment: %s", environment)
	}

	files, err := os.ReadDir(path)
	if err != nil {
		return "", err
	}

	configName := "config"
	for _, file := range files {
		if file.IsDir() {
			continue
		}
		if file.Name() == fmt.Sprintf("config-%s.%s", environment, ConfigType) {
			configName = file.Name()[:len(file.Name())-len(ConfigType)-1]
			break
		}
	}

	return configName, nil
}

func BindEnv(v *viper.Viper) error {
	var config Config
	keys := reflect.ValueOf(&config).Elem()
	for i := 0; i < keys.NumField(); i++ {
		key := keys.Type().Field(i).Tag.Get("mapstructure")
		err := v.BindEnv(key)
		if err != nil {
			return err
		}
	}
	return nil
}
