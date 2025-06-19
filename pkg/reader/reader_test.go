package reader

import (
	"testing"

	"github.com/elastiflow/snmp/pkg/def"
	"github.com/stretchr/testify/assert"
)

// TestReadSNMPTrapEnterprises tests the ReadSNMPTrapEnterprises function
func TestReadSNMPTrapEnterprises(t *testing.T) {
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
			enterprises, err := ReadSNMPTrapEnterprises("testdata/traps/rules", tt.enterpriseFile)

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

// TestReadSNMPDefinitions tests the ReadSNMPDefinitions function
func TestReadSNMPDefinitions(t *testing.T) {
	tests := []struct {
		name             string
		devicesDir       string
		deviceGroupsDir  string
		objectGroupsDir  string
		objectsDir       string
		wantErr          string
		wantDevices      map[string]def.Device
		wantDeviceGroups map[string]def.DeviceGroup
		wantObjectGroups map[string]def.ObjectGroup
		wantObjects      map[string]def.Object
	}{
		{
			name:            "valid definitions",
			devicesDir:      "testdata/devices",
			deviceGroupsDir: "testdata/device_groups",
			objectGroupsDir: "testdata/object_groups",
			objectsDir:      "testdata/objects",
			wantDevices: map[string]def.Device{
				"test_device": {
					IP:      "192.168.1.1",
					Port:    161,
					Version: "2c",
					Timeout: 5,
					Retries: 3,
				},
			},
			wantDeviceGroups: map[string]def.DeviceGroup{
				"test_device_group": {
					ObjectGroups: []string{"test_object_group"},
				},
			},
			wantObjectGroups: map[string]def.ObjectGroup{
				"test_object_group": {
					Objects: []string{"test_object"},
				},
			},
			wantObjects: map[string]def.Object{
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
			name:            "devices directory does not exist",
			devicesDir:      "testdata/nonexistent",
			deviceGroupsDir: "testdata/device_groups",
			objectGroupsDir: "testdata/object_groups",
			objectsDir:      "testdata/objects",
			wantErr:         "failed to read and validate devices: failed to walk directory",
		},
		{
			name:            "device groups directory does not exist",
			devicesDir:      "testdata/devices",
			deviceGroupsDir: "testdata/nonexistent",
			objectGroupsDir: "testdata/object_groups",
			objectsDir:      "testdata/objects",
			wantErr:         "failed to read and validate device groups: failed to walk directory",
		},
		{
			name:            "object groups directory does not exist",
			devicesDir:      "testdata/devices",
			deviceGroupsDir: "testdata/device_groups",
			objectGroupsDir: "testdata/nonexistent",
			objectsDir:      "testdata/objects",
			wantErr:         "failed to read and validate object groups: failed to walk directory",
		},
		{
			name:            "objects directory does not exist",
			devicesDir:      "testdata/devices",
			deviceGroupsDir: "testdata/device_groups",
			objectGroupsDir: "testdata/object_groups",
			objectsDir:      "testdata/nonexistent",
			wantErr:         "failed to read and validate objects: failed to walk directory",
		},
		{
			name:            "invalid definitions",
			devicesDir:      "testdata/devices_duplicate_ips",
			deviceGroupsDir: "testdata/device_groups",
			objectGroupsDir: "testdata/object_groups",
			objectsDir:      "testdata/objects",
			wantErr:         "failed to validate definitions",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			devices, deviceGroups, objectGroups, objects, err := ReadSNMPDefinitions(
				tt.devicesDir,
				tt.deviceGroupsDir,
				tt.objectGroupsDir,
				tt.objectsDir,
			)

			if tt.wantErr != "" {
				assert.Error(t, err)
				assert.ErrorContains(t, err, tt.wantErr)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.wantDevices, devices)
				assert.Equal(t, tt.wantDeviceGroups, deviceGroups)
				assert.Equal(t, tt.wantObjectGroups, objectGroups)
				assert.Equal(t, tt.wantObjects, objects)
			}
		})
	}
}

// TestReadSNMPEnums tests the ReadSNMPEnums function
func TestReadSNMPEnums(t *testing.T) {
	tests := []struct {
		name             string
		enumDir          string
		wantIntegerEnums map[string]def.IntegerEnum
		wantBitMapEnums  map[string]def.BitMapEnum
		wantOidEnums     map[string]def.OidEnum
		wantErr          string
	}{
		{
			name:    "valid enums",
			enumDir: "testdata/enums",
			wantIntegerEnums: map[string]def.IntegerEnum{
				"test_integer_enum": {
					1: "one",
					2: "two",
					3: "three",
				},
			},
			wantBitMapEnums: map[string]def.BitMapEnum{
				"test_bitmap_enum": {
					1: "one",
					2: "two",
					3: "three",
				},
			},
			wantOidEnums: map[string]def.OidEnum{
				"test_oid_enum": ".1.3.6.1.2.1",
			},
		},
		{
			name:    "missing integer enums",
			enumDir: "testdata/enums_missing_integer_enums",
			wantErr: "failed to read and validate integer enums: failed to read file testdata/enums_missing_integer_enums/test_integer_enum.yml",
		},
		{
			name:    "missing bit map enums",
			enumDir: "testdata/enums_missing_bitmap_enums",
			wantErr: "failed to read and validate bit map enums: failed to read file testdata/enums_missing_bitmap_enums/test_bitmap_enum.yml",
		},
		{
			name:    "missing oid enums",
			enumDir: "testdata/enums_missing_oid_enums",
			wantErr: "failed to read and validate oid enums: failed to read file testdata/enums_missing_oid_enums/test_oid_enum.yml",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			integerEnums, bitMapEnums, oidEnums, err := ReadSNMPEnums(tt.enumDir)
			if tt.wantErr != "" {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.wantIntegerEnums, integerEnums)
				assert.Equal(t, tt.wantBitMapEnums, bitMapEnums)
				assert.Equal(t, tt.wantOidEnums, oidEnums)
			}
		})
	}
}

// TestValidateSNMPDefinitions tests the validateSNMPDefinitions function
func TestValidateSNMPDefinitions(t *testing.T) {
	tests := []struct {
		name         string
		devices      map[string]def.Device
		deviceGroups map[string]def.DeviceGroup
		objectGroups map[string]def.ObjectGroup
		objects      map[string]def.Object
		wantErr      string
	}{
		{
			name: "valid definitions",
			devices: map[string]def.Device{
				"device1": {
					IP: "192.168.1.1",
				},
				"device2": {
					IP: "192.168.1.2",
				},
			},
			deviceGroups: map[string]def.DeviceGroup{
				"device_group1": {
					ObjectGroups: []string{"object_group1"},
				},
			},
			objectGroups: map[string]def.ObjectGroup{
				"object_group1": {
					Objects: []string{"object1"},
				},
			},
			objects: map[string]def.Object{
				"object1": {},
			},
		},
		{
			name: "duplicate IP address",
			devices: map[string]def.Device{
				"device1": {
					IP: "192.168.1.1",
				},
				"device2": {
					IP: "192.168.1.1", // Same IP as device1
				},
			},
			deviceGroups: map[string]def.DeviceGroup{
				"device_group1": {
					ObjectGroups: []string{"object_group1"},
				},
			},
			objectGroups: map[string]def.ObjectGroup{
				"object_group1": {
					Objects: []string{"object1"},
				},
			},
			objects: map[string]def.Object{
				"object1": {},
			},
			wantErr: "device device2 has the same IP address (192.168.1.1) as device device1",
		},
		{
			name: "undefined object group",
			devices: map[string]def.Device{
				"device1": {
					IP: "192.168.1.1",
				},
			},
			deviceGroups: map[string]def.DeviceGroup{
				"device_group1": {
					ObjectGroups: []string{"undefined_object_group"}, // This object group doesn't exist
				},
			},
			objectGroups: map[string]def.ObjectGroup{
				"object_group1": {
					Objects: []string{"object1"},
				},
			},
			objects: map[string]def.Object{
				"object1": {},
			},
			wantErr: "device group device_group1 references an undefined object group: undefined_object_group",
		},
		{
			name: "undefined object",
			devices: map[string]def.Device{
				"device1": {
					IP: "192.168.1.1",
				},
			},
			deviceGroups: map[string]def.DeviceGroup{
				"device_group1": {
					ObjectGroups: []string{"object_group1"},
				},
			},
			objectGroups: map[string]def.ObjectGroup{
				"object_group1": {
					Objects: []string{"undefined_object"}, // This object doesn't exist
				},
			},
			objects: map[string]def.Object{
				"object1": {},
			},
			wantErr: "object group object_group1 references an undefined object: undefined_object",
		},
		{
			name: "multiple validation failures",
			devices: map[string]def.Device{
				"device1": {
					IP: "192.168.1.1",
				},
				"device2": {
					IP: "192.168.1.1", // Same IP as device1
				},
			},
			deviceGroups: map[string]def.DeviceGroup{
				"device_group1": {
					ObjectGroups: []string{"undefined_object_group"}, // This object group doesn't exist
				},
			},
			objectGroups: map[string]def.ObjectGroup{
				"object_group1": {
					Objects: []string{"undefined_object"}, // This object doesn't exist
				},
			},
			objects: map[string]def.Object{
				"object1": {},
			},
			wantErr: "found 4 invalid definitions",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := validateSNMPDefinitions(tt.devices, tt.deviceGroups, tt.objectGroups, tt.objects)

			if tt.wantErr != "" {
				assert.Error(t, err)
				assert.ErrorContains(t, err, tt.wantErr)
			} else {
				assert.NoError(t, err)
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
			dirPath: "testdata/objects",
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
			dirPath: "testdata/objects_unreadable",
			wantErr: "failed to read definitions",
		},
		{
			name:    "duplicate object keys",
			dirPath: "testdata/objects_duplicate_keys",
			wantErr: "failed to unmarshal file",
		},
		{
			name:    "duplicate objects",
			dirPath: "testdata/objects_duplicate",
			wantErr: "found duplicate object with id test_object",
		},
		{
			name:    "invalid object definition",
			dirPath: "testdata/objects_invalid",
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
