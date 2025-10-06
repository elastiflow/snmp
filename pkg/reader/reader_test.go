package reader

import (
	"testing"

	"github.com/elastiflow/snmp/pkg/def"
	"github.com/stretchr/testify/assert"
)

func TestReadTrapEnterprises(t *testing.T) {
	tests := []struct {
		name           string
		enterpriseFile string
		wantErr        string
		want           map[string]def.Enterprise
	}{
		{
			name:           "valid enterprises file",
			enterpriseFile: "testdata/traps/enterprises_valid.yml",
			want:           map[string]def.Enterprise{".1.2.840.10036.1.6": "valid.yml"},
		},
		{
			name:           "enterprise file does not exist",
			enterpriseFile: "testdata/traps/enterprises_does_not_exist.yml",
			wantErr:        "failed to read file testdata/traps/enterprises_does_not_exist.yml",
		},
		{
			name:           "enterprise file cannot be unmarshalled",
			enterpriseFile: "testdata/traps/enterprises_invalid.yml",
			wantErr:        "failed to unmarshal file testdata/traps/enterprises_invalid.yml",
		},
		{
			name:           "invalid rule in enterprises file",
			enterpriseFile: "testdata/traps/enterprises_invalid_rule.yml",
			wantErr:        "failed to validate enterprise .1.2.840.10036.1.6",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			enterprises, err := ReadTrapEnterprises("testdata/traps/rules", tt.enterpriseFile)

			if tt.wantErr != "" {
				assert.Error(t, err)
				assert.ErrorContains(t, err, tt.wantErr)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.want, enterprises)
			}
		})
	}
}

func TestReadDefinitions(t *testing.T) {
	tests := []struct {
		name            string
		devicesDir      string
		deviceGroupsDir string
		objectGroupsDir string
		objectsDir      string
		defaultsDir     string
		wantErr         string
		wantDefinitions *def.Definitions
	}{
		{
			name:            "valid definitions",
			devicesDir:      "testdata/valid/devices",
			deviceGroupsDir: "testdata/valid/device_groups",
			objectGroupsDir: "testdata/valid/object_groups",
			objectsDir:      "testdata/valid/objects",
			defaultsDir:     "testdata/valid/defaults",
			wantDefinitions: &def.Definitions{
				Devices: map[string]def.Device{
					"test_device": {
						IP:            "192.168.1.1",
						Port:          161,
						Version:       "2c",
						Timeout:       5,
						Retries:       3,
						PollIntervals: map[string]uint64{"configuration": 60},
					},
				},
				DeviceGroups: map[string]def.DeviceGroup{
					"test_device_group": {
						ObjectGroups: []string{"test_object_group"},
					},
				},
				ObjectGroups: map[string]def.ObjectGroup{
					"test_object_group": {
						Objects: []string{"test_object"},
					},
				},
				Objects: map[string]def.Object{
					"test_object": {
						Mib:                "TEST-MIB",
						ObjectName:         "sysDescr",
						DiscoveryAttribute: "sysDescr",
						Attributes: def.ObjectAttributes{
							"sysDescr": {
								Oid:    ".1.3.6.1.2.1.1.1.0",
								Name:   "System Description",
								Syntax: "DisplayString",
							},
						},
					},
				},
				ObjectTypes: map[string]def.ObjectType{
					"configuration": {PollInterval: 86400},
					"telemetry":     {PollInterval: 60},
				},
				DefaultDevice: &def.Device{
					Port:         161,
					Timeout:      3000,
					Retries:      2,
					Version:      "2c",
					Communities:  []string{"public"},
					PollInterval: 60,
				},
			},
		},
		{
			name:            "valid definitions without defaults",
			devicesDir:      "testdata/valid/devices_without_object_type",
			deviceGroupsDir: "testdata/valid/device_groups",
			objectGroupsDir: "testdata/valid/object_groups",
			objectsDir:      "testdata/valid/objects",
			defaultsDir:     "testdata/invalid/defaults",
			wantDefinitions: &def.Definitions{
				Devices: map[string]def.Device{
					"test_device": {
						IP:      "192.168.1.1",
						Port:    161,
						Version: "2c",
						Timeout: 5,
						Retries: 3,
					},
				},
				DeviceGroups: map[string]def.DeviceGroup{
					"test_device_group": {
						ObjectGroups: []string{"test_object_group"},
					},
				},
				ObjectGroups: map[string]def.ObjectGroup{
					"test_object_group": {
						Objects: []string{"test_object"},
					},
				},
				Objects: map[string]def.Object{
					"test_object": {
						Mib:                "TEST-MIB",
						ObjectName:         "sysDescr",
						DiscoveryAttribute: "sysDescr",
						Attributes: def.ObjectAttributes{
							"sysDescr": {
								Oid:    ".1.3.6.1.2.1.1.1.0",
								Name:   "System Description",
								Syntax: "DisplayString",
							},
						},
					},
				},
			},
		},
		{
			name:       "devices directory does not exist",
			devicesDir: "testdata/nonexistent",
			wantErr:    "failed to read and validate devices: failed to walk directory",
		},
		{
			name:            "device groups directory does not exist",
			devicesDir:      "testdata/valid/devices",
			deviceGroupsDir: "testdata/nonexistent",
			wantErr:         "failed to read and validate device groups: failed to walk directory",
		},
		{
			name:            "object groups directory does not exist",
			devicesDir:      "testdata/valid/devices",
			deviceGroupsDir: "testdata/valid/device_groups",
			objectGroupsDir: "testdata/nonexistent",
			wantErr:         "failed to read and validate object groups: failed to walk directory",
		},
		{
			name:            "objects directory does not exist",
			devicesDir:      "testdata/valid/devices",
			deviceGroupsDir: "testdata/valid/device_groups",
			objectGroupsDir: "testdata/valid/object_groups",
			objectsDir:      "testdata/nonexistent",
			wantErr:         "failed to read and validate objects: failed to walk directory",
		},
		{
			name:            "invalid definitions",
			devicesDir:      "testdata/invalid/devices",
			deviceGroupsDir: "testdata/valid/device_groups",
			objectGroupsDir: "testdata/valid/object_groups",
			objectsDir:      "testdata/valid/objects",
			defaultsDir:     "testdata/valid/defaults",
			wantErr:         "failed to validate definitions",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := ReadDefinitions(
				tt.devicesDir,
				tt.deviceGroupsDir,
				tt.objectGroupsDir,
				tt.objectsDir,
				tt.defaultsDir,
			)

			if tt.wantErr != "" {
				assert.Error(t, err)
				assert.ErrorContains(t, err, tt.wantErr)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.wantDefinitions, got)
			}
		})
	}
}

func TestReadEnums(t *testing.T) {
	tests := []struct {
		name    string
		enumDir string
		want    *def.Enums
		wantErr string
	}{
		{
			name:    "valid enums",
			enumDir: "testdata/valid/enums",
			want: &def.Enums{
				Integers: map[string]def.IntegerEnum{".1.2.3.4": {1: "one", 2: "two", 3: "three"}},
				BitMaps:  map[string]def.BitMapEnum{".1.2.3.4": {1: "one", 2: "two", 3: "three"}},
				Oids:     map[string]def.OidEnum{".1.2.3.4": ".1.3.6.1.2.1"},
			},
		},
		{
			name:    "missing integer enums",
			enumDir: "testdata/invalid/enums_missing_integer_enums",
			wantErr: "failed to read and validate integer enums: failed to walk directory: error while walking directory",
		},
		{
			name:    "missing bit map enums",
			enumDir: "testdata/invalid/enums_missing_bitmap_enums",
			wantErr: "failed to read and validate bit map enums: failed to walk directory: error while walking directory",
		},
		{
			name:    "missing oid enums",
			enumDir: "testdata/invalid/enums_missing_oid_enums",
			wantErr: "failed to read and validate oid enums: failed to walk directory: error while walking directory",
		},
		{
			name:    "invalid integer oid",
			enumDir: "testdata/invalid/enums_invalid_integer_oid",
			wantErr: "invalid oid invalid_oid in integer enums",
		},
		{
			name:    "invalid bit map oid",
			enumDir: "testdata/invalid/enums_invalid_bitmap_oid",
			wantErr: "invalid oid invalid_oid in bit map enums",
		},
		{
			name:    "invalid oid",
			enumDir: "testdata/invalid/enums_invalid_oid_oid",
			wantErr: "invalid oid invalid_oid in oid enums",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := ReadEnums(tt.enumDir)
			if tt.wantErr != "" {
				assert.Error(t, err)
				assert.ErrorContains(t, err, tt.wantErr)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.want, got)
			}
		})
	}
}

// TestRead tests the Read function
func TestRead(t *testing.T) {
	tests := []struct {
		name    string
		dirPath string
		wantErr string
		want    map[string]def.Object
	}{
		{
			name:    "valid objects",
			dirPath: "testdata/valid/objects",
			want: map[string]def.Object{
				"test_object": {
					Mib:                "TEST-MIB",
					ObjectName:         "sysDescr",
					DiscoveryAttribute: "sysDescr",
					Attributes: def.ObjectAttributes{
						"sysDescr": {
							Oid:    ".1.3.6.1.2.1.1.1.0",
							Name:   "System Description",
							Syntax: "DisplayString",
						},
					},
				},
			},
		},
		{
			name:    "failed to walk directory",
			dirPath: "testdata/does_not_exist",
			wantErr: "failed to walk directory: error while walking directory 'testdata/does_not_exist'",
		},
		{
			name:    "failed to read file",
			dirPath: "testdata/invalid/objects_unreadable",
			wantErr: "failed to read definitions",
		},
		{
			name:    "duplicate object keys",
			dirPath: "testdata/invalid/objects_duplicate_keys",
			wantErr: "failed to unmarshal file",
		},
		{
			name:    "duplicate objects",
			dirPath: "testdata/invalid/objects_duplicate",
			wantErr: "found duplicate object with id test_object",
		},
		{
			name:    "invalid object definition",
			dirPath: "testdata/invalid/objects_invalid",
			wantErr: "failed to validate object with id test_object_invalid: jsonschema: '/attributes/sysDescr/syntax' does not validate",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			devices, err := Read[def.Object](tt.dirPath)

			if tt.wantErr != "" {
				assert.Error(t, err)
				assert.ErrorContains(t, err, tt.wantErr)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.want, devices)
			}
		})
	}
}
