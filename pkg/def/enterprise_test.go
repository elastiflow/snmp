package def

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestEnterprise_Validate(t *testing.T) {
	tests := []struct {
		name       string
		enterprise Enterprise
		wantErr    bool
		errMsg     string
	}{
		{
			name:       "empty enterprise",
			enterprise: Enterprise(""),
			wantErr:    true,
			errMsg:     "enterprise value cannot be empty",
		},
		{
			name:       "enterprise without yml extension",
			enterprise: Enterprise("not_a_yaml_file.txt"),
			wantErr:    true,
			errMsg:     "enterprise value must be a path to a YAML file (with .yml or .yaml extension): not_a_yaml_file.txt",
		},
		{
			name:       "enterprise with yml extension but file doesn't exist",
			enterprise: Enterprise("non_existent.yml"),
			wantErr:    true,
			errMsg:     "enterprise file does not exist",
		},
		{
			name:       "valid enterprise with yml extension",
			enterprise: Enterprise("valid.yml"),
			wantErr:    false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.enterprise.Validate("testdata/traps/rules")
			if tt.wantErr {
				assert.Error(t, err)
				assert.Contains(t, err.Error(), tt.errMsg)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}
