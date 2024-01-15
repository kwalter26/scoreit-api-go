// Package security provides utilities for building
// authorization structures in a Go application using
// casbin.
package security

import (
	"bytes"
	"embed"
	"encoding/csv"
	casbin "github.com/casbin/casbin/v2"
	"github.com/casbin/casbin/v2/model"
	"github.com/kwalter26/scoreit-api-go/util"
	"github.com/rs/zerolog/log"
	"strings"

	"os"
)

// res is a file server that serves embedded static files.
//
//go:embed resources
var res embed.FS

// getFileContent reads the content of a file specified by filePath.
// If the file exists, it returns the file contents as a byte slice ([]byte).
// If the file does not exist, it loads the default file specified by defaultPath and returns its contents.
// The fileType parameter is used for logging purposes.
// Returns []byte and error.
func getFileContent(filePath string, defaultPath string, fileType string) ([]byte, error) {
	var fileBytes []byte
	var err error

	if _, fileErr := os.Stat(filePath); fileErr == nil {
		log.Info().Msgf("Loading %s file %s", fileType, filePath)
		fileBytes, err = os.ReadFile(filePath)
	} else {
		log.Warn().Msgf("%s File '%s' does not exist. Loading default %s file %s", fileType, filePath, fileType, defaultPath)
		fileBytes, err = res.ReadFile(defaultPath)
	}

	return fileBytes, err
}

// getModelBytes reads the content of the Casbin model file specified by config.CasbinModelPath.
// If the file exists, it returns the file contents as a byte slice ([]byte).
// If the file does not exist, it loads the default model file "resources/authz_model.conf" and returns its contents.
// The "Model" parameter is used for logging purposes.
// Requires a valid util.Config instance as a parameter.
// Returns the byte slice of the model file and an error.
func getModelBytes(config util.Config) ([]byte, error) {
	return getFileContent(config.CasbinModelPath, "resources/authz_model.conf", "Model")
}

// getPolicyBytes reads the content of the Casbin Policy file specified by config.CasbinPolicyPath.
// If the file exists, it returns the file contents as a byte slice ([]byte).
// If the file does not exist, it loads the default policy file "resources/authz_policy.csv" and returns its contents.
// The "Policy" parameter is used for logging purposes.
// Requires a valid util.Config instance as a parameter.
// Returns the byte slice of the policy file and an error.
func getPolicyBytes(config util.Config) ([]byte, error) {
	return getFileContent(config.CasbinPolicyPath, "resources/authz_policy.csv", "Policy")
}

// NewEnforcer creates a new casbin enforcer with the provided configuration.
// It loads the model bytes and policy bytes from the specified paths in the config.
// Then it initializes the casbin enforcer with the loaded model.
// Next, it parses the policy bytes as a CSV file and adds the policies to the enforcer.
// If any policy does not contain the required 4 parameters, it logs a warning.
// If any error occurs during the process, it returns nil and the error.
// Otherwise, it returns the initialized enforcer and nil error.
func NewEnforcer(config util.Config) (*casbin.Enforcer, error) {
	modelBytes, err := getModelBytes(config)
	if err != nil {
		return nil, err
	}

	policyBytes, err := getPolicyBytes(config)
	if err != nil {
		return nil, err
	}

	modelFromString, err := model.NewModelFromString(string(modelBytes))
	if err != nil {
		return nil, err
	}

	csvReader := csv.NewReader(bytes.NewReader(policyBytes))
	policies, err := csvReader.ReadAll()
	if err != nil {
		return nil, err
	}

	e, err := casbin.NewEnforcer(modelFromString)
	if err != nil {
		return nil, err
	}

	for _, policy := range policies {
		if len(policy) == 4 {
			if ok, err := e.AddNamedPolicy(
				strings.TrimSpace(policy[0]),
				strings.TrimSpace(policy[1]),
				strings.TrimSpace(policy[2]),
				strings.TrimSpace(policy[3]),
			); ok {
				log.Info().Msgf("Added policy: %s", policy)
			} else {
				log.Warn().Msgf("Failed to add policy: %s with error %s", policy, err)
				return nil, err
			}
		} else {
			log.Warn().Msgf("Policy '%s' does not contain required 4 params", policy)
		}
	}
	return e, nil
}

// Role type tracks the role of a user
type Role string

// UserRole is the role for normal users
var UserRole Role = "user"

// AdminRole is the role for admin users
var AdminRole Role = "admin"

// UserRoles includes all roles a user can have
var UserRoles = []Role{UserRole}

// AdminRoles includes all roles an admin can have
var AdminRoles = []Role{AdminRole}
