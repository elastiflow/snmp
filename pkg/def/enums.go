package def

import "fmt"

type IntegerEnum map[int]string

func (i IntegerEnum) Kind() string {
	return "integer_enum"
}

func (i IntegerEnum) Validate() error {
	if i == nil {
		return fmt.Errorf("integer enum cannot be empty")
	}

	for integer, value := range i {
		if value == "" {
			return fmt.Errorf("integer enum value for %d cannot be empty (null values must be stringified)", integer)
		}
	}
	return nil
}

type BitMapEnum map[int]string

func (b BitMapEnum) Kind() string {
	return "bit_map_enum"
}

func (b BitMapEnum) Validate() error {
	if b == nil {
		return fmt.Errorf("bit map enum cannot be empty")
	}

	for integer, value := range b {
		if value == "" {
			return fmt.Errorf("integer enum value for %d cannot be empty (null values must be stringified)", integer)
		}
	}
	return nil
}

type OidEnum string

func (o OidEnum) Kind() string {
	return "oid_enum"
}

func (o OidEnum) Validate() error {
	if o == "" {
		return fmt.Errorf("oid enum value cannot be empty")
	}
	return nil
}
