package reader

import (
	"fmt"
	"os"
	"strings"

	"github.com/elastiflow/snmp/pkg/def"
	"gopkg.in/yaml.v3"
)

type Definition interface {
	Type() string
	Validate() error
}

func ReadAndValidateDefinitions(
	devicesDir string,
	deviceGroupsDir string,
	objectGroupsDir string,
	objectsDir string,
) (
	map[string]def.Device,
	map[string]def.DeviceGroup,
	map[string]def.ObjectGroup,
	map[string]def.Object,
	error,
) {
	devices, err := ReadAndValidate[def.Device](devicesDir)
	if err != nil {
		return nil, nil, nil, nil, fmt.Errorf("failed to read and validate devices: %w", err)
	}

	deviceGroups, err := ReadAndValidate[def.DeviceGroup](deviceGroupsDir)
	if err != nil {
		return nil, nil, nil, nil, fmt.Errorf("failed to read and validate device groups: %w", err)
	}

	objectGroups, err := ReadAndValidate[def.ObjectGroup](objectGroupsDir)
	if err != nil {
		return nil, nil, nil, nil, fmt.Errorf("failed to read and validate object groups: %w", err)
	}

	objects, err := ReadAndValidate[def.Object](objectsDir)
	if err != nil {
		return nil, nil, nil, nil, fmt.Errorf("failed to read and validate objects: %w", err)
	}

	if err = validateDefinitions(devices, deviceGroups, objectGroups, objects); err != nil {
		return nil, nil, nil, nil, fmt.Errorf("failed to validate definitions: %w", err)
	}

	return devices, deviceGroups, objectGroups, objects, err
}

func ReadAndValidateEnums(enumsDir string) (
	map[string]def.IntegerEnum,
	map[string]def.BitMapEnum,
	map[string]def.OidEnum,
	error,
) {
	integerEnums, err := ReadAndValidate[def.IntegerEnum](fmt.Sprintf("%s/integer", enumsDir))
	if err != nil {
		return nil, nil, nil, fmt.Errorf("failed to read and validate integer enums: %w", err)
	}

	bitMapEnums, err := ReadAndValidate[def.BitMapEnum](fmt.Sprintf("%s/bitmap", enumsDir))
	if err != nil {
		return nil, nil, nil, fmt.Errorf("failed to read and validate bit map enums: %w", err)
	}

	oidEnums, err := ReadAndValidate[def.OidEnum](fmt.Sprintf("%s/oid", enumsDir))
	if err != nil {
		return nil, nil, nil, fmt.Errorf("failed to read and validate oid enums: %w", err)
	}

	return integerEnums, bitMapEnums, oidEnums, nil
}

func ReadAndValidate[T Definition](dirPath string) (map[string]T, error) {
	filePaths, err := walkDirectory(dirPath)
	if err != nil {
		return nil, fmt.Errorf("failed to walk directory: %w", err)
	}

	definitions := make(map[string]T)

	for _, filePath := range filePaths {
		definitionsInFile, parseErr := read[T](filePath)
		if parseErr != nil {
			return nil, fmt.Errorf("failed to parse definitions: %w", err)
		}

		err = validateUniqueness(definitions, definitionsInFile)
		if err != nil {
			return nil, fmt.Errorf("failed to validate definitions in file %s: %w", filePath, err)
		}

		err = validateEach(definitionsInFile)
		if err != nil {
			return nil, fmt.Errorf("failed to validate definitions in file %s: %w", filePath, err)
		}

		mergeDefinitions(definitions, definitionsInFile)
	}

	return definitions, nil
}

func read[T Definition](filePath string) (map[string]T, error) {
	fileBytes, err := os.ReadFile(filePath)
	if err != nil {
		return nil, fmt.Errorf("failed to read file %s: %w", filePath, err)
	}

	definitions := make(map[string]T)
	err = yaml.Unmarshal(fileBytes, &definitions)
	if err != nil {
		return nil, fmt.Errorf("failed to unmarshal file %s: %w", filePath, err)
	}

	return definitions, nil
}

func validateUniqueness[T Definition](definitionsA, definitionsB map[string]T) error {
	for id := range definitionsB {
		if definition, ok := definitionsA[id]; ok {
			return fmt.Errorf("found duplicate %s with id %s", definition.Type(), id)
		}
	}
	return nil
}

func validateEach[T Definition](definitions map[string]T) error {
	for id, definition := range definitions {
		err := definition.Validate()
		if err != nil {
			return fmt.Errorf("failed to validate %s with id %s: %w", definition.Type(), id, err)
		}
	}
	return nil
}

func mergeDefinitions[T Definition](definitionsA, definitionsB map[string]T) {
	for id, definition := range definitionsB {
		definitionsA[id] = definition
	}
}

func validateDefinitions(
	devices map[string]def.Device,
	deviceGroups map[string]def.DeviceGroup,
	objectGroups map[string]def.ObjectGroup,
	objects map[string]def.Object,
) error {
	invalidDefinitions := make([]string, 0)

	// Ensure every device has a unique IP address.
	for deviceName, device := range devices {
		for otherDeviceName, otherDevice := range devices {
			if deviceName != otherDeviceName && device.IP == otherDevice.IP {
				errorMsg := fmt.Sprintf("device %s has the same IP address (%s) as device %s", deviceName, device.IP, otherDeviceName)
				invalidDefinitions = append(invalidDefinitions, errorMsg)
			}
		}
	}

	// Ensure every object group in every device group is defined.
	for deviceGroupName, deviceGroup := range deviceGroups {
		for _, objectGroup := range deviceGroup.ObjectGroups {
			if _, ok := objectGroups[objectGroup]; !ok {
				errorMsg := fmt.Sprintf("device group %s references an undefined object group: %s", deviceGroupName, objectGroup)
				invalidDefinitions = append(invalidDefinitions, errorMsg)
			}
		}
	}

	// Ensure every object in every object group is defined.
	for objectGroupName, objectGroup := range objectGroups {
		for _, object := range objectGroup.Objects {
			if _, ok := objects[object]; !ok {
				errorMsg := fmt.Sprintf("object group %s references an undefined object: %s", objectGroupName, object)
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
