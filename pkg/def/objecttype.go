package def

import (
	"fmt"
	"strings"
)

// ObjectTypeDefinition represents the structure of object_types.yml
type ObjectTypeDefinition struct {
	ObjectTypes         []ObjectTypeConfig `yaml:"object_types,omitempty" json:"object_types,omitempty"`
	DefaultPollInterval uint64             `yaml:"default_poll_interval,omitempty" json:"default_poll_interval,omitempty"`
}

// ObjectTypeConfig represents a single object type configuration
type ObjectTypeConfig struct {
	Name                string `yaml:"name" json:"name"`
	Description         string `yaml:"description,omitempty" json:"description,omitempty"`
	DefaultPollInterval uint64 `yaml:"default_poll_interval" json:"default_poll_interval"`
}

// Validate validates the object types definition against its JSON schema.
func (o ObjectTypeDefinition) Validate() error {
	return validate(o, "schemas/object_types.json")
}

// ValidObjectTypes returns a slice of valid object type names
func (o ObjectTypeDefinition) ValidObjectTypes() []string {
	validTypes := make([]string, 0, len(o.ObjectTypes))
	for _, objType := range o.ObjectTypes {
		validTypes = append(validTypes, objType.Name)
	}
	return validTypes
}

// IsValidObjectType checks if the given object type is valid
func (o ObjectTypeDefinition) IsValidObjectType(objectType string) bool {
	if objectType == "" {
		return true // Empty type is allowed (will use default)
	}
	for _, objType := range o.ObjectTypes {
		if objType.Name == objectType {
			return true
		}
	}
	return false
}

// DefaultPollIntervalForType returns the default polling interval for a specific object type
func (o ObjectTypeDefinition) DefaultPollIntervalForType(objectType string) uint64 {
	if objectType == "" {
		return o.DefaultPollInterval
	}
	for _, objType := range o.ObjectTypes {
		if objType.Name == objectType {
			return objType.DefaultPollInterval
		}
	}
	return o.DefaultPollInterval
}

// ValidateObjectType validates that an object type is valid and returns an error if not
func (o ObjectTypeDefinition) ValidateObjectType(objectType string) error {
	if objectType == "" {
		return nil // Empty type is allowed (will use defaults to ensure backwards compatibility)
	}
	if !o.IsValidObjectType(objectType) {
		validTypes := o.ValidObjectTypes()
		return fmt.Errorf("invalid object type '%s'. Valid types are: %s",
			objectType, strings.Join(validTypes, ", "))
	}
	return nil
}

// ValidatePollIntervals validates that all poll interval keys are valid object types
func (o ObjectTypeDefinition) ValidatePollIntervals(pollIntervals map[string]uint64) error {
	if pollIntervals == nil {
		return nil
	}
	for objectType := range pollIntervals {
		if err := o.ValidateObjectType(objectType); err != nil {
			return fmt.Errorf("invalid poll interval key: %w", err)
		}
	}
	return nil
}
