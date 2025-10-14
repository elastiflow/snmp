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
	assert.Equal(t, "object", object.Kind())
}

func TestObject_GetDiscoveryOid(t *testing.T) {
	tests := []struct {
		name          string
		object        Object
		expectedOid   string
		expectedError string
	}{
		{
			name: "should return the discovery oid",
			object: Object{
				ObjectName:         "testObject",
				DiscoveryAttribute: "ifNumber",
				Attributes: ObjectAttributes{
					"ifNumber": {
						Oid: ".1.3.6.1.2.1.2.1",
					},
				},
			},
			expectedOid: ".1.3.6.1.2.1.2.1",
		},
		{
			name: "should return error if discovery attribute not present",
			object: Object{
				ObjectName:         "testObject",
				DiscoveryAttribute: "nonExistentAttribute",
				Attributes: ObjectAttributes{
					"ifNumber": {
						Oid: ".1.3.6.1.2.1.2.1",
					},
				},
			},
			expectedError: "discovery attribute not present for object \"testObject\"",
		},
		{
			name: "should return error if discovery attribute does not exist",
			object: Object{
				ObjectName:         "testObject",
				DiscoveryAttribute: "ifNumber",
				Attributes:         ObjectAttributes{},
			},
			expectedError: "discovery attribute not present for object \"testObject\"",
		},
		{
			name: "should return error if discovery attribute tag not present",
			object: Object{
				ObjectName:         "testObject",
				DiscoveryAttribute: "ifNumber",
			},
			expectedError: "discovery attribute not present for object \"testObject\"",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			oid, err := tt.object.GetDiscoveryOid()
			assert.Equal(t, tt.expectedOid, oid)
			if tt.expectedError != "" {
				assert.EqualError(t, err, tt.expectedError)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

func TestObject_GetOverrideOids(t *testing.T) {
	tests := []struct {
		name         string
		object       Object
		expectedOids []string
	}{
		{
			name: "should return all override OIDs",
			object: Object{
				Attributes: ObjectAttributes{
					"attr1": {
						Oid: ".1.3.6.1.2.1.2.1",
						Overrides: ObjectAttributeOverride{
							Object:    "someObject",
							Attribute: "someAttribute",
						},
					},
					"attr2": {
						Oid: ".1.3.6.1.2.1.2.2",
						Overrides: ObjectAttributeOverride{
							Object:    "anotherObject",
							Attribute: "anotherAttribute",
						},
					},
					"attr3": {
						Oid: ".1.3.6.1.2.1.2.3",
						// No overrides
					},
				},
			},
			expectedOids: []string{".1.3.6.1.2.1.2.1", ".1.3.6.1.2.1.2.2"},
		},
		{
			name: "should return only valid override OIDs",
			object: Object{
				Attributes: ObjectAttributes{
					"attr1": {
						Oid: ".1.3.6.1.2.1.2.1",
						Overrides: ObjectAttributeOverride{
							Object:    "someObject",
							Attribute: "", // Missing attribute
						},
					},
					"attr2": {
						Oid: ".1.3.6.1.2.1.2.2",
						Overrides: ObjectAttributeOverride{
							Object:    "", // Missing object
							Attribute: "someAttribute",
						},
					},
					"attr3": {
						Oid: ".1.3.6.1.2.1.2.3",
						Overrides: ObjectAttributeOverride{
							Object:    "validObject",
							Attribute: "validAttribute",
						},
					},
				},
			},
			expectedOids: []string{".1.3.6.1.2.1.2.3"},
		},
		{
			name: "should return empty list when no overrides exist",
			object: Object{
				Attributes: ObjectAttributes{
					"attr1": {
						Oid: ".1.3.6.1.2.1.2.1",
					},
					"attr2": {
						Oid: ".1.3.6.1.2.1.2.2",
					},
				},
			},
			expectedOids: []string{},
		},
		{
			name:         "should return empty list for empty attributes",
			object:       Object{Attributes: ObjectAttributes{}},
			expectedOids: []string{},
		},
		{
			name:         "should return empty list for nil attributes",
			object:       Object{},
			expectedOids: []string{},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.ElementsMatch(t, tt.expectedOids, tt.object.GetOverrideOids())
		})
	}
}
