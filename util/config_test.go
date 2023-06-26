package util

import (
	"github.com/spf13/viper"
	"github.com/stretchr/testify/require"
	"os"
	"testing"
)

// Test func LoadConfig with testing.env file
func TestLoadConfig(t *testing.T) {
	config, err := LoadConfig("../", true)
	require.NoError(t, err)
	require.NotEmpty(t, config)
}

func TestLoadConfigNotFound(t *testing.T) {
	config, err := LoadConfig("../../", true)
	require.NoError(t, err)
	require.Empty(t, config)
}

func TestLoadConfigNotFoundWithEnv(t *testing.T) {
	err := os.Setenv("DB_SOURCE", "asdfasdf")
	require.NoError(t, err)
	config, err := LoadConfig("../../", true)
	require.NoError(t, err)
	require.NotEmpty(t, config)
	require.Equal(t, "asdfasdf", config.DBSource)
}

func TestLoadConfigGetOS(t *testing.T) {
	fakeSource := "asdfasdf"
	err := os.Setenv("DB_SOURCE", fakeSource)
	require.NoError(t, err)

	config, err := LoadConfig("../", true)

	require.NoError(t, err)
	require.Equal(t, fakeSource, config.DBSource)
}

func TestBindEnv(t *testing.T) {
	config := Config{}
	_, err := BindEnv(config, viper.New())
	require.NoError(t, err)
}
