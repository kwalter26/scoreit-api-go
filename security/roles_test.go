package security

import (
	"embed"
	"github.com/kwalter26/scoreit-api-go/test"
	"github.com/kwalter26/scoreit-api-go/util"
	"github.com/stretchr/testify/require"
	"os"
	"testing"
)

const (
	validModelPath    = "../test/resources/test_authz_model.conf"
	validPolicyPath   = "../test/resources/test_authz_policy.csv"
	invalidModelPath  = "path/to/invalid/model"
	invalidPolicyPath = "path/to/invalid/policy"
	badPolicyPath     = "../test/bad_authz_policy.csv"
	badModelPath      = "../test/bad_authz_model.conf"
	testModelConfFile = "test_authz_model.conf"
	testPolicyCsvFile = "test_authz_policy.csv"
)

func TestNewEnforcer(t *testing.T) {
	tests := []struct {
		name    string
		config  util.Config
		fs      *embed.FS
		wantErr bool
	}{
		{
			name:    "valid case",
			config:  util.Config{CasbinModelPath: validModelPath, CasbinPolicyPath: validPolicyPath},
			fs:      test.SecurityResources(),
			wantErr: false,
		},
		{
			name:    "invalid model path",
			config:  util.Config{CasbinModelPath: invalidModelPath, CasbinPolicyPath: validPolicyPath},
			fs:      test.SecurityResources(),
			wantErr: false,
		},
		{
			name:    "invalid policy path",
			config:  util.Config{CasbinModelPath: validModelPath, CasbinPolicyPath: invalidPolicyPath},
			fs:      test.SecurityResources(),
			wantErr: false,
		},
		{
			name:    "invalid model and policy paths. uses embedded fs",
			config:  util.Config{CasbinModelPath: invalidModelPath, CasbinPolicyPath: invalidPolicyPath},
			fs:      test.SecurityResources(),
			wantErr: false,
		},
		{
			name:    "invalid embedded fs",
			config:  util.Config{CasbinModelPath: invalidModelPath, CasbinPolicyPath: invalidPolicyPath},
			fs:      &embed.FS{},
			wantErr: true,
		},
		{
			name:    "bad policy invalid csv",
			config:  util.Config{CasbinModelPath: invalidModelPath, CasbinPolicyPath: badPolicyPath},
			fs:      test.SecurityResources(),
			wantErr: true,
		},
		{
			name:    "bad model embedded fs",
			config:  util.Config{CasbinModelPath: badModelPath, CasbinPolicyPath: invalidPolicyPath},
			fs:      test.SecurityResources(),
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Create a file named test_authz_model.conf
			fileA, errA := os.Create(testModelConfFile)
			if errA != nil {
				t.Fatal(errA)
			}
			defer func(fileA *os.File) {
				err := fileA.Close()
				require.NoError(t, err)
			}(fileA)

			// Create a file named test_authz_policy.csv
			fileB, errB := os.Create(testPolicyCsvFile)
			if errB != nil {
				t.Fatal(errB)
			}
			defer func(fileB *os.File) {
				err := fileB.Close()
				require.NoError(t, err)
			}(fileB)

			_, err := NewEnforcer(tt.config, tt.fs)
			if (err != nil) != tt.wantErr {
				t.Errorf("NewEnforcer() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			err = os.Remove(testModelConfFile)
			require.NoError(t, err)
			err = os.Remove(testPolicyCsvFile)
			require.NoError(t, err)
		})
	}

}
