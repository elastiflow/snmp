// Package def provides the definitions that are to be validated.
package def

import (
	"fmt"
	"strings"
)

type Definitions struct {
	Devices       map[string]*Device
	DeviceGroups  map[string]DeviceGroup
	ObjectGroups  map[string]ObjectGroup
	Objects       map[string]Object
	ObjectTypes   map[string]ObjectType
	DefaultDevice *Device
}

func (d *Definitions) ApplyDefaults() {
	// Create a default device if none is defined
	if d.DefaultDevice == nil {
		d.DefaultDevice = &Device{}
	}

	// Apply hardcoded defaults to the default device
	d.DefaultDevice.applyHardcodedDefaults(d.ObjectTypes)

	// Apply defaults to each device
	for _, device := range d.Devices {
		device.applyDefaults(d.DefaultDevice)
	}
}

// Validate checks if the definitions are valid. All devices should be populated with required fields.
// All device and object groups should reference existing object groups and objects, respectively.
// All objects should reference existing object types.
// All attribute overrides should reference objects and attributes that exist within the device group
func (d *Definitions) Validate() error {
	err := ValidateDevices(d.Devices, d.DeviceGroups, d.ObjectTypes)
	if err != nil {
		return fmt.Errorf("failed to validate devices: %w", err)
	}

	invalidDefinitions := make([]string, 0)

	// Ensure every object group in every device group is defined.
	// Ensure every overridden object in a devivce group exists
	for deviceGroupName, deviceGroup := range d.DeviceGroups {
		objectsInDeviceGroup := make(map[string]Object)
		for _, objectGroupName := range deviceGroup.ObjectGroups {
			objectGroup, ok := d.ObjectGroups[objectGroupName]
			if !ok {
				errorMsg := fmt.Sprintf("device group %q references an undefined object group: %q", deviceGroupName, objectGroupName)
				invalidDefinitions = append(invalidDefinitions, errorMsg)
				continue
			}
			for _, objectName := range objectGroup.Objects {
				object, ok := d.Objects[objectName]
				if !ok {
					continue
				}
				objectsInDeviceGroup[objectName] = object
			}
		}

		for objectName, object := range objectsInDeviceGroup {
			for attributeName, attribute := range object.Attributes {
				targetObjectName := attribute.Overrides.Object
				targetAttributeName := attribute.Overrides.Attribute
				// no overide on this attribute
				if targetObjectName == "" && targetAttributeName == "" {
					continue
				}
				targetObject, ok := objectsInDeviceGroup[targetObjectName]
				// object not in device group
				if !ok {
					invalidDefinitions = append(invalidDefinitions, fmt.Sprintf(
						"device group %q object %q attribute %q overrides object %q, which is not included in the device group",
						deviceGroupName, objectName, attributeName, targetObjectName,
					))
					continue
				}
				// object in device group, but attribute does not exist for object
				if _, ok := targetObject.Attributes[targetAttributeName]; !ok {
					invalidDefinitions = append(invalidDefinitions, fmt.Sprintf(
						"device group %q object %q attribute %q overrides an undefined attribute %q on object %q",
						deviceGroupName, objectName, attributeName, targetAttributeName, targetObjectName,
					))
				}
			}
		}
	}

	// Ensure every object in every object group is defined.
	for objectGroupName, objectGroup := range d.ObjectGroups {
		for _, object := range objectGroup.Objects {
			if _, ok := d.Objects[object]; !ok {
				errorMsg := fmt.Sprintf("object group %q references an undefined object: %q", objectGroupName, object)
				invalidDefinitions = append(invalidDefinitions, errorMsg)
			}
		}
	}

	// Ensure every object has a valid object type (if defined)
	for _, object := range d.Objects {
		if object.Type != "" {
			if _, ok := d.ObjectTypes[object.Type]; !ok {
				errorMsg := fmt.Sprintf("object %q references an undefined object type: %q", object.ObjectName, object.Type)
				invalidDefinitions = append(invalidDefinitions, errorMsg)
			}
		}
	}

	if len(invalidDefinitions) > 0 {
		return fmt.Errorf("found %d invalid definitions:\n%s",
			len(invalidDefinitions),
			strings.Join(invalidDefinitions, "\n"),
		)
	}

	return nil
}

type Enums struct {
	Integers map[string]IntegerEnum
	BitMaps  map[string]BitMapEnum
	Oids     map[string]OidEnum
}
