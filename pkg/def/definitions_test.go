package def

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestDefinitions_Validate(t *testing.T) {
	tests := []struct {
		name        string
		definitions Definitions
		wantErr     string
	}{
		{
			name: "valid definitions",
			definitions: Definitions{
				Devices: map[string]Device{
					"device1": {IP: "127.0.0.1", DeviceGroups: []string{"dg"}, PollIntervals: map[string]uint64{"configuration": 60}},
					"device2": {IP: "127.0.0.2", DeviceGroups: []string{"dg"}},
				},
				DeviceGroups: map[string]DeviceGroup{"dg": {ObjectGroups: []string{"og"}}},
				ObjectGroups: map[string]ObjectGroup{"og": {Objects: []string{"object1", "object2"}}},
				Objects: map[string]Object{
					"object1": {Mib: "IF-MIB", ObjectName: "system"},
					"object2": {Mib: "IF-MIB", ObjectName: "ifEntry"},
				},
				ObjectTypes: map[string]ObjectType{
					"configuration": {PollInterval: 60},
					"performance":   {PollInterval: 300},
				},
			},
		},
		{
			name: "invalid device",
			definitions: Definitions{
				Devices: map[string]Device{
					"device1": {IP: "127.0.0.1", DeviceGroups: []string{"dg"}, PollIntervals: map[string]uint64{"configuration": 60}},
					"device2": {IP: "127.0.0.1", DeviceGroups: []string{"dg"}},
				},
			},
			wantErr: "failed to validate devices",
		},
		{
			name: "invalid device group",
			definitions: Definitions{
				Devices: map[string]Device{
					"device1": {IP: "127.0.0.1", DeviceGroups: []string{"dg"}},
					"device2": {IP: "127.0.0.2", DeviceGroups: []string{"dg"}},
				},
				DeviceGroups: map[string]DeviceGroup{"dg": {ObjectGroups: []string{"og"}}},
			},
			wantErr: "device group dg references an undefined object group: og",
		},
		{
			name: "invalid object group",
			definitions: Definitions{
				Devices: map[string]Device{
					"device1": {IP: "127.0.0.1", DeviceGroups: []string{"dg"}},
					"device2": {IP: "127.0.0.2", DeviceGroups: []string{"dg"}},
				},
				DeviceGroups: map[string]DeviceGroup{"dg": {ObjectGroups: []string{"og"}}},
				ObjectGroups: map[string]ObjectGroup{"og": {Objects: []string{"object1", "object2"}}},
			},
			wantErr: "object group og references an undefined object: object1",
		},
		{
			name: "invalid object",
			definitions: Definitions{
				Devices: map[string]Device{
					"device1": {IP: "127.0.0.1", DeviceGroups: []string{"dg"}},
					"device2": {IP: "127.0.0.2", DeviceGroups: []string{"dg"}},
				},
				DeviceGroups: map[string]DeviceGroup{"dg": {ObjectGroups: []string{"og"}}},
				ObjectGroups: map[string]ObjectGroup{"og": {Objects: []string{"object1"}}},
				Objects: map[string]Object{
					"object1": {Mib: "IF-MIB", ObjectName: "system", Type: "object_type"},
				},
			},
			wantErr: "object system references an undefined object type: object_type",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.definitions.Validate()
			if tt.wantErr != "" {
				assert.ErrorContains(t, err, tt.wantErr)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}
