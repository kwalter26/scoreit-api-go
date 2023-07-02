package util

import (
	"fmt"
	"github.com/spf13/viper"
	"path/filepath"
	"reflect"
	"time"
)

type Config struct {
	Environment         Environment   `mapstructure:"ENVIRONMENT"`
	DBDriver            string        `mapstructure:"DB_DRIVER"`
	DBSource            string        `mapstructure:"DB_SOURCE"`
	MigrationUrl        string        `mapstructure:"MIGRATION_URL"`
	HttpServerAddress   string        `mapstructure:"HTTP_SERVER_ADDRESS"`
	AccessTokenDuration time.Duration `mapstructure:"ACCESS_TOKEN_DURATION"`
	TokenSymmetricKey   string        `mapstructure:"TOKEN_SYMMETRIC_KEY"`
}

type Environment string

const (
	Testing     Environment = "testing"
	Development Environment = "development"
)

func LoadConfig(path string, isTesting bool) (config Config, err error) {
	v := viper.New()
	c, err2 := BindEnv(config, v)
	if err2 != nil {
		return c, err2
	}
	v.AutomaticEnv()
	absPath, err := filepath.Abs(path)
	filename := "app.env"
	if isTesting {
		filename = "testing.env"
	}

	v.AddConfigPath(absPath)
	v.SetConfigName(filename)
	v.SetConfigType("env")

	err = v.ReadInConfig()
	if err != nil {
		if _, ok := err.(viper.ConfigFileNotFoundError); ok {
			fmt.Printf("Config file not found: %s/%s\n", absPath, filename)
		} else {
			fmt.Printf("Config file was found but another error was produced: %s/%s\n", absPath, filename)
			return Config{}, err
		}
	}

	err = v.Unmarshal(&config)
	return config, err
}

func BindEnv(config Config, v *viper.Viper) (Config, error) {
	keys := reflect.ValueOf(&config).Elem()
	for i := 0; i < keys.NumField(); i++ {
		key := keys.Type().Field(i).Tag.Get("mapstructure")
		err := v.BindEnv(key)
		if err != nil {
			return config, err
		}
	}
	return Config{}, nil
}
