package def

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"gopkg.in/yaml.v3"
)

func TestObjectTypeDefinition_Validate(t *testing.T) {
	tests := []struct {
		name    string
		def     ObjectTypeDefinition
		wantErr bool
	}{
		{
			name: "ok",
			def: ObjectTypeDefinition{
				DefaultPollInterval: 60,
				ObjectTypes: []ObjectTypeConfig{
					{Name: "interfaces", Description: "Interface-related objects", DefaultPollInterval: 60},
					{Name: "routes", DefaultPollInterval: 300},
				},
			},
			wantErr: false,
		},
		{
			name: "missing name",
			def: ObjectTypeDefinition{
				ObjectTypes: []ObjectTypeConfig{
					{Name: "", DefaultPollInterval: 60},
				},
			},
			wantErr: true,
		},
		{
			name: "empty list",
			def: ObjectTypeDefinition{
				ObjectTypes: []ObjectTypeConfig{},
			},
			wantErr: true,
		},
		{
			name: "non-positive interval",
			def: ObjectTypeDefinition{
				ObjectTypes: []ObjectTypeConfig{
					{Name: "bad", DefaultPollInterval: 0},
				},
			},
			wantErr: true,
		},
	}

	for _, tc := range tests {
		tc := tc // capture range var
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel() // safe: each case uses its own value

			err := tc.def.Validate()

			if tc.wantErr && err == nil {
				t.Fatalf("expected error, got nil")
			}
			if !tc.wantErr && err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
		})
	}
}

func TestObjectTypeDefinition_GetValidObjectTypes(t *testing.T) {
	tests := []struct {
		name     string
		otd      ObjectTypeDefinition
		expected []string
	}{
		{
			name: "empty object types",
			otd: ObjectTypeDefinition{
				ObjectTypes: []ObjectTypeConfig{},
			},
			expected: []string{},
		},
		{
			name: "multiple object types",
			otd: ObjectTypeDefinition{
				ObjectTypes: []ObjectTypeConfig{
					{Name: "configuration"},
					{Name: "telemetry"},
					{Name: "status"},
					{Name: "performance"},
				},
			},
			expected: []string{"configuration", "telemetry", "status", "performance"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := tt.otd.ValidObjectTypes()
			assert.Equal(t, tt.expected, result)
		})
	}
}

func TestObjectTypeDefinition_IsValidObjectType(t *testing.T) {
	otd := ObjectTypeDefinition{
		ObjectTypes: []ObjectTypeConfig{
			{Name: "configuration"},
			{Name: "telemetry"},
			{Name: "status"},
			{Name: "performance"},
		},
	}

	tests := []struct {
		name       string
		objectType string
		expected   bool
	}{
		{
			name:       "valid object type",
			objectType: "configuration",
			expected:   true,
		},
		{
			name:       "valid object type - telemetry",
			objectType: "telemetry",
			expected:   true,
		},
		{
			name:       "invalid object type",
			objectType: "invalid_type",
			expected:   false,
		},
		{
			name:       "empty object type",
			objectType: "",
			expected:   true, // Empty type is allowed
		},
		{
			name:       "case sensitive",
			objectType: "Configuration",
			expected:   false, // Case sensitive
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := otd.IsValidObjectType(tt.objectType)
			assert.Equal(t, tt.expected, result)
		})
	}
}

func TestObjectTypeDefinition_GetDefaultPollIntervalForType(t *testing.T) {
	otd := ObjectTypeDefinition{
		ObjectTypes: []ObjectTypeConfig{
			{Name: "configuration", DefaultPollInterval: 86400},
			{Name: "telemetry", DefaultPollInterval: 60},
			{Name: "status", DefaultPollInterval: 300},
			{Name: "performance", DefaultPollInterval: 120},
		},
		DefaultPollInterval: 300,
	}

	tests := []struct {
		name       string
		objectType string
		expected   uint64
	}{
		{
			name:       "configuration type",
			objectType: "configuration",
			expected:   86400,
		},
		{
			name:       "telemetry type",
			objectType: "telemetry",
			expected:   60,
		},
		{
			name:       "status type",
			objectType: "status",
			expected:   300,
		},
		{
			name:       "performance type",
			objectType: "performance",
			expected:   120,
		},
		{
			name:       "empty type",
			objectType: "",
			expected:   300, // Default poll interval
		},
		{
			name:       "invalid type",
			objectType: "invalid",
			expected:   300, // Default poll interval
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := otd.DefaultPollIntervalForType(tt.objectType)
			assert.Equal(t, tt.expected, result)
		})
	}
}

func TestObjectTypeDefinition_ValidateObjectType(t *testing.T) {
	otd := ObjectTypeDefinition{
		ObjectTypes: []ObjectTypeConfig{
			{Name: "configuration"},
			{Name: "telemetry"},
			{Name: "status"},
			{Name: "performance"},
		},
	}

	tests := []struct {
		name        string
		objectType  string
		expectError bool
	}{
		{
			name:        "valid object type",
			objectType:  "configuration",
			expectError: false,
		},
		{
			name:        "valid object type - telemetry",
			objectType:  "telemetry",
			expectError: false,
		},
		{
			name:        "empty object type",
			objectType:  "",
			expectError: false,
		},
		{
			name:        "invalid object type",
			objectType:  "invalid_type",
			expectError: true,
		},
		{
			name:        "case sensitive invalid",
			objectType:  "Configuration",
			expectError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := otd.ValidateObjectType(tt.objectType)
			if tt.expectError {
				assert.Error(t, err)
				assert.Contains(t, err.Error(), "invalid object type")
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

func TestObjectTypeDefinition_ValidatePollIntervals(t *testing.T) {
	otd := ObjectTypeDefinition{
		ObjectTypes: []ObjectTypeConfig{
			{Name: "configuration"},
			{Name: "telemetry"},
			{Name: "status"},
			{Name: "performance"},
		},
	}

	tests := []struct {
		name          string
		pollIntervals map[string]uint64
		expectError   bool
	}{
		{
			name: "valid poll intervals",
			pollIntervals: map[string]uint64{
				"configuration": 86400,
				"telemetry":     60,
				"status":        300,
				"performance":   120,
			},
			expectError: false,
		},
		{
			name:          "nil poll intervals",
			pollIntervals: nil,
			expectError:   false,
		},
		{
			name:          "empty poll intervals",
			pollIntervals: map[string]uint64{},
			expectError:   false,
		},
		{
			name: "invalid poll interval key",
			pollIntervals: map[string]uint64{
				"configuration": 86400,
				"invalid_type":  300,
			},
			expectError: true,
		},
		{
			name: "multiple invalid keys",
			pollIntervals: map[string]uint64{
				"invalid_type1": 86400,
				"invalid_type2": 300,
			},
			expectError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := otd.ValidatePollIntervals(tt.pollIntervals)
			if tt.expectError {
				assert.Error(t, err)
				assert.Contains(t, err.Error(), "invalid poll interval key")
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

func TestObjectTypeDefinition_YAMLUnmarshal(t *testing.T) {
	yamlContent := `
object_types:
  - name: configuration
    description: "Static objects that describe how the device is configured"
    default_poll_interval: 86400
  - name: telemetry
    description: "Dynamic objects that provide real-time telemetry data"
    default_poll_interval: 60
  - name: status
    description: "Objects that indicate device or interface status"
    default_poll_interval: 300
  - name: performance
    description: "Objects that provide performance metrics"
    default_poll_interval: 120
default_poll_interval: 300
`

	var otd ObjectTypeDefinition
	err := yaml.Unmarshal([]byte(yamlContent), &otd)
	require.NoError(t, err)

	assert.Len(t, otd.ObjectTypes, 4)
	assert.Equal(t, uint64(300), otd.DefaultPollInterval)

	// Test individual object types
	assert.Equal(t, "configuration", otd.ObjectTypes[0].Name)
	assert.Equal(t, uint64(86400), otd.ObjectTypes[0].DefaultPollInterval)
	assert.Equal(t, "Static objects that describe how the device is configured", otd.ObjectTypes[0].Description)

	assert.Equal(t, "telemetry", otd.ObjectTypes[1].Name)
	assert.Equal(t, uint64(60), otd.ObjectTypes[1].DefaultPollInterval)

	assert.Equal(t, "status", otd.ObjectTypes[2].Name)
	assert.Equal(t, uint64(300), otd.ObjectTypes[2].DefaultPollInterval)

	assert.Equal(t, "performance", otd.ObjectTypes[3].Name)
	assert.Equal(t, uint64(120), otd.ObjectTypes[3].DefaultPollInterval)

	// Test validation methods
	assert.True(t, otd.IsValidObjectType("configuration"))
	assert.True(t, otd.IsValidObjectType("telemetry"))
	assert.False(t, otd.IsValidObjectType("invalid"))

	err = otd.ValidateObjectType("configuration")
	assert.NoError(t, err)

	err = otd.ValidateObjectType("invalid")
	assert.Error(t, err)
}
