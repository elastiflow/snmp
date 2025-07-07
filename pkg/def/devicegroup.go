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

// DeviceGroupName is a string that represents the name of a device group.
type DeviceGroupName string

func (d DeviceGroupName) Type() string {
	return "device_group"
}

func (d DeviceGroupName) Validate() error {
	if d == "" {
		return fmt.Errorf("device group cannot be empty")
	}
	return nil
}
