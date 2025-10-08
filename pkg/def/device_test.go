package def

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestValidateDevices(t *testing.T) {
	tests := []struct {
		name         string
		devices      map[string]Device
		deviceGroups map[string]DeviceGroup
		objectTypes  map[string]ObjectType
		expected     string
	}{
		{
			name: "valid devices",
			devices: map[string]Device{
				"device1": {IP: "127.0.0.1", DeviceGroups: []string{"dg"}, PollIntervals: map[string]uint64{"configuration": 60}},
				"device2": {IP: "127.0.0.2", DeviceGroups: []string{"dg"}},
			},
			deviceGroups: map[string]DeviceGroup{"dg": {ObjectGroups: []string{"og"}}},
			objectTypes:  map[string]ObjectType{"configuration": {PollInterval: 60}},
		},
		{
			name:     "undefined device group",
			devices:  map[string]Device{"device1": {IP: "127.0.0.1", DeviceGroups: []string{"dg"}}},
			expected: "device \"device1\" references an undefined device group: \"dg\"",
		},
		{
			name: "duplicate IP addresses",
			devices: map[string]Device{
				"device1": {IP: "127.0.0.1"},
				"device2": {IP: "127.0.0.1"},
			},
			expected: "has the same IP address (127.0.0.1)",
		},
		{
			name: "undefined object type",
			devices: map[string]Device{
				"device1": {IP: "127.0.0.1", DeviceGroups: []string{"dg"}, PollIntervals: map[string]uint64{"configuration": 60}},
			},
			deviceGroups: map[string]DeviceGroup{"dg": {ObjectGroups: []string{"og"}}},
			expected:     "found 1 device definition errors:\ndevice \"device1\" references an undefined object type: \"configuration\"",
		},
		{
			name:    "empty device map",
			devices: map[string]Device{},
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

func TestDevice_Type(t *testing.T) {
	device := Device{}
	assert.Equal(t, "device", device.Kind())
}
