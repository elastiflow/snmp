package def

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestValidateDevices(t *testing.T) {
	tests := []struct {
		name         string
		devices      map[string]*Device
		deviceGroups map[string]DeviceGroup
		objectTypes  map[string]ObjectType
		expected     string
	}{
		{
			name:         "valid devices",
			devices:      map[string]*Device{"device1": device1, "device2": device2},
			deviceGroups: map[string]DeviceGroup{"dg": {ObjectGroups: []string{"og"}}},
			objectTypes:  map[string]ObjectType{"configuration": {PollInterval: 60}},
		},
		{
			name:     "invalid device",
			devices:  map[string]*Device{"device1": {IP: "127.0.0.1"}},
			expected: "missing properties: 'version', 'device_groups'",
		},
		{
			name:     "undefined device group",
			devices:  map[string]*Device{"device1": {IP: "127.0.0.1", DeviceGroups: []string{"dg"}}},
			expected: "device \"device1\" references an undefined device group: \"dg\"",
		},
		{
			name: "duplicate IP addresses",
			devices: map[string]*Device{
				"device1": {IP: "127.0.0.1"},
				"device2": {IP: "127.0.0.1"},
			},
			expected: "has the same IP address (127.0.0.1)",
		},
		{
			name:         "undefined object type",
			devices:      map[string]*Device{"device1": device1, "device2": device2},
			deviceGroups: map[string]DeviceGroup{"dg": {ObjectGroups: []string{"og"}}},
			expected:     "found 1 device definition errors:\ndevice \"device1\" references an undefined object type: \"configuration\"",
		},
		{
			name:    "empty device map",
			devices: map[string]*Device{},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := ValidateDevices(tt.devices, tt.deviceGroups, tt.objectTypes)
			if tt.expected == "" {
				assert.Nil(t, err)
			} else {
				assert.ErrorContains(t, err, tt.expected)
			}
		})
	}
}

func TestDevice_Validate(t *testing.T) {
	tests := []struct {
		name     string
		device   *Device
		expected error
	}{
		{
			name: "basic device",
			device: &Device{
				IP:           "192.168.1.1",
				Port:         161,
				Timeout:      5,
				Retries:      3,
				Version:      "2c",
				Communities:  []string{"public"},
				DeviceGroups: []string{"network_devices"},
				PollInterval: 60,
			},
			expected: nil,
		},
		{
			name: "device with v3 credentials",
			device: &Device{
				IP:           "192.168.1.2",
				Port:         161,
				Timeout:      10,
				Retries:      2,
				Version:      "3",
				DeviceGroups: []string{"secure_devices"},
				PollInterval: 300,
				V3Credentials: []V3Credential{
					{
						Username:                 "admin",
						AuthenticationProtocol:   "SHA",
						PrivacyProtocol:          "AES",
						AuthenticationPassphrase: "authpass",
						PrivacyPassphrase:        "privpass",
					},
				},
			},
			expected: nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equal(t, tt.expected, tt.device.Validate())
		})
	}
}

func TestDevice_Kind(t *testing.T) {
	device := Device{}
	assert.Equal(t, "device", device.Kind())
}

func TestDevice_applyHardcodedDefaults(t *testing.T) {
	tests := []struct {
		name        string
		device      *Device
		objectTypes map[string]ObjectType
		expected    *Device
	}{
		{
			name:   "empty device",
			device: &Device{},
			objectTypes: map[string]ObjectType{
				"interface": {PollInterval: 30},
				"system":    {PollInterval: 60},
			},
			expected: &Device{
				Port:               161,
				Timeout:            3000,
				Retries:            2,
				PollInterval:       60,
				MaxOIDs:            64,
				MaxConcurrentPolls: 4,
				PollIntervals: map[string]uint64{
					"interface": 30,
					"system":    60,
				},
			},
		},
		{
			name: "partially filled device",
			device: &Device{
				Port:         162,
				Timeout:      5000,
				PollInterval: 120,
			},
			objectTypes: map[string]ObjectType{
				"interface": {PollInterval: 30},
			},
			expected: &Device{
				Port:               162,
				Timeout:            5000,
				Retries:            2,
				PollInterval:       120,
				MaxOIDs:            64,
				MaxConcurrentPolls: 4,
				PollIntervals: map[string]uint64{
					"interface": 30,
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.device.applyHardcodedDefaults(tt.objectTypes)
			assert.Equal(t, tt.expected, tt.device)
		})
	}
}

func TestDevice_applyDefaults(t *testing.T) {
	tests := []struct {
		name          string
		device        *Device
		defaultDevice *Device
		expected      *Device
	}{
		{
			name:   "empty device",
			device: &Device{},
			defaultDevice: &Device{
				Port:               161,
				Timeout:            3000,
				Retries:            2,
				ExponentialTimeout: true,
				Version:            "2c",
				Communities:        []string{"public"},
				DeviceGroups:       []string{"default"},
				PollInterval:       60,
				MaxOIDs:            64,
				V3Credentials:      []V3Credential{{Username: "admin"}},
				MaxRepetitions:     10,
				MaxConcurrentPolls: 4,
				PollIntervals: map[string]uint64{
					"interface": 30,
					"system":    60,
				},
			},
			expected: &Device{
				Port:               161,
				Timeout:            3000,
				Retries:            2,
				ExponentialTimeout: true,
				Version:            "2c",
				Communities:        []string{"public"},
				DeviceGroups:       []string{"default"},
				PollInterval:       60,
				MaxOIDs:            64,
				V3Credentials:      []V3Credential{{Username: "admin"}},
				MaxRepetitions:     10,
				MaxConcurrentPolls: 4,
				PollIntervals: map[string]uint64{
					"interface": 30,
					"system":    60,
				},
			},
		},
		{
			name: "partially filled device",
			device: &Device{
				Port:         162,
				Timeout:      5000,
				PollInterval: 120,
				PollIntervals: map[string]uint64{
					"interface": 45,
				},
			},
			defaultDevice: &Device{
				Port:               161,
				Timeout:            3000,
				Retries:            2,
				ExponentialTimeout: true,
				Version:            "2c",
				Communities:        []string{"public"},
				DeviceGroups:       []string{"default"},
				PollInterval:       60,
				MaxOIDs:            64,
				V3Credentials:      []V3Credential{{Username: "admin"}},
				MaxRepetitions:     10,
				MaxConcurrentPolls: 4,
				PollIntervals: map[string]uint64{
					"interface": 30,
					"system":    60,
				},
			},
			expected: &Device{
				Port:               162,
				Timeout:            5000,
				Retries:            2,
				ExponentialTimeout: true,
				Version:            "2c",
				Communities:        []string{"public"},
				DeviceGroups:       []string{"default"},
				PollInterval:       120,
				MaxOIDs:            64,
				V3Credentials:      []V3Credential{{Username: "admin"}},
				MaxRepetitions:     10,
				MaxConcurrentPolls: 4,
				PollIntervals: map[string]uint64{
					"interface": 45,
					"system":    60,
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.device.applyDefaults(tt.defaultDevice)
			assert.Equal(t, tt.expected, tt.device)
		})
	}
}
