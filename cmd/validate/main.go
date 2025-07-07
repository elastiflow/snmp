package main

import (
	"log"

	"github.com/elastiflow/snmp/pkg/def"
	"github.com/elastiflow/snmp/pkg/reader"
)

func main() {
	enterprises, enterprisesErr := reader.ReadSNMPTrapEnterprises("traps/rules", "traps/enterprises.yml")
	if enterprisesErr != nil {
		log.Println(enterprisesErr)
	} else {
		log.Printf("Validated %d enterprises", len(enterprises))
	}

	devices, deviceGroups, objectGroups, objects, defErr := reader.ReadSNMPDefinitions(
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

	integerEnums, bitMapEnums, oidEnums, enumErr := reader.ReadSNMPEnums("enums")
	if enumErr != nil {
		log.Println(enumErr)
	} else {
		log.Printf("Validated %d integer enums", len(integerEnums))
		log.Printf("Validated %d bit map enums", len(bitMapEnums))
		log.Printf("Validated %d oid enums", len(oidEnums))
	}

	sysObjectIDMapping, sysErr := reader.Read[def.DeviceGroupName]("autodiscover")
	if sysErr != nil {
		log.Println(sysErr)
	} else {
		log.Printf("Validated %d sysObjectIDs", len(sysObjectIDMapping))
	}

	if enterprisesErr != nil || defErr != nil || enumErr != nil || sysErr != nil {
		log.Fatal("Validation failed")
	}

}
