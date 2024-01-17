package security

import (
	"embed"
	"github.com/kwalter26/scoreit-api-go/test"
	"github.com/kwalter26/scoreit-api-go/util"
	"github.com/stretchr/testify/require"
	"os"
	"testing"
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
			config:  util.Config{CasbinModelPath: "../test/resources/test_authz_model.conf", CasbinPolicyPath: "../test/resources/test_authz_policy.csv"},
			fs:      test.SecurityResources(),
			wantErr: false,
		},
		{
			name:    "invalid model path",
			config:  util.Config{CasbinModelPath: "path/to/invalid/model", CasbinPolicyPath: "../test/test_authz_policy.csv"},
			fs:      test.SecurityResources(),
			wantErr: false,
		},
		{
			name:    "invalid policy path",
			config:  util.Config{CasbinModelPath: "../test/test_authz_model.conf", CasbinPolicyPath: "path/to/invalid/policy"},
			fs:      test.SecurityResources(),
			wantErr: false,
		},
		{
			name:    "invalid model and policy paths. uses embedded fs",
			config:  util.Config{CasbinModelPath: "path/to/invalid/model", CasbinPolicyPath: "path/to/invalid/policy"},
			fs:      test.SecurityResources(),
			wantErr: false,
		},
		{
			name:    "invalid embedded fs",
			config:  util.Config{CasbinModelPath: "path/to/invalid/model", CasbinPolicyPath: "path/to/invalid/policy"},
			fs:      &embed.FS{},
			wantErr: true,
		},
		{
			name:    "bad policy invalid csv",
			config:  util.Config{CasbinModelPath: "path/to/invalid/model", CasbinPolicyPath: "../test/bad_authz_policy.csv"},
			fs:      test.SecurityResources(),
			wantErr: true,
		},
		{
			name:    "bad model embedded fs",
			config:  util.Config{CasbinModelPath: "../test/bad_authz_model.conf", CasbinPolicyPath: "path/to/invalid/policy"},
			fs:      test.SecurityResources(),
			wantErr: true,
		},
	}

	// Create a file named test_authz_model.conf

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			fileA, errA := os.Create("test_authz_model.conf")
			if errA != nil {
				t.Fatal(errA)
			}
			defer func(fileA *os.File) {
				err := fileA.Close()
				require.NoError(t, err)
			}(fileA)

			// Create a file named test_authz_policy.csv
			fileB, errB := os.Create("test_authz_policy.csv")
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
			err = os.Remove("test_authz_model.conf")
			require.NoError(t, err)
			err = os.Remove("test_authz_policy.csv")
			require.NoError(t, err)
		})
	}

}
