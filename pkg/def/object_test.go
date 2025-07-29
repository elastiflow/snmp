package def

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestObject_Validate(t *testing.T) {
	tests := []struct {
		name        string
		object      *Object
		expectedErr string
	}{
		{
			name: "valid IF-MIB::system",
			object: &Object{
				Mib:                "IF-MIB",
				ObjectName:         "system",
				Augments:           "SNMPv2-MIB::system",
				DiscoveryAttribute: "ifNumber",
				Attributes: ObjectAttributes{
					"ifNumber": {
						Oid:    ".1.3.6.1.2.1.2.1",
						Name:   "system.netifs",
						Syntax: "Integer32",
						Metric: "gauge",
					},
				},
			},
		},
		{
			name: "valid IF-MIB::ifEntry",
			object: &Object{
				Mib:        "IF-MIB",
				ObjectName: "ifEntry",
				Index: []ObjectIndex{
					{
						Type:   "Integer32",
						Oid:    ".1.3.6.1.2.1.2.2.1.1",
						Name:   "netif",
						Syntax: "InterfaceIndex",
					},
				},
				DiscoveryAttribute: "ifDescr",
				Attributes: ObjectAttributes{
					"ifDescr": {
						Oid:    ".1.3.6.1.2.1.2.2.1.2",
						Tag:    true,
						Name:   "netif.descr",
						Syntax: "DisplayString",
					},
				},
			},
		},
		{
			name: "invalid syntax field",
			object: &Object{
				Mib:                "IF-MIB",
				ObjectName:         "system",
				Augments:           "SNMPv2-MIB::system",
				DiscoveryAttribute: "ifNumber",
				Attributes: ObjectAttributes{
					"ifNumber": {
						Oid:    ".1.3.6.1.2.1.2.1",
						Name:   "system.netifs",
						Syntax: "invalid",
						Metric: "gauge",
					},
				},
			},
			expectedErr: "value must be one of \"skip\", \"none\", \"ip_address\", \"Integer\", ",
		},
		{
			name: "invalid metric field",
			object: &Object{
				Mib:                "IF-MIB",
				ObjectName:         "system",
				Augments:           "SNMPv2-MIB::system",
				DiscoveryAttribute: "ifNumber",
				Attributes: ObjectAttributes{
					"ifNumber": {
						Oid:    ".1.3.6.1.2.1.2.1",
						Name:   "system.netifs",
						Syntax: "Integer32",
						Metric: "invalid",
					},
				},
			},
			expectedErr: "value must be one of \"gauge\", \"counter\"",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.expectedErr != "" {
				err := tt.object.Validate()
				assert.ErrorContains(t, err, tt.expectedErr)
			} else {
				assert.NoError(t, tt.object.Validate())
			}
		})
	}
}

func TestObject_Type(t *testing.T) {
	object := Object{}
	assert.Equal(t, "object", object.Type())
}
