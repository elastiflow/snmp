package def

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

var device1 = &Device{
	IP:            "127.0.0.1",
	Port:          161,
	Version:       "2c",
	Communities:   []string{"public"},
	DeviceGroups:  []string{"dg"},
	PollIntervals: map[string]uint64{"configuration": 60},
}

var device2 = &Device{
	IP:      "fc00::",
	Port:    161,
	Version: "3",
	V3Credentials: []V3Credential{
		{
			Username:                 "username",
			AuthenticationProtocol:   "sha",
			AuthenticationPassphrase: "passphrase",
			PrivacyProtocol:          "aes",
			PrivacyPassphrase:        "privacy passphrase",
		},
	},
	DeviceGroups: []string{"dg"},
}

func TestDefinitions_Validate(t *testing.T) {
	tests := []struct {
		name        string
		definitions Definitions
		wantErr     string
	}{
		{
			name: "valid definitions",
			definitions: Definitions{
				Devices:      map[string]*Device{"device1": device1, "device2": device2},
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
				Devices: map[string]*Device{
					"device1": {IP: "127.0.0.1", DeviceGroups: []string{"dg"}, PollIntervals: map[string]uint64{"configuration": 60}},
					"device2": {IP: "127.0.0.1", DeviceGroups: []string{"dg"}},
				},
			},
			wantErr: "failed to validate devices",
		},
		{
			name: "invalid device group",
			definitions: Definitions{
				Devices:      map[string]*Device{"device1": device1, "device2": device2},
				DeviceGroups: map[string]DeviceGroup{"dg": {ObjectGroups: []string{"og"}}},
				ObjectTypes: map[string]ObjectType{
					"configuration": {PollInterval: 60},
					"performance":   {PollInterval: 300},
				},
			},
			wantErr: "device group \"dg\" references an undefined object group: \"og\"",
		},
		{
			name: "invalid object group",
			definitions: Definitions{
				Devices:      map[string]*Device{"device1": device1, "device2": device2},
				DeviceGroups: map[string]DeviceGroup{"dg": {ObjectGroups: []string{"og"}}},
				ObjectGroups: map[string]ObjectGroup{"og": {Objects: []string{"object1", "object2"}}},
				ObjectTypes: map[string]ObjectType{
					"configuration": {PollInterval: 60},
					"performance":   {PollInterval: 300},
				},
			},
			wantErr: "object group \"og\" references an undefined object: \"object1\"",
		},
		{
			name: "invalid object",
			definitions: Definitions{
				Devices:      map[string]*Device{"device1": device1, "device2": device2},
				DeviceGroups: map[string]DeviceGroup{"dg": {ObjectGroups: []string{"og"}}},
				ObjectGroups: map[string]ObjectGroup{"og": {Objects: []string{"object1"}}},
				Objects: map[string]Object{
					"object1": {Mib: "IF-MIB", ObjectName: "system", Type: "object_type"},
				},
				ObjectTypes: map[string]ObjectType{
					"configuration": {PollInterval: 60},
					"performance":   {PollInterval: 300},
				},
			},
			wantErr: "object \"system\" references an undefined object type: \"object_type\"",
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

func TestDefinitions_ApplyDefaults(t *testing.T) {
	tests := []struct {
		name        string
		definitions Definitions
		want        Definitions
	}{
		{
			name: "nil DefaultDevice",
			definitions: Definitions{
				Devices: map[string]*Device{
					"device1": {
						IP:           "127.0.0.1",
						DeviceGroups: []string{"dg1"},
					},
				},
				ObjectTypes: map[string]ObjectType{
					"configuration": {PollInterval: 60},
					"performance":   {PollInterval: 300},
				},
			},
			want: Definitions{
				Devices: map[string]*Device{
					"device1": {
						IP:                 "127.0.0.1",
						Port:               161,
						Timeout:            3000,
						Retries:            2,
						PollInterval:       60,
						MaxOIDs:            64,
						MaxConcurrentPolls: 4,
						DeviceGroups:       []string{"dg1"},
						PollIntervals: map[string]uint64{
							"configuration": 60,
							"performance":   300,
						},
					},
				},
				DefaultDevice: &Device{
					Port:               161,
					Timeout:            3000,
					Retries:            2,
					PollInterval:       60,
					MaxOIDs:            64,
					MaxConcurrentPolls: 4,
					PollIntervals: map[string]uint64{
						"configuration": 60,
						"performance":   300,
					},
				},
				ObjectTypes: map[string]ObjectType{
					"configuration": {PollInterval: 60},
					"performance":   {PollInterval: 300},
				},
			},
		},
		{
			name: "existing DefaultDevice",
			definitions: Definitions{
				Devices: map[string]*Device{
					"device1": {
						IP:           "127.0.0.1",
						DeviceGroups: []string{"dg1"},
					},
				},
				DefaultDevice: &Device{
					Port:           162,
					Timeout:        5000,
					Retries:        3,
					PollInterval:   120,
					Communities:    []string{"public"},
					DeviceGroups:   []string{"default"},
					Version:        "2c",
					MaxRepetitions: 10,
				},
				ObjectTypes: map[string]ObjectType{
					"configuration": {PollInterval: 60},
					"performance":   {PollInterval: 300},
				},
			},
			want: Definitions{
				Devices: map[string]*Device{
					"device1": {
						IP:                 "127.0.0.1",
						Port:               162,
						Timeout:            5000,
						Retries:            3,
						PollInterval:       120,
						MaxOIDs:            64,
						MaxConcurrentPolls: 4,
						DeviceGroups:       []string{"dg1"},
						Version:            "2c",
						Communities:        []string{"public"},
						MaxRepetitions:     10,
						PollIntervals: map[string]uint64{
							"configuration": 60,
							"performance":   300,
						},
					},
				},
				DefaultDevice: &Device{
					Port:               162,
					Timeout:            5000,
					Retries:            3,
					PollInterval:       120,
					MaxOIDs:            64,
					MaxConcurrentPolls: 4,
					Communities:        []string{"public"},
					DeviceGroups:       []string{"default"},
					Version:            "2c",
					MaxRepetitions:     10,
					PollIntervals: map[string]uint64{
						"configuration": 60,
						"performance":   300,
					},
				},
				ObjectTypes: map[string]ObjectType{
					"configuration": {PollInterval: 60},
					"performance":   {PollInterval: 300},
				},
			},
		},
		{
			name: "multiple devices with different settings",
			definitions: Definitions{
				Devices: map[string]*Device{
					"device1": {
						IP:           "127.0.0.1",
						DeviceGroups: []string{"dg1"},
						Port:         8161,
					},
					"device2": {
						IP:           "127.0.0.2",
						DeviceGroups: []string{"dg2"},
						Timeout:      1000,
					},
				},
				ObjectTypes: map[string]ObjectType{
					"configuration": {PollInterval: 60},
					"performance":   {PollInterval: 300},
				},
			},
			want: Definitions{
				Devices: map[string]*Device{
					"device1": {
						IP:                 "127.0.0.1",
						Port:               8161, // Custom value preserved
						Timeout:            3000,
						Retries:            2,
						PollInterval:       60,
						MaxOIDs:            64,
						MaxConcurrentPolls: 4,
						DeviceGroups:       []string{"dg1"},
						PollIntervals: map[string]uint64{
							"configuration": 60,
							"performance":   300,
						},
					},
					"device2": {
						IP:                 "127.0.0.2",
						Port:               161,
						Timeout:            1000, // Custom value preserved
						Retries:            2,
						PollInterval:       60,
						MaxOIDs:            64,
						MaxConcurrentPolls: 4,
						DeviceGroups:       []string{"dg2"},
						PollIntervals: map[string]uint64{
							"configuration": 60,
							"performance":   300,
						},
					},
				},
				DefaultDevice: &Device{
					Port:               161,
					Timeout:            3000,
					Retries:            2,
					PollInterval:       60,
					MaxOIDs:            64,
					MaxConcurrentPolls: 4,
					PollIntervals: map[string]uint64{
						"configuration": 60,
						"performance":   300,
					},
				},
				ObjectTypes: map[string]ObjectType{
					"configuration": {PollInterval: 60},
					"performance":   {PollInterval: 300},
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.definitions.ApplyDefaults()
			assert.Equal(t, tt.want, tt.definitions)
		})
	}
}
