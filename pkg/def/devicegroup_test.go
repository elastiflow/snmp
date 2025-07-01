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
	assert.Equal(t, "device_group", deviceGroup.Type())
}

func TestSysObjectID_Type(t *testing.T) {
	sysObjectID := SysObjectID("sysObjectID")
	assert.Equal(t, "sys_object_id", sysObjectID.Type())
}

func TestSysObjectID_Validate(t *testing.T) {
	tests := []struct {
		name        string
		sysObjectID SysObjectID
		expected    error
	}{
		{
			name:        "valid sysObjectID",
			sysObjectID: SysObjectID("sysObjectID"),
			expected:    nil,
		},
		{
			name:        "empty sysObjectID",
			sysObjectID: SysObjectID(""),
			expected:    fmt.Errorf("device group cannot be empty"),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equal(t, tt.expected, tt.sysObjectID.Validate())
		})
	}
}
