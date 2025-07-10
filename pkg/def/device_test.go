package def

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestValidateDevices(t *testing.T) {
	tests := []struct {
		name     string
		devices  map[string]Device
		expected error
	}{
		{
			name: "valid devices",
			devices: map[string]Device{
				"device1": {
					IP: "192.168.1.1",
				},
				"device2": {
					IP: "192.168.1.2",
				},
			},
		},
		{
			name: "duplicate IP addresses",
			devices: map[string]Device{
				"device1": {
					IP: "192.168.1.1",
				},
				"device2": {
					IP: "192.168.1.1", // Same IP as device1
				},
			},
			expected: fmt.Errorf("found 2 invalid device definitions:\ndevice device1 has the same IP address (192.168.1.1) as device device2\ndevice device2 has the same IP address (192.168.1.1) as device device1"),
		},
		{
			name:    "empty device map",
			devices: map[string]Device{},
		},
		{
			name: "single device",
			devices: map[string]Device{
				"device1": {
					IP: "192.168.1.1",
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := ValidateDevices(tt.devices)
			if tt.expected == nil {
				assert.Nil(t, err)
			} else {
				assert.Equal(t, tt.expected.Error(), err.Error())
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
	assert.Equal(t, "device", device.Type())
}
