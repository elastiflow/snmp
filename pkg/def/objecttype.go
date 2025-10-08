package def

type ObjectType struct {
	PollInterval uint64 `yaml:"poll_interval,omitempty" json:"poll_interval,omitempty"`
}

func (o ObjectType) Validate() error {
	return nil
}

func (o ObjectType) Kind() string {
	return "object_type"
}
