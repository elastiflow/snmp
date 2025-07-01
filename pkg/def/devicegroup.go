package def

import "fmt"

type DeviceGroup struct {
	ObjectGroups []string `yaml:"object_groups,omitempty" json:"object_groups,omitempty"`
}

func (d DeviceGroup) Validate() error {
	return nil
}

func (d DeviceGroup) Type() string {
	return "device_group"
}

// SysObjectID is a string that represents the device group for a specific sysObjectID.
type SysObjectID string

func (s SysObjectID) Type() string {
	return "sys_object_id"
}

func (s SysObjectID) Validate() error {
	if s == "" {
		return fmt.Errorf("device group cannot be empty")
	}
	return nil
}
