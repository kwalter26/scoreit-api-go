package security

import (
	"embed"
	"fmt"
	casbin "github.com/casbin/casbin/v2"
	"github.com/kwalter26/scoreit-api-go/util"
	"os"
)

//go:embed resources
var res embed.FS

func NewEnforcer(config util.Config) (*casbin.Enforcer, error) {

	if _, err := os.Stat(config.CasbinModelPath); err == nil {
		fmt.Printf("File exists\n")
	} else {
		fmt.Printf("File does not exist\n")
		model, err := res.ReadFile("resources/authz_model.conf")
		if err != nil {
			return nil, err
		}
		err = os.WriteFile(config.CasbinModelPath, model, 0644)
		if err != nil {
			return nil, err
		}
	}

	if _, err := os.Stat(config.CasbinPolicyPath); err == nil {
		fmt.Printf("File exists\n")
	} else {
		fmt.Printf("File does not exist\n")
		policy, err := res.ReadFile("resources/authz_policy.csv")
		if err != nil {
			return nil, err
		}
		err = os.WriteFile(config.CasbinPolicyPath, policy, 0644)
		if err != nil {
			return nil, err
		}
	}

	e, err := casbin.NewEnforcer(config.CasbinModelPath, config.CasbinPolicyPath)
	if err != nil {
		return nil, err
	}
	return e, nil
}

type Role string

var UserRole Role = "user"
var AdminRole Role = "admin"
var UserRoles = []Role{UserRole}
var AdminRoles = []Role{AdminRole}
