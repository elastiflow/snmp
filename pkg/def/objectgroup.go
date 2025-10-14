package def

import "errors"

type ObjectGroup struct {
	Objects []string `yaml:"objects,omitempty" json:"objects,omitempty"`
}

func (o ObjectGroup) Validate() error {
	if len(o.Objects) == 0 {
		return errors.New("objects is required")
	}
	return nil
}

func (o ObjectGroup) Kind() string {
	return "object_group"
}
