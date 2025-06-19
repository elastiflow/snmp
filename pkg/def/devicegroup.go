package def

type DeviceGroup struct {
	ObjectGroups []string `yaml:"object_groups,omitempty" json:"object_groups,omitempty"`
}

func (d DeviceGroup) Validate() error {
	return nil
}

func (d DeviceGroup) Type() string {
	return "device_group"
}
