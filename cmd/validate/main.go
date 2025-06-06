package main

import (
	"log"

	"github.com/elastiflow/snmp/pkg/reader"
	"github.com/elastiflow/snmp/pkg/types"
)

func main() {
	deviceGroupDefinitions, err := reader.ReadAndValidate[types.DeviceGroup]("device_groups")
	if err != nil {
		panic(err)
	}

	log.Printf("found and validated %d device group definitions", len(deviceGroupDefinitions))

	objectGroupDefinitions, err := reader.ReadAndValidate[types.ObjectGroup]("object_groups")
	if err != nil {
		panic(err)
	}

	log.Printf("found and validated %d object group definitions", len(objectGroupDefinitions))

	objectDefinitions, err := reader.ReadAndValidate[types.Object]("objects")
	if err != nil {
		panic(err)
	}

	log.Printf("found and validated %d object definitions", len(objectDefinitions))

	if err = reader.Validate(deviceGroupDefinitions, objectGroupDefinitions, objectDefinitions); err != nil {
		panic(err)
	}

	integerEnumDefinitions, err := reader.ReadAndValidate[types.IntegerEnum]("enums/integer")
	if err != nil {
		panic(err)
	}

	log.Printf("found and validated %d integer enum definitions", len(integerEnumDefinitions))

	bitMapEnumDefinitions, err := reader.ReadAndValidate[types.BitMapEnum]("enums/bitmap")
	if err != nil {
		panic(err)
	}

	log.Printf("found and validated %d bit map enum definitions", len(bitMapEnumDefinitions))

	oidEnumDefinitions, err := reader.ReadAndValidate[types.OidEnum]("enums/oid")
	if err != nil {
		panic(err)
	}

	log.Printf("found and validated %d oid enum definitions", len(oidEnumDefinitions))
}
