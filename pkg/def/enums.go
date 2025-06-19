package def

import "fmt"

type IntegerEnum map[int]string

func (i IntegerEnum) Type() string {
	return "integer_enum"
}

func (i IntegerEnum) Validate() error {
	for integer, value := range i {
		if value == "" {
			return fmt.Errorf("integer enum value for %d cannot be empty (null values must be stringified)", integer)
		}
	}
	return nil
}

type BitMapEnum map[int]string

func (b BitMapEnum) Type() string {
	return "bit_map_enum"
}

func (b BitMapEnum) Validate() error {
	for integer, value := range b {
		if value == "" {
			return fmt.Errorf("integer enum value for %d cannot be empty (null values must be stringified)", integer)
		}
	}
	return nil
}

type OidEnum string

func (o OidEnum) Type() string {
	return "oid_enum"
}

func (o OidEnum) Validate() error {
	if o == "" {
		return fmt.Errorf("oid enum value cannot be empty")
	}
	return nil
}
