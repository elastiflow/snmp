package main

import (
	"log"

	"github.com/elastiflow/snmp/pkg/def"
	"github.com/elastiflow/snmp/pkg/reader"
)

func main() {
	enterprises, enterprisesErr := reader.ReadTrapEnterprises("traps/rules", "traps/enterprises.yml")
	if enterprisesErr != nil {
		log.Println(enterprisesErr)
	} else {
		log.Printf("Validated %d enterprises", len(enterprises))
	}

	definitions, defErr := reader.ReadDefinitions(
		"devices", "device_groups", "object_groups", "objects", "defaults",
	)
	if defErr != nil {
		log.Println(defErr)
	} else {
		log.Printf("Validated %d devices", len(definitions.Devices))
		log.Printf("Validated %d device groups", len(definitions.DeviceGroups))
		log.Printf("Validated %d object groups", len(definitions.ObjectGroups))
		log.Printf("Validated %d objects", len(definitions.Objects))
		log.Printf("Validated %d object types", len(definitions.ObjectTypes))
		if definitions.DefaultDevice != nil {
			log.Print("Validated 1 default device")
		}
	}

	enums, enumErr := reader.ReadEnums("enums")
	if enumErr != nil {
		log.Println(enumErr)
	} else {
		log.Printf("Validated %d integer enums", len(enums.Integers))
		log.Printf("Validated %d bit map enums", len(enums.BitMaps))
		log.Printf("Validated %d oid enums", len(enums.Oids))
	}

	sysObjectIDMapping, sysErr := reader.Read[def.DeviceGroupName]("autodiscover/sysoids")
	if sysErr != nil {
		log.Println(sysErr)
	} else {
		log.Printf("Validated %d sysObjectIDs", len(sysObjectIDMapping))
	}

	if enterprisesErr != nil || defErr != nil || enumErr != nil || sysErr != nil {
		log.Fatal("Validation failed")
	}

}
