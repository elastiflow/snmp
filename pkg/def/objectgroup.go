package def

type ObjectGroup struct {
	Objects []string `yaml:"objects,omitempty" json:"objects,omitempty"`
}

func (o ObjectGroup) Validate() error {
	return nil
}

func (o ObjectGroup) Type() string {
	return "object_group"
}
