package reader

import (
	"fmt"
	"os"

	"github.com/elastiflow/snmp/pkg/def"
	"gopkg.in/yaml.v3"
)

const (
	defaultDeviceFileName      = "device.yml"
	defaultObjectTypesFileName = "object_types.yml"
)

func ReadTrapEnterprises(rulesDir, enterpriseFile string) (map[string]def.Enterprise, error) {
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

func ReadDefinitions(
	devicesDir string,
	deviceGroupsDir string,
	objectGroupsDir string,
	objectsDir string,
	defaultsDir string,
) (*def.Definitions, error) {
	d := &def.Definitions{}
	var err error

	d.Devices, err = Read[def.Device](devicesDir)
	if err != nil {
		return nil, fmt.Errorf("failed to read and validate devices: %w", err)
	}

	d.DeviceGroups, err = Read[def.DeviceGroup](deviceGroupsDir)
	if err != nil {
		return nil, fmt.Errorf("failed to read and validate device groups: %w", err)
	}

	d.ObjectGroups, err = Read[def.ObjectGroup](objectGroupsDir)
	if err != nil {
		return nil, fmt.Errorf("failed to read and validate object groups: %w", err)
	}

	d.Objects, err = Read[def.Object](objectsDir)
	if err != nil {
		return nil, fmt.Errorf("failed to read and validate objects: %w", err)
	}

	objectTypes, err := read[def.ObjectType](fmt.Sprintf("%s/%s", defaultsDir, defaultObjectTypesFileName))
	if err == nil && len(objectTypes) > 0 {
		d.ObjectTypes = objectTypes
	}

	defaultDevices, err := read[def.Device](fmt.Sprintf("%s/%s", defaultsDir, defaultDeviceFileName))
	if err == nil && len(defaultDevices) > 0 {
		for _, defaultDevice := range defaultDevices {
			d.DefaultDevice = &defaultDevice
			break
		}
	}

	if err = d.Validate(); err != nil {
		return nil, fmt.Errorf("failed to validate definitions: %w", err)
	}

	return d, nil
}

func ReadEnums(enumsDir string) (*def.Enums, error) {
	integerEnums, err := Read[def.IntegerEnum](fmt.Sprintf("%s/integer", enumsDir))
	if err != nil {
		return nil, fmt.Errorf("failed to read and validate integer enums: %w", err)
	}

	for oid := range integerEnums {
		if !isValidOID(oid) {
			return nil, fmt.Errorf("invalid oid %s in integer enums", oid)
		}
	}

	bitMapEnums, err := Read[def.BitMapEnum](fmt.Sprintf("%s/bitmap", enumsDir))
	if err != nil {
		return nil, fmt.Errorf("failed to read and validate bit map enums: %w", err)
	}

	for oid := range bitMapEnums {
		if !isValidOID(oid) {
			return nil, fmt.Errorf("invalid oid %s in bit map enums", oid)
		}
	}

	oidEnums, err := Read[def.OidEnum](fmt.Sprintf("%s/oid", enumsDir))
	if err != nil {
		return nil, fmt.Errorf("failed to read and validate oid enums: %w", err)
	}

	for oid := range oidEnums {
		if !isValidOID(oid) {
			return nil, fmt.Errorf("invalid oid %s in oid enums", oid)
		}
	}

	return &def.Enums{
		Integers: integerEnums,
		BitMaps:  bitMapEnums,
		Oids:     oidEnums,
	}, nil
}

type Validator interface {
	Validate() error
	Kind() string
}

func Read[T Validator](dirPath string) (map[string]T, error) {
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

		for id, definition := range definitionsInFile {
			// ensure that each definition id is unique
			if _, ok := definitions[id]; ok {
				return nil, fmt.Errorf("found duplicate %s with id %s", definition.Kind(), id)
			}

			// ensure that each definition is valid
			err := definition.Validate()
			if err != nil {
				return nil, fmt.Errorf("failed to validate %s with id %s: %w", definition.Kind(), id, err)
			}

			definitions[id] = definition
		}
	}

	return definitions, nil
}

func read[T Validator](filePath string) (map[string]T, error) {
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
