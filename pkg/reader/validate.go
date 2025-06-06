package reader

import (
	"fmt"
	"strings"

	"github.com/elastiflow/snmp/pkg/types"
)

func Validate(
	deviceGroups map[string]types.DeviceGroup,
	objectGroups map[string]types.ObjectGroup,
	objects map[string]types.Object,
) error {
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
		invalidDefinitionsStrings := strings.Join(invalidDefinitions, "\n")
		return fmt.Errorf(fmt.Sprintf("found %d invalid definitions:\n%s", len(invalidDefinitions), invalidDefinitionsStrings))
	}

	return nil
}
