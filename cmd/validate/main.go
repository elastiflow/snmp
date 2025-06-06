package main

import (
	"log"

	"github.com/elastiflow/snmp/pkg/reader"
)

func main() {
	devices, deviceGroups, objectGroups, objects, defErr := reader.ReadAndValidateDefinitions(
		"devices", "device_groups", "object_groups", "objects",
	)
	if defErr != nil {
		log.Println(defErr)
	} else {
		log.Printf("Validated %d devices", len(devices))
		log.Printf("Validated %d device groups", len(deviceGroups))
		log.Printf("Validated %d object groups", len(objectGroups))
		log.Printf("Validated %d objects", len(objects))
	}

	integerEnums, bitMapEnums, oidEnums, enumErr := reader.ReadAndValidateEnums("enums")
	if enumErr != nil {
		log.Println(enumErr)
	} else {
		log.Printf("Validated %d integer enums", len(integerEnums))
		log.Printf("Validated %d bit map enums", len(bitMapEnums))
		log.Printf("Validated %d oid enums", len(oidEnums))
	}

	if defErr != nil || enumErr != nil {
		log.Fatal("Validation failed")
	}
}
