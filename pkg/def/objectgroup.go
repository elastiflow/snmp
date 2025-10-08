package def

type ObjectGroup struct {
	Objects []string `yaml:"objects,omitempty" json:"objects,omitempty"`
}

func (o ObjectGroup) Validate() error {
	return nil
}

func (o ObjectGroup) Kind() string {
	return "object_group"
}
