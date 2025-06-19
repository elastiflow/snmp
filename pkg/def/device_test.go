package def

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

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
