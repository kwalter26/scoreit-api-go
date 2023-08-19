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

//go:embed resources
var res embed.FS

func NewEnforcer(config util.Config) (*casbin.Enforcer, error) {

	var modelBytes []byte
	var policyBytes []byte

	if _, err := os.Stat(config.CasbinModelPath); err == nil {
		log.Info().Msgf("Loading file %s", config.CasbinModelPath)
		modelBytes, err = os.ReadFile(config.CasbinModelPath)
	} else {
		log.Warn().Msgf("File '%s' does not exist", config.CasbinModelPath)
		log.Info().Msgf("Loading default file %s", "resources/authz_model.conf")
		modelBytes, err = res.ReadFile("resources/authz_model.conf")
		if err != nil {
			return nil, err
		}
	}

	if _, err := os.Stat(config.CasbinPolicyPath); err == nil {
		log.Info().Msgf("Loading file %s", config.CasbinPolicyPath)
		policyBytes, err = os.ReadFile(config.CasbinPolicyPath)
	} else {
		log.Warn().Msgf("File '%s' does not exist", config.CasbinPolicyPath)
		log.Info().Msgf("Loading default file %s", "resources/authz_policy.csv")
		policyBytes, err = res.ReadFile("resources/authz_policy.csv")
		if err != nil {
			return nil, err
		}
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
				log.Warn().Msgf("Failed to add policy: %s", policy)
				return nil, err
			}
		} else {
			log.Warn().Msgf("Failed to add policy: %s", policy)
		}
	}

	return e, nil
}

type Role string

var UserRole Role = "user"
var AdminRole Role = "admin"
var UserRoles = []Role{UserRole}
var AdminRoles = []Role{AdminRole}
