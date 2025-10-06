package def

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestDeviceGroup_Validate(t *testing.T) {
	tests := []struct {
		name        string
		deviceGroup *DeviceGroup
		expected    error
	}{
		{
			name: "empty device group",
			deviceGroup: &DeviceGroup{
				ObjectGroups: []string{},
			},
			expected: nil,
		},
		{
			name: "device group with object groups",
			deviceGroup: &DeviceGroup{
				ObjectGroups: []string{"network_interfaces", "system_metrics"},
			},
			expected: nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equal(t, tt.expected, tt.deviceGroup.Validate())
		})
	}
}

func TestDeviceGroup_Type(t *testing.T) {
	deviceGroup := DeviceGroup{}
	assert.Equal(t, "device_group", deviceGroup.Kind())
}

func TestDeviceGroupName_Type(t *testing.T) {
	dg := DeviceGroupName("device_group")
	assert.Equal(t, "device_group", dg.Kind())
}

func TestSysObjectID_Validate(t *testing.T) {
	tests := []struct {
		name            string
		deviceGroupName DeviceGroupName
		expected        error
	}{
		{
			name:            "valid deviceGroupName",
			deviceGroupName: DeviceGroupName("deviceGroupName"),
			expected:        nil,
		},
		{
			name:            "empty deviceGroupName",
			deviceGroupName: DeviceGroupName(""),
			expected:        fmt.Errorf("device group cannot be empty"),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equal(t, tt.expected, tt.deviceGroupName.Validate())
		})
	}
}
