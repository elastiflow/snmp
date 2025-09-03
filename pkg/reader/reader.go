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

func ReadSNMPTrapEnterprises(rulesDir, enterpriseFile string) (map[string]def.Enterprise, error) {
	fileBytes, err := os.ReadFile(enterpriseFile)
	if err != nil {
		return nil, fmt.Errorf("failed to read file %s: %w", enterpriseFile, err)
	}

	enterprises := make(map[string]def.Enterprise)
	err = yaml.Unmarshal(fileBytes, &enterprises)
	if err != nil {
		return nil, fmt.Errorf("failed to unmarshal file %s: %w", enterpriseFile, err)
	}

	for enterpriseName, enterprise := range enterprises {
		err = enterprise.Validate(rulesDir)
		if err != nil {
			return nil, fmt.Errorf("failed to validate enterprise %s: %w", enterpriseName, err)
		}
	}

	return enterprises, nil
}

func ReadSNMPDefinitions(
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
	devices, err := Read[def.Device](devicesDir)
	if err != nil {
		return nil, nil, nil, nil, fmt.Errorf("failed to read and validate devices: %w", err)
	}

	deviceGroups, err := Read[def.DeviceGroup](deviceGroupsDir)
	if err != nil {
		return nil, nil, nil, nil, fmt.Errorf("failed to read and validate device groups: %w", err)
	}

	objectGroups, err := Read[def.ObjectGroup](objectGroupsDir)
	if err != nil {
		return nil, nil, nil, nil, fmt.Errorf("failed to read and validate object groups: %w", err)
	}

	objects, err := Read[def.Object](objectsDir)
	if err != nil {
		return nil, nil, nil, nil, fmt.Errorf("failed to read and validate objects: %w", err)
	}

	if err = validateSNMPDefinitions(devices, deviceGroups, objectGroups, objects); err != nil {
		return nil, nil, nil, nil, fmt.Errorf("failed to validate definitions: %w", err)
	}

	return devices, deviceGroups, objectGroups, objects, err
}

func ReadSNMPEnums(enumsDir string) (
	map[string]def.IntegerEnum,
	map[string]def.BitMapEnum,
	map[string]def.OidEnum,
	error,
) {
	integerEnums, err := Read[def.IntegerEnum](fmt.Sprintf("%s/integer", enumsDir))
	if err != nil {
		return nil, nil, nil, fmt.Errorf("failed to read and validate integer enums: %w", err)
	}

	for oid := range integerEnums {
		if !isValidOID(oid) {
			return nil, nil, nil, fmt.Errorf("invalid oid %s in integer enums", oid)
		}
	}

	bitMapEnums, err := Read[def.BitMapEnum](fmt.Sprintf("%s/bitmap", enumsDir))
	if err != nil {
		return nil, nil, nil, fmt.Errorf("failed to read and validate bit map enums: %w", err)
	}

	for oid := range bitMapEnums {
		if !isValidOID(oid) {
			return nil, nil, nil, fmt.Errorf("invalid oid %s in bit map enums", oid)
		}
	}

	oidEnums, err := Read[def.OidEnum](fmt.Sprintf("%s/oid", enumsDir))
	if err != nil {
		return nil, nil, nil, fmt.Errorf("failed to read and validate oid enums: %w", err)
	}

	for oid := range oidEnums {
		if !isValidOID(oid) {
			return nil, nil, nil, fmt.Errorf("invalid oid %s in oid enums", oid)
		}
	}

	return integerEnums, bitMapEnums, oidEnums, nil
}

func Read[T Definition](dirPath string) (map[string]T, error) {
	filePaths, walkDirErr := walkDirectory(dirPath)
	if walkDirErr != nil {
		return nil, fmt.Errorf("failed to walk directory: %w", walkDirErr)
	}

	definitions := make(map[string]T)

	for _, filePath := range filePaths {
		definitionsInFile, readErr := read[T](filePath)
		if readErr != nil {
			return nil, fmt.Errorf("failed to read definitions: %w", readErr)
		}

		err := validateUniqueness(definitions, definitionsInFile)
		if err != nil {
			return nil, fmt.Errorf("failed to validate definitions in file %s: %w", filePath, err)
		}

		err = validateDefinition(definitionsInFile)
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

func validateDefinition[T Definition](definitions map[string]T) error {
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

func validateSNMPDefinitions(
	devices map[string]def.Device,
	deviceGroups map[string]def.DeviceGroup,
	objectGroups map[string]def.ObjectGroup,
	objects map[string]def.Object,
) error {
	err := def.ValidateDevices(devices)
	if err != nil {
		return fmt.Errorf("failed to validate devices: %w", err)
	}

	invalidDefinitions := make([]string, 0)

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
