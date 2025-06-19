package def

import (
	"github.com/stretchr/testify/assert"

	"testing"
)

func TestObject_Validate(t *testing.T) {
	tests := []struct {
		name     string
		object   *Object
		expected error
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
					},
				},
			},
			expected: nil,
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
			expected: nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equal(t, tt.expected, tt.object.Validate())
		})
	}
}

func TestObject_Type(t *testing.T) {
	object := Object{}
	assert.Equal(t, "object", object.Type())
}
