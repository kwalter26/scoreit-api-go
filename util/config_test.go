package util

import (
	"github.com/spf13/viper"
	"github.com/stretchr/testify/require"
	"os"
	"testing"
)

// Test func LoadConfig with testing.env file
func TestLoadConfig(t *testing.T) {
	config, err := LoadConfig("testdata")
	require.NoError(t, err)
	require.NotEmpty(t, config)
	require.Equal(t, "postgres", config.DBDriver)
}

func TestLoadConfigNotFound(t *testing.T) {
	config, err := LoadConfig(".")
	require.NoError(t, err)
	require.Empty(t, config)
}

func TestLoadConfigInvalidPath(t *testing.T) {
	config, err := LoadConfig("xyz")
	require.Error(t, err)
	require.Empty(t, config)
}

func TestLoadConfigNotFoundWithEnv(t *testing.T) {
	err := os.Setenv("DB_SOURCE", "asdfasdf")
	require.NoError(t, err)
	config, err := LoadConfig("testdata")
	require.NoError(t, err)
	require.NotEmpty(t, config)
	require.Equal(t, "asdfasdf", config.DBSource)
}

func TestLoadConfigWithEnv(t *testing.T) {
	err := os.Setenv("ENVIRONMENT", "testing")
	require.NoError(t, err)
	config, err := LoadConfig("testdata")
	require.NoError(t, err)
	require.NotEmpty(t, config)
	require.Equal(t, "testing", config.Environment)
	require.Equal(t, "postgres-TEST", config.DBDriver)
}

func TestBindEnv(t *testing.T) {
	err := BindEnv(viper.New())
	require.NoError(t, err)
}

func TestLoadConfigInvalidEnv(t *testing.T) {
	err := os.Setenv("ENVIRONMENT", "invalid")
	require.NoError(t, err)
	_, err = LoadConfig("testdata")
	require.Error(t, err)
}

func TestIsValidEnvironment(t *testing.T) {
	require.True(t, IsValidEnvironment("production"))
	require.True(t, IsValidEnvironment("staging"))
	require.True(t, IsValidEnvironment("local"))
	require.True(t, IsValidEnvironment("testing"))
	require.True(t, IsValidEnvironment("development"))
	require.False(t, IsValidEnvironment("invalid"))
}

func TestFindConfig(t *testing.T) {
	_, err := findConfig("testdata", "testing")
	require.NoError(t, err)
	_, err = findConfig("testdata", "invalid")
	require.Error(t, err)
	_, err = findConfig("invalidpath", "testing")
	require.Error(t, err)
}

func TestLoadConfigInvalidData(t *testing.T) {
	err := os.Setenv("ENVIRONMENT", "development")
	require.NoError(t, err)
	// Attempt to load the configuration
	_, err = LoadConfig("testdata/bad_data")
	require.Error(t, err)
	require.Contains(t, err.Error(), "unable to decode into struct")
}

func TestLoadConfigInvalidFile(t *testing.T) {
	err := os.Setenv("ENVIRONMENT", "testing")
	require.NoError(t, err)
	// Attempt to load the configuration
	_, err = LoadConfig("testdata/bad_data")
	require.Error(t, err)
	require.Contains(t, err.Error(), "error reading config file")
}
